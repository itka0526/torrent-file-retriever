package src

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent"
)

const downloadPath = "./downloads/"

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
	err := filepath.Walk(downloadPath, func(path string, info fs.FileInfo, err error) error {
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
	newFile, err := os.Create(downloadPath + fh.Filename)
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

func DeleteFile(mfi MyFileInfo) error {
	if mfi.IsDir {
		return os.RemoveAll(mfi.Path)
	}
	return os.Remove(mfi.Path)
}

func GetFile(mfi MyFileInfo) ([]byte, error) {
	if mfi.IsDir {
		var buf bytes.Buffer
		writer := zip.NewWriter(&buf)
		if err := zipDir(writer, mfi); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	fb, err := os.ReadFile(mfi.Path)
	if err != nil {
		return nil, err
	}

	return fb, nil
}

func zipDir(writer *zip.Writer, mfi MyFileInfo) error {
	return filepath.Walk(mfi.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Store

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(mfi.Path), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func DownloadMagnet(url string) error {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = downloadPath
	cfg.ListenPort = 0

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return err
	}

	t, err := client.AddMagnet(url)
	if err != nil {
		return err
	}
	<-t.GotInfo()
	t.DownloadAll()
	return nil
}
