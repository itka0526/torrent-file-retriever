package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/anacrolix/torrent"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.POST("/api/metadata", getMetadata)

	// router.POST("/download", download)

	// router.Run(":3000")
	// router := gin.New()

	// // router.LoadHTMLGlob("./templates/**/*.html")

	// router.GET("/", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "Hell World", "Gicn")
	// 	// ctx.HTML(http.StatusOK, "index.html", gin.H{
	// 	// 	"title": "No",
	// 	// })
	// })
	// // router.POST("/download", download)

	// router.Run(":3000")

	// router := gin.New()

	// router.LoadHTMLGlob("client/*.html")

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{
	// 		"title": "Torrent Downloader",
	// 	})
	// })

	router.Run(":3000")
}

type ReqBodyMD struct {
	URL string `json:"URL"`
}

func getMetadata(ctx *gin.Context) {
	reqBodyMD := ReqBodyMD{}

	err := ctx.BindJSON(&reqBodyMD)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You fucker, you fucked up!", "data": nil})
		return
	}

	client, _ := torrent.NewClient(nil)
	defer client.Close()

	t, err := client.AddMagnet(reqBodyMD.URL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "No torrents found.", "data": nil})
		return
	}

	<-t.GotInfo()

	ctx.JSON(http.StatusOK, gin.H{"message": "Torrent found", "data": t.Info()})
}

func download(ctx *gin.Context) {
	data, _ := io.ReadAll(ctx.Request.Body)

	fmt.Println(string(data))
}

func torrentExample() {
	client, _ := torrent.NewClient(nil)
	defer client.Close()

	t, _ := client.AddMagnet("magnet:?xt=urn:btih:485CDA1FC4DF71F5F81D5918D761DD8F199FF581&dn=Titans+2018+S01E06+XviD-AFG+%5Beztv%5D&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftorrent.gresille.org%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce")
	fmt.Println("Retrieved metadata...")

	<-t.GotInfo()
	msg := t.Info()

	fmt.Println(msg.Name)

	client.Close()

	// t.DownloadAll()
	// client.WaitAll()
	// log.Print("ermahgerd, torrent downloaded")
}
