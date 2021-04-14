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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()

	app.Use(cors.New())
	videos := app.Group("/videos")
	videos.Get("/", listVideos)
	videos.Get("video/:id", getVideo)

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

	videoPath := fmt.Sprintf("./assets/%s.mp4", c.Params("id"))
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

		chunkSize := (end - start) + 1

		f, err := os.Open(videoPath)

		if err != nil {
			log.Print(fmt.Sprintf("error opening file '%s': %s", videoPath, err))
			return err
		}

		f.Seek(int64(start), 0)
		log.Printf("Reading file from %d", start)

		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		c.Set("Accept-Ranges", "bytes")
		c.Set("Content-Length", strconv.Itoa(int(chunkSize)))
		c.Set("Content-Type", "video/mp4")

		// Send in 1% at a time.
		c.Status(206)
		c.SendStream(f, int(fileSize)/100)

	} else {

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
