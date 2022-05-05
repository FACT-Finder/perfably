package state

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/FACT-Finder/perfably/config"
	"github.com/coreos/go-semver/semver"
	"github.com/rs/zerolog/log"
)

type MetaValue struct {
	Value string `json:"value"`
	URL   string `json:"url,omitempty"`
}

type (
	DataPoint map[string]float64
	MetaPoint map[string]MetaValue
)

type State struct {
	Projects map[string]*Project
}

type Project struct {
	Lock     *sync.RWMutex
	Versions semver.Versions
	Data     map[semver.Version]*VersionData
	Writer   io.WriteCloser
}

func (s *State) Close() {
	for _, p := range s.Projects {
		p.Writer.Close()
	}
}

func (p *Project) Add(line *VersionDataLine) error {
	err := json.NewEncoder(p.Writer).Encode(line)
	if err != nil {
		return err
	}

	p.addInternal(line)
	p.sortVersions()
	return nil
}

func (p *Project) sortVersions() {
	sort.Sort(p.Versions)
}

func (p *Project) addInternal(line *VersionDataLine) {
	data, ok := p.Data[line.Version]
	if !ok {
		data = &VersionData{Values: DataPoint{}, Meta: MetaPoint{}}
		p.Data[line.Version] = data
		p.Versions = append(p.Versions, &line.Version)
	}
	if line.Meta != nil {
		for key, value := range line.Meta {
			data.Meta[key] = value
		}
	}
	if line.Values != nil {
		for key, value := range line.Values {
			data.Values[key] = value
		}
	}
}

type VersionDataLine struct {
	VersionData
	Version semver.Version `json:"version"`
}

type VersionData struct {
	Values DataPoint `json:"values,omitempty"`
	Meta   MetaPoint `json:"meta,omitempty"`
}

func ReadState(config *config.Config, directory string) (*State, error) {
	err := os.MkdirAll(directory, 0o755)
	if err != nil {
		return nil, err
	}
	state := &State{Projects: map[string]*Project{}}

	for name := range config.Projects {
		var err error
		state.Projects[name], err = readProject(directory, name)
		if err != nil {
			return state, err
		}
	}

	return state, nil
}

func readProject(directory, name string) (*Project, error) {
	file := filepath.Join(directory, fmt.Sprintf("%s.v1.jsonl", name))
	handle, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)

	project := &Project{
		Lock:     &sync.RWMutex{},
		Versions: make(semver.Versions, 0),
		Data:     make(map[semver.Version]*VersionData),
		Writer:   handle,
	}

	if os.IsNotExist(err) {
		log.Info().Str("file", file).Msg("No data available")
		return project, nil
	}
	if err != nil {
		return nil, err
	}
	now := time.Now()
	r := bufio.NewReader(handle)
	rows := 0
	errs := 0

	for {
		line, err := r.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return project, err
		}
		data := &VersionDataLine{}
		err = json.Unmarshal([]byte(line), data)
		if err != nil {
			fmt.Println(err)
			errs++
			continue
		}
		rows++
		project.addInternal(data)
	}
	duration := time.Since(now)
	rowsPerSecond := (float64(rows) + float64(errs)) / duration.Seconds()
	log.Info().
		Str("took", duration.String()).
		Int("rows", rows).
		Int("errRows", errs).
		Int("versions", len(project.Versions)).
		Float64("rows/s", math.Round(rowsPerSecond)).
		Msgf("Parsed %s", file)
	project.sortVersions()
	return project, nil
}
