package youtube

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// AddVideos 指定した日付の指定したチャンネルの動画を取得する
func AddVideos(w http.ResponseWriter, r *http.Request) {

	err := godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	connectionName := os.Getenv("POSTGRES_INSTANCE_CONNECTION_NAME")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DBNAME")
	dev := os.Getenv("DEV")
	dsn := fmt.Sprintf("user=%s password=%s host=/cloudsql/%s/ dbname=%s", dbUser, dbPassword, connectionName, dbName)
	if dev == "true" {
		dsn += "sslmode=disable"
	}
	flag.Parse()

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// Start making YouTube API calls.
	now := time.Now()
	yesterday := now.Add(-time.Duration(3) * time.Hour)
	layout := "2006-01-02T15:04:05+09:00"
	layoutedDate := yesterday.Format(layout)

	db := connectDB(dsn)
	channelIDs := []Channel{}
	tmpVideos := []Video{}
	db.Find(&channelIDs)
	db.Find(&tmpVideos)
	tmpIDs := make([]string, len(tmpVideos))
	for _, tmpVideo := range tmpVideos {
		tmpIDs = append(tmpIDs, tmpVideo.VideoID)
	}
	tx := db.Begin()

	for _, channelID := range channelIDs {
		activityCall := service.Activities.List("contentDetails").ChannelId(channelID.ChannelID).PublishedAfter(layoutedDate)
		activitiesResponse, err := activityCall.Do()
		if err != nil {
			// The channels.list method call returned an error.
			log.Fatalf("Error making API call to list activities: %v", err.Error())
		}

		for _, activity := range activitiesResponse.Items {
			uploaded := activity.ContentDetails.Upload
			if uploaded == nil {
				continue
			}
			videoID := activity.ContentDetails.Upload.VideoId
			videoCall := service.Videos.List("snippet, contentDetails").Id(videoID)
			videoResponce, err := videoCall.Do()
			if err != nil {
				log.Fatalf("Error making API call to video: %v", err.Error())
			}
			videoSnippet := videoResponce.Items[0].Snippet
			videoDetails := videoResponce.Items[0].ContentDetails
			video := NewVideo(videoID, videoSnippet, videoDetails)
			for _, IDs := range tmpIDs {
				if video == nil {
					continue
				}
				if IDs == video.VideoID {
					video = nil
				}
			}

			// insert into DB
			if video != nil {
				db.Debug().Create(&video)
			}
		}
	}
	tx.Commit()
	w.Write([]byte("OK"))
}
