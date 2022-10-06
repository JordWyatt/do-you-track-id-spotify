package fsutils

import (
	"log"
	"os"
	"path"
)

const (
	ConfigDirectoryName  = ".do-you-spotify"
	TrackHistoryFileName = "trackHistory.json"
	credentialsFileName  = "credentials.json"
)

func Init() error {
	// Creates a directory, along with any necessary parents, and returns nil, or else returns an error.
	// If directory already exists, MkdirAll does nothing and returns nil.
	err := os.MkdirAll(getConfigDirectoryPath(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func GetCredentialsFilePath() string {
	homeDirectoryPath := GetUserHomeDirectory()
	return path.Join(homeDirectoryPath, ConfigDirectoryName, credentialsFileName)
}

func getConfigDirectoryPath() string {
	homeDirectoryPath := GetUserHomeDirectory()
	return path.Join(homeDirectoryPath, ConfigDirectoryName)
}

// TODO: Improve error handling
func GetUserHomeDirectory() string {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return homeDirPath
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
