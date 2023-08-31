package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/anacrolix/torrent"
	"github.com/gin-gonic/gin"

	"github.com/itka0526/gtorrent/ws"
)

type Data struct {
	URL string `json:"URL"`
}

func main() {
	router := gin.New()
	hub := ws.NewHub()
	files, err := ws.NewFileHandler()
	if err != nil {
		log.Panic("File handler could not be created. ", err)
	}

	go hub.Run(files)

	router.GET("/api/ws", func(ctx *gin.Context) { ws.ServeWs(hub, ctx.Writer, ctx.Request) })
	router.POST("/api/download", download)

	router.Run(":3000")
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
