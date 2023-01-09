package APIs

import (
	DTOs "YTSearchAPI/DTOs"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// flags
var (
	query      = flag.String("query", "Cricket", "Search term")
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
	apiKey     = flag.String("apiKey", "AIzaSyBEHkskW5N6D1aT3v52CvLIAW2DL0sAO7Y", "YT credential API Key")
)

// db details
const (
	host     = "192.168.18.88"
	port     = 55000
	user     = "postgres"
	password = "postgrespw"
	dbname   = "VideoData"
)

func FetchAndStoreVideos(wg *sync.WaitGroup) {

	flag.Parse()
	defer wg.Done()

	client := &http.Client{
		Transport: &transport.APIKey{Key: *apiKey},
	}

	service, err := youtube.New(client)
	// fmt.Printf("%T | %v", err, err)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
		wg.Done()
	}

	// Make the API call to YouTube.
	var part []string
	part = append(part, "id")
	part = append(part, "snippet")

	for {
		call := service.Search.List(part).
			Q(*query).
			MaxResults(*maxResults).
			PublishedAfter(time.Now().Add(-1 * time.Minute).Format("2006-01-02T15:04:05Z07:00")).
			PublishedBefore(time.Now().Format("2006-01-02T15:04:05Z07:00")).
			Type("video").
			Order("date")

		response, err := call.Do()

		if err != nil {
			log.Fatalf("Error calling  YouTube: %v", err)
			break
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
				video.Thumbnails = item.Snippet.Thumbnails.High.Url
				videos = append(videos, video)
			}
		}

		done, err := printNStoreInDB(videos)
		if err != nil || !done {
			log.Fatalf("Error : %v", err)
			break
		}
		if done {
			fmt.Println("Succesfully pushed in DB")
		}
		// this helps in not exhausting the daily YT api quota,
		// also buys some time to get fresh uploads
		time.Sleep(time.Minute)
	}

}

// Print and store in DB
func printNStoreInDB(videos []DTOs.Video) (bool, error) {
	for _, val := range videos {
		fmt.Println(val.Thumbnails, "-", val.PublishDate)
		fmt.Print("\n")
	}

	// DB related code

	// defining connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// we open and close connection for every tranche of videos
	// opening a DB connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to verify supplied argument")
		return false, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to establish a connection with DB")
		return false, err
	}

	sqlStatement := `INSERT INTO ytvideos 
	(vidId, title, description, publishDate, thumbnails)				
	VALUES ($1, $2, $3, $4, $5)`

	for _, val := range videos {

		_, err := db.Exec(sqlStatement, val.VidId, val.Title, val.Description,
			val.PublishDate, val.Thumbnails)
		if err != nil {
			fmt.Println("Failed to make an entry in DB", err)
			// return false, err
		}
	}

	return true, nil
}
