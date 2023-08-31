package ws

import (
	"log"
	"os"
	"path/filepath"
)

type MyFileInfo struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
}

var files []MyFileInfo

func NewFileHandler() ([]MyFileInfo, error) {
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Panic("Could not get path correctly. ", err)
	}

	f, err := getFiles(projectRoot + "/downloads/")
	files = f

	return files, err
}

func getFiles(path string) (files []MyFileInfo, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file := MyFileInfo{Path: path, Name: info.Name(), IsDir: info.IsDir(), Size: info.Size()}
		files = append(files, file)
		return nil
	})
	return
}

// TODO: if update occurs in the files variable we should send a signal over the channel then afterwards we must
// send it over websocket to the client I think its stupid to constantly bombard client every 500ms or so with no updates
