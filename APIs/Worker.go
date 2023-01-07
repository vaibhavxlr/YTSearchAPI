package APIs

import (
	DTOs "YTSearchAPI/DTOs"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "Cricket", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
	apiKey     = flag.String("apiKey", "AIzaSyBEHkskW5N6D1aT3v52CvLIAW2DL0sAO7Y", "YT credential API Key")
)

func FetchAndStoreVideos(wg *sync.WaitGroup) {

	flag.Parse()
	client := &http.Client{
		Transport: &transport.APIKey{Key: *apiKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
		wg.Done()
	}

	// Make the API call to YouTube.
	var part []string
	// part = append(part, "id")
	part = append(part, "snippet")

	for {
		call := service.Search.List(part).
			Q(*query).
			MaxResults(*maxResults).
			PublishedAfter(time.Now().Format("2006-01-02T15:04:05Z07:00")).
			Type("video").
			Order("date")

		response, err := call.Do()

		if err != nil {
			log.Fatalf("Error calling  YouTube: %v", err)
			wg.Done()
		}

		// At max we are fetching 25 video items
		videos := make([]DTOs.Video, 0, 25)

		for _, item := range response.Items {
			// Just to be extra sure
			if item.Id.Kind == "youtube#video" {
				var video DTOs.Video
				video.VidId = item.Id.VideoId
				video.Description = item.Snippet.Description
				video.PublishDate = item.Snippet.PublishedAt
				video.Title = item.Snippet.Title

				thumbnails := make([]DTOs.Thumbnail, 0)
				var thumbnail DTOs.Thumbnail
				thumbnail.URL = item.Snippet.Thumbnails.Default.Url
				thumbnails = append(thumbnails, thumbnail)
				thumbnail.URL = item.Snippet.Thumbnails.Medium.Url
				thumbnails = append(thumbnails, thumbnail)
				thumbnail.URL = item.Snippet.Thumbnails.High.Url
				thumbnails = append(thumbnails, thumbnail)
				video.Thumbnails = thumbnails
			}
		}

		done, err := printNStoreInDB(videos)
		if err != nil || !done {
			log.Fatalf("Error creating new YouTube client: %v", err)
			wg.Done()
		}
		if done {
			fmt.Println("Succesfully pushed in DB")
		}
		time.Sleep(time.Minute)
	}

}

// Print and store in DB
func printNStoreInDB(videos []DTOs.Video) (bool, error){

	retun false, nil
}

