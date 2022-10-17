package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"jordwyatt.github.com/do-you-spotify/pkg/fsutils"
	"jordwyatt.github.com/do-you-spotify/pkg/oauth"
	"jordwyatt.github.com/do-you-spotify/pkg/trackid"
	"jordwyatt.github.com/do-you-spotify/pkg/trackstore"
)

const (
	SPOTIFY_MAX_TRACKS_ADD_PER_PLAYLIST_REQUEST = 100
)

func main() {
	fsutils.Init()

	credentialsFileExists, err := fsutils.Exists(fsutils.GetCredentialsFilePath())
	if err != nil {
		log.Fatalf("Error checking if file '%v' exists: %v", fsutils.GetCredentialsFilePath(), credentialsFileExists)
	}

	if !credentialsFileExists {
		log.Println("Fetching credentials...")
		oauth.GetCredentials()
		log.Println("Credentials successfully fetched!")
	}

	fsTrackStore, err := trackstore.NewFsTrackStore()
	if err != nil {
		log.Fatalf("Error initializing track store: %v", err)
	}

	log.Println("Beginning playlist update...")
	err = updatePlaylist(fsTrackStore)
	if err != nil {
		log.Fatalf("an error occurred when updating the playlist: %v", err)
	}
}

func updatePlaylist(trackstore trackstore.TrackStore) error {
	ctx, spotifyClient := getSpotifyClient()
	trackIdClient := trackid.NewTrackIdClient()
	err := addTracksPlayedTodayToPlaylist(ctx, trackIdClient, spotifyClient, trackstore)

	if err != nil {
		return err
	}

	return nil
}

func getSpotifyClient() (context.Context, *spotify.Client) {
	file, _ := ioutil.ReadFile(fsutils.GetCredentialsFilePath())
	token := &oauth2.Token{}
	_ = json.Unmarshal([]byte(file), token)
	ctx := context.Background()
	client := spotify.New(spotifyauth.New(spotifyauth.WithScopes(spotifyauth.ScopePlaylistModifyPrivate)).Client(ctx, token))
	return ctx, client
}

// TODO: Refactor
func addTracksPlayedTodayToPlaylist(ctx context.Context, trackIdClient *trackid.TrackIdClient, spotifyClient *spotify.Client, trackStore trackstore.TrackStore) error {
	tracks, err := trackIdClient.GetTracks()

	if err != nil {
		return err
	}

	trackIds := []string{}

	log.Println("Fetching Spotify IDs for tracks, this can take a while...")
	for i, track := range tracks {
		index := i + 1
		if track.SongLink == "" {
			log.Printf("[%v/%v] Skipping track '%v - %v' as SongLink is nil", index, len(tracks), track.Artist, track.Title)
			continue
		}
		trackId, _ := getSpotifyTrackIdFromAuddLink(track.SongLink)
		if trackId != "" {
			trackIds = append(trackIds, trackId)
			log.Printf("[%v/%v] Retrieved Spotify ID for track '%v - %v'\n", index, len(tracks), track.Artist, track.Title)
		} else {
			log.Printf("[%v/%v] Could not retrieve Spotify ID for track '%v - %v', skipping.\n", index, len(tracks), track.Artist, track.Title)
		}
	}

	log.Printf("Finished fetching Spotify track IDs. Out of %v tracks played, %v Spotify IDs were retrieved.\n", len(tracks), len(trackIds))

	log.Println("Filtering out tracks already in playlist.")
	trackIdsNotInPlaylist := getTrackIdsToAdd(trackIds, trackStore)

	log.Printf("There are a total of %v tracks to add to playlist.\n", len(trackIdsNotInPlaylist))

	for i := 0; i < len(trackIdsNotInPlaylist); i += SPOTIFY_MAX_TRACKS_ADD_PER_PLAYLIST_REQUEST {
		tracksToAdd := paginate(trackIdsNotInPlaylist, i, SPOTIFY_MAX_TRACKS_ADD_PER_PLAYLIST_REQUEST)
		err = addTracksToPlaylist(ctx, tracksToAdd, spotifyClient)
		if err != nil {
			return err
		}
	}

	trackStore.AddTracks(trackIdsNotInPlaylist)

	log.Printf("Finished adding tracks to playlist")

	return nil
}

func addTracksToPlaylist(ctx context.Context, trackIds []string, client *spotify.Client) error {
	spotifyPlaylistId := spotify.ID(os.Getenv("SPOTIFY_PLAYLIST_ID"))
	spotifyTrackIds := []spotify.ID{}

	if len(trackIds) == 0 {
		return nil
	}

	for _, trackId := range trackIds {
		spotifyTrackIds = append(spotifyTrackIds, spotify.ID(trackId))
	}

	log.Printf("Adding %v tracks to playlist with ID: %v\n", len(spotifyTrackIds), spotifyPlaylistId)
	_, err := client.AddTracksToPlaylist(ctx, spotifyPlaylistId, spotifyTrackIds...)

	if err != nil {
		return fmt.Errorf("an error occured when adding to the playlist: %v", err)
	}

	return nil
}

// TODO: Improve error handling
func getSpotifyTrackIdFromAuddLink(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to fetch content from %v", url)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)
	regex, err := regexp.Compile("https?://open.spotify.com/track/([a-zA-Z0-9]+)")
	if err != nil {
		log.Fatal(err)
	}

	matches := regex.FindStringSubmatch(bodyString)

	if len(matches) == 2 {
		return matches[1], nil
	}

	return "", nil

}

func getTrackIdsToAdd(trackIds []string, trackStore trackstore.TrackStore) []string {
	trackIdsNotInPlaylist := []string{}

	for _, trackId := range trackIds {
		if !trackStore.HasTrack(trackId) {
			trackIdsNotInPlaylist = append(trackIdsNotInPlaylist, trackId)
		}
	}

	return trackIdsNotInPlaylist
}

func paginate(x []string, skip int, size int) []string {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}
