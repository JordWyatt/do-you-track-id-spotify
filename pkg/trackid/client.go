package trackid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TrackIdClient struct {
	BaseUrl string
}

type Track struct {
	Artist           string  `json:"artist"`
	PlayedDate       string  `json:"played_date"`
	Timecode         string  `json:"timecode"`
	SongLink         string  `json:"song_link"`
	Label            string  `json:"label"`
	Score            float64 `json:"score"`
	ReceivedDatetime string  `json:"received_datetime"`
	Album            string  `json:"album"`
	ReleaseDate      string  `json:"release_date"`
	PlayedDatetime   string  `json:"played_datetime"`
	OutOf            float64 `json:"out_of"`
	Title            string  `json:"title"`
}

type Response struct {
	Message string   `json:"message"`
	Tracks  []*Track `json:"tracks"`
}

func NewTrackIdClient() *TrackIdClient {
	return &TrackIdClient{
		BaseUrl: "https://3rqvxp6o77.execute-api.eu-central-1.amazonaws.com/api/",
	}
}

func (t *TrackIdClient) GetTodaysTracks() ([]*Track, error) {
	fmt.Println("Fetching tracks played today from Do You Track ID API...")
	resp, err := http.Get(fmt.Sprintf("%s/today", t.BaseUrl))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	trackIdResponse := Response{}

	err = json.NewDecoder(resp.Body).Decode(&trackIdResponse)

	if err != nil {
		return nil, err
	}

	log.Println("Successfully fetched tracks!")
	return trackIdResponse.Tracks, nil
}
