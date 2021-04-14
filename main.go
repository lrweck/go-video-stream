// Go implementation of a video-streaming app
// Originally from https://www.smashingmagazine.com/2021/04/building-video-streaming-app-nuxtjs-node-express/
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	// Starts a new Mux
	app := fiber.New()

	// Uses CORS
	app.Use(cors.New())

	// Groups in /videos
	videos := app.Group("/videos")
	videos.Get("/", listVideos)
	videos.Get("video/:id", getVideo)

	// Create os.Signal channel to intercept errors and interruptions
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Println("Received shutdown command...")
		err := app.Shutdown()
		if err != nil {
			log.Fatal("Error shutting down:", err)
			return
		}

	}()

	if err := app.Listen(":8000"); err != nil {
		log.Panic(err)
	}

	log.Println("Shutdown complete. Bye bye!")

}

func listVideos(c *fiber.Ctx) error {
	return c.JSON(getVideos())
}

func getVideo(c *fiber.Ctx) error {

	videoName := c.Params("id")
	videoPath := fmt.Sprintf("./assets/%s.mp4", videoName)

	videoStat, err := os.Stat(videoPath)
	if err != nil {
		log.Print(fmt.Sprintf("error opening file '%s': %s", videoPath, err))
		return err
	}

	fileSize := videoStat.Size()
	videoRange := c.Get("Range")

	if videoRange != "" {

		// If range is specified, parse header and continue from where the user
		// left off
		parts := strings.Split(strings.Replace(videoRange, "bytes=", "", -1), "-")
		start, _ := strconv.Atoi(parts[0])
		end := int(fileSize - 1)
		if parts[1] != "" {
			end, _ = strconv.Atoi(parts[1])
		}

		remainingSize := (end - start) + 1

		f, err := os.Open(videoPath)
		if err != nil {
			log.Print(fmt.Sprintf("error opening file '%s': %s", videoPath, err))
			return err
		}

		// Seek file to start of Range
		f.Seek(int64(start), 0)

		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		c.Set("Accept-Ranges", "bytes")
		c.Set("Content-Length", strconv.Itoa(int(remainingSize)))
		c.Set("Content-Type", "video/mp4")

		dur := 1 * time.Minute
		for _, v := range getVideos() {
			if v.Id == videoName {
				dur = v.Duration
				break
			}
		}

		// Send in 20 seconds at a time.
		// (100MB / 100 seconds = 1MB/s) * 20
		buf := fileSize / int64(dur.Seconds()) * 20

		// Set header to Patial Content and send stream.
		c.Status(206)
		c.SendStream(f, int(buf))

	} else {

		// If no header Range is found, stream whole video.
		c.Set("Content-Length", strconv.Itoa(int(fileSize)))
		c.Set("Content-Type", "video/mp4")

		f, err := os.Open(videoPath)
		if err != nil {
			return err
		}

		c.Status(200)
		c.SendStream(f)
	}

	return nil
}
