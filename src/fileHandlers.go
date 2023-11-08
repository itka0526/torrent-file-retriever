package src

import (
	"encoding/json"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const p = "./downloads/"

type MyFileInfo struct {
	Path    string    `json:"path"`
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modified_date"`
	IsDir   bool      `json:"is_directory"`
}

func GetFiles() []byte {
	// TODO: handle errors
	list, _ := GetFilesRaw()
	data, _ := json.Marshal(list)
	return data
}

func GetFilesRaw() ([]MyFileInfo, error) {
	files := make([]MyFileInfo, 0)
	err := filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, MyFileInfo{
			Path:    path,
			Name:    info.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func SaveFile(f *multipart.File, fh *multipart.FileHeader) error {
	newFile, err := os.Create(p + fh.Filename)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, cpErr := io.Copy(newFile, *f)
	if cpErr != nil {
		return err
	}
	return nil
}

func DeleteFile(p string) error {
	return os.Remove(p)
}
