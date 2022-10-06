package trackstore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"jordwyatt.github.com/do-you-spotify/pkg/fsutils"
)

func getTrackStorePath() string {
	homeDirectoryPath := fsutils.GetUserHomeDirectory()
	return path.Join(homeDirectoryPath, fsutils.ConfigDirectoryName, "trackStore.json")
}

type FsTrackStore struct {
	filePath string
	store    map[string]string
}

func NewFsTrackStore() (*FsTrackStore, error) {
	filePath := getTrackStorePath()

	trackStoreExists, err := fsutils.Exists(filePath)
	if err != nil {
		return nil, err
	}

	trackStore := &FsTrackStore{
		filePath: getTrackStorePath(),
		store:    map[string]string{},
	}

	if !trackStoreExists {
		initialiseStoreFile()
	} else {
		trackStore.loadStore()
	}

	return trackStore, nil
}

func (s *FsTrackStore) AddTrack(trackId string) error {
	s.store[trackId] = time.Now().String()
	err := s.writeStore()
	return err
}

func (s *FsTrackStore) AddTracks(trackIds []string) error {
	for _, trackId := range trackIds {
		s.store[trackId] = time.Now().String()
	}
	err := s.writeStore()
	return err
}

func (s *FsTrackStore) HasTrack(trackId string) bool {
	if _, ok := s.store[trackId]; ok {
		return true
	}

	return false
}

func initialiseStoreFile() {
	emptyFile, err := os.Create(getTrackStorePath())
	if err != nil {
		log.Fatalf(err.Error())
	}
	emptyFile.Close()
}

// TODO: Error handling
func (s *FsTrackStore) loadStore() {
	file, _ := ioutil.ReadFile(getTrackStorePath())
	_ = json.Unmarshal([]byte(file), &s.store)
}

// TODO: Error handling
func (s *FsTrackStore) writeStore() error {
	file, err := json.MarshalIndent(s.store, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(getTrackStorePath(), file, 0644)
	if err != nil {
		return err
	}

	return nil
}
