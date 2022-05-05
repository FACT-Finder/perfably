package auth

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/FACT-Finder/perfably/token"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

const (
	fileName         = "users"
	tokenLength      = 32
	passwordStrength = 12
)

func Parse(directory string) (*Auth, error) {
	auth := &Auth{
		lock:      new(sync.Mutex),
		realm:     "perfably",
		users:     map[string]string{},
		directory: directory,
	}

	_, err := auth.Reload()
	return auth, err
}

type Auth struct {
	lock      *sync.Mutex
	users     map[string]string
	realm     string
	directory string
}

func (a *Auth) Secure(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok {
			unauthorized(w, a.realm)
			return
		}

		a.lock.Lock()
		hash, ok := a.users[user]
		a.lock.Unlock()

		if !ok {
			unauthorized(w, a.realm)
			return
		}

		if !token.ComparePassword([]byte(hash), []byte(pass)) {
			unauthorized(w, a.realm)
			return
		}

		handler(w, r)
	}
}

func (a *Auth) mutate(f func(map[string]string) error) error {
	err := os.MkdirAll(a.directory, 0o755)
	if err != nil {
		return err
	}
	a.lock.Lock()
	defer a.lock.Unlock()
	err = f(a.users)
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile(a.directory, "")
	if err != nil {
		return err
	}
	for name, hash := range a.users {
		fmt.Fprintf(file, "%s:%s\n", name, hash)
	}
	to := filepath.Join(a.directory, fileName)
	return os.Rename(file.Name(), to)
}

func (a *Auth) CreateWithPW(name, password string) (string, error) {
	hashedPassword := token.CreatePassword(password, passwordStrength)
	err := a.mutate(func(m map[string]string) error {
		_, ok := m[name]
		if ok {
			return fmt.Errorf("token '%s' already exists", name)
		}
		m[name] = string(hashedPassword)
		return nil
	})
	return password, err
}

func (a *Auth) Create(name string) (string, error) {
	password := token.GenerateRandomString(tokenLength)
	return a.CreateWithPW(name, password)
}

func (a *Auth) Remove(name string) error {
	return a.mutate(func(m map[string]string) error {
		delete(m, name)
		return nil
	})
}

func (a *Auth) Names() []string {
	a.lock.Lock()
	defer a.lock.Unlock()
	result := []string{}
	for name := range a.users {
		result = append(result, name)
	}
	return result
}

func (a *Auth) Reload() (amount int, err error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.users, err = readUserCSV(a.directory)
	amount = len(a.users)
	return amount, err
}

func (a *Auth) HotReload() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	f := filepath.Join(a.directory, fileName)
	log.Info().Str("file", f).Msg("Start User Hot Reload")
	err = watcher.Add(a.directory)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case e := <-watcher.Events:
				if e.Op == fsnotify.Chmod {
					break
				}
				if !strings.HasSuffix(e.Name, fileName) {
					continue
				}
				amount, err := a.Reload()
				if err != nil {
					log.Err(err).Msg("could not reload auth")
				} else {
					log.Info().Int("amount", amount).Str("file", f).Msg("Reloaded Users")
				}
			case err := <-watcher.Errors:
				log.Err(err).Msg("could not fsnotify")
				return
			}
		}
	}()
	return nil
}

func unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(401)
	io.WriteString(w, "unauthorized")
}

type UserPW struct {
	Name string
	Pass string
}

func readUserCSV(directory string) (map[string]string, error) {
	r, err := os.Open(filepath.Join(directory, fileName))
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()

	reader := csv.NewReader(r)
	reader.Comma = ':'
	reader.Comment = '#'
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, record := range records {
		if len(record) != 2 {
			return nil, errors.New("malformed users file")
		}
		result[record[0]] = record[1]
	}
	return result, nil
}
