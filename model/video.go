package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"google.golang.org/api/youtube/v3"
)

// Video VIDEOテーブル
type Video struct {
	gorm.Model
	VideoID     string
	ChannelID   string
	Title       string
	URL         string
	Length      string
	PublishedAt time.Time `gorm:"default:NULL"`
}

// NewVideo 新規Video
func NewVideo(videoID string, snippet *youtube.VideoSnippet, details *youtube.VideoContentDetails) *Video {
	URL := "https://www.youtube.com/watch?v=" + videoID
	video := &Video{VideoID: videoID,
		ChannelID:   snippet.ChannelId,
		Title:       snippet.Title,
		URL:         URL,
		Length:      details.Duration,
		PublishedAt: parseDate(snippet.PublishedAt)}
	return video
}

func parseDate(dateString string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	time, _ := time.Parse(layout, dateString)
	return time
}
