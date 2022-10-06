package oauth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"jordwyatt.github.com/do-you-spotify/pkg/fsutils"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const (
	redirectURI = "http://localhost:8080/callback"
)

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopePlaylistModifyPrivate))
	ch    = make(chan *oauth2.Token)
	state = "abc123"
)

func GetCredentials() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	log.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	token := <-ch

	err := persistCredentials(token)

	if err != nil {
		log.Fatalf("Error persisting credentials, exiting: %v", err)
	}
}

func persistCredentials(token *oauth2.Token) error {
	file, _ := json.MarshalIndent(token, "", " ")
	err := ioutil.WriteFile(fsutils.GetCredentialsFilePath(), file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	ch <- token
}
