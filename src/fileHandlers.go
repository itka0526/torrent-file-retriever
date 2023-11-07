package src

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"time"
)

const p = "./downloads/"

type MyFileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modified_date"`
	IsDir   bool      `json:"is_directory"`
}

func GetFiles() []byte {
	// TODO: handle errors
	list, _ := GetFilesRaw()
	myFilesInfos := make([]MyFileInfo, len(list))

	for i, item := range list {
		fileInfo, _ := item.Info()
		myFilesInfos[i] = MyFileInfo{
			Name:    fileInfo.Name(),
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			IsDir:   fileInfo.IsDir(),
		}
	}
	data, _ := json.Marshal(myFilesInfos)
	return data
}

func GetFilesRaw() ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(p)

	fmt.Println(entries)
	if err != nil {
		return entries, err
	}

	return entries, nil
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
