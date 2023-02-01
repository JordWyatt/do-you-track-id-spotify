package trackid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

func (t *TrackIdClient) GetTracks() ([]*Track, error) {
	log.Println("Fetching tracks played today from Do You Track ID API...")

	endpoint := "/today"

	// Hack - allows the program to run for previous dates in the event a previous execution has failed
	if len(os.Args) > 1 {
		endpoint = fmt.Sprintf("/archive/%v", os.Args[1])
		log.Printf("Date override provided, fetching tracks played on %s\n", os.Args[1])
	} else {
		log.Printf("Fetching tracks played today (%s)\n", time.Now().Format("02/01/2006"))
	}

	log.Printf("Calling %s", endpoint)

	resp, err := http.Get(fmt.Sprintf("%s%s", t.BaseUrl, endpoint))
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
