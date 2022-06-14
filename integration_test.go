package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/FACT-Finder/perfably/api"
	"github.com/FACT-Finder/perfably/auth"
	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/state"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestIntegration(t *testing.T) {
	workDir := filepath.Join(getCurrentPath(), "work")
	scenarioDir := filepath.Join(getCurrentPath(), "scenario")

	update := os.Getenv("UPDATE") == "true"

	scenarioFiles, err := ioutil.ReadDir(scenarioDir)
	require.NoError(t, err)

	for _, file := range scenarioFiles {
		if file.IsDir() {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			clearDir(workDir)
			scenarioFile := filepath.Join(scenarioDir, file.Name())
			body, err := ioutil.ReadFile(scenarioFile)
			require.NoError(t, err)
			var test Test
			err = yaml.Unmarshal(body, &test)
			require.NoError(t, err)

			users, err := auth.Parse(workDir)
			require.NoError(t, err)

			for name, pass := range test.Users {
				_, err := users.CreateWithPW(name, pass, 4)
				require.NoError(t, err)
			}
			s, err := state.ReadState(&test.Config, workDir)
			defer func() {
				s.Close()
			}()
			require.NoError(t, err)

			handler := api.New(&test.Config, s, users)

			for i, step := range test.Steps {
				t.Run(fmt.Sprintf("%d_%s", i, step.Name), func(t *testing.T) {
					if step.Restart {
						users, err = auth.Parse(workDir)
						require.NoError(t, err)

						s.Close()
						s, err = state.ReadState(&test.Config, workDir)
						require.NoError(t, err)

						handler = api.New(&test.Config, s, users)
						return
					}

					if step.HTTP != "" {
						parts := strings.Split(step.HTTP, " ")

						var body io.Reader
						if step.RequestBody != "" {
							b, err := json.Marshal(step.RequestBody)
							require.NoError(t, err)
							body = bytes.NewReader(b)
						}

						req, err := http.NewRequest(parts[0], parts[1], body)
						require.NoError(t, err)

						if step.Auth != "" {
							userPW := strings.Split(step.Auth, ":")
							req.SetBasicAuth(userPW[0], userPW[1])
						}

						recorder := httptest.NewRecorder()

						handler.ServeHTTP(recorder, req)

						if update {
							step.Status = recorder.Code
							step.ResponseBody = ""
							if recorder.Body.Len() != 0 {
								actual := prettyPrint(t, recorder.Body.Bytes())
								step.ResponseBody = string(actual)
							}
						}

						require.Equal(t, step.Status, recorder.Code)

						if step.ResponseBody != "" {
							actual := prettyPrint(t, recorder.Body.Bytes())
							require.Equal(t, step.ResponseBody, string(actual))
						}
					}

					if step.FileContent != nil && len(step.FileContent) > 0 {
						for file, content := range step.FileContent {
							actual, err := ioutil.ReadFile(filepath.Join(workDir, file))
							require.NoError(t, err)
							if update {
								step.FileContent[file] = string(actual)
							}
							require.Equal(t, content, string(actual))
						}
					}
				})
			}

			if update {
				newBytes, err := yaml.Marshal(test)
				require.NoError(t, err)
				err = ioutil.WriteFile(scenarioFile, newBytes, 0o644)
				require.NoError(t, err)
			}
		})
	}
	clearDir(workDir)
}

type Test struct {
	Config config.Config     `yaml:"config"`
	Users  map[string]string `yaml:"users"`
	Steps  []*TestStep       `yaml:"steps"`
}

type TestStep struct {
	Name         string            `yaml:"name"`
	Restart      bool              `yaml:"restart,omitempty"`
	HTTP         string            `yaml:"http,omitempty"`
	Auth         string            `yaml:"auth,omitempty"`
	RequestBody  interface{}       `yaml:"request_body,omitempty"`
	Status       int               `yaml:"status,omitempty"`
	ResponseBody string            `yaml:"response_body,omitempty"`
	FileContent  map[string]string `yaml:"file_content"`
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return filepath.Dir(filename)
}

func clearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func prettyPrint(t *testing.T, b []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	require.NoError(t, err)
	return out.Bytes()
}
