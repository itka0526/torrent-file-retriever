package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anacrolix/torrent"
	"github.com/gin-gonic/gin"

	"github.com/itka0526/gtorrent/ws"
)

type MyFileInfo struct {
	path  string
	name  string
	size  int64
	isDir bool
}

var files []MyFileInfo

type Data struct {
	URL string `json:"URL"`
}

func main() {
	router := gin.New()

	if f, err := getFiles("./downloads"); err == nil {
		files = f
	}

	hub := ws.NewHub()
	go hub.Run()

	router.GET("/api/ws", func(ctx *gin.Context) { ws.ServeWs(hub, ctx.Writer, ctx.Request) })
	router.POST("/api/download", download)

	router.Run(":3000")
}

func getFiles(path string) (files []MyFileInfo, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, MyFileInfo{path: path, name: info.Name(), isDir: info.IsDir(), size: info.Size()})
		return nil
	})
	return
}

func download(ctx *gin.Context) {
	data, err := validateJSON(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid URL"})
		return
	}

	cfg := torrent.NewDefaultClientConfig()
	// Download files to /downloads
	cfg.DataDir = "./downloads"
	// System pick any available port.
	cfg.ListenPort = 0

	client, err := torrent.NewClient(cfg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	t, _ := client.AddMagnet(data.URL)
	// Wait till the channel closes, the channel closes when it retrieves the metadata
	<-t.GotInfo()
	// save to disk and after the user downloads keep for 20 minutes and delete
	t.DownloadAll()
}

func validateJSON(body io.ReadCloser) (data Data, err error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return data, err
	}

	return data, err
}
