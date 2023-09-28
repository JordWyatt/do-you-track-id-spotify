# Do!! You!!! Radio Track ID Spotify Playlist Generator

This project:

1. Pulls data on what tracks were played today (the day the script runs) by pulling data from https://doyoutrackid.com/
2. Pulls the Spotify track IDs where possible
3. Adds them to a Spotify playlist

It avoids adding duplicates by maintaining a map of track IDs already added on the file system.

**The code is not great, it was hacked together for a bit of fun. This is also currently a WIP and may be prone to errors**.

### Current Playlist

[Do!! You!!! Radio Track IDs](https://open.spotify.com/playlist/44FmPU3PZFiIvU9TtQuigu?si=c070f50e559d4f4b)

### Gaps / TODO

1. Improve logging, make it consistent.
2. Spotify currently limits playlists to 10k tracks. Add logic to create a new playlist when the current playlist is too large. This is a manual effort for now.
3. Using an on disk store to maintain a list of track IDs already added is risky (e.g disk failure) and will also not scale well. If I can be bothered, move this to DynamoDB.
4. Refactoring.

### Installation

If you'd like to run this yourself, follow the instructions below.

#### Prerequisites

1. Create an application on the [Spotify developer portal](https://developer.spotify.com) to obtain a client ID and client secret. Set your redirect URL to `http://localhost:8080/callback`.
2. Create a playlist you'd like to add tracks to, [obtain the playlist ID](https://clients.caster.fm/knowledgebase/110/How-to-find-Spotify-playlist-ID.html)
3. Export the following environment variables with your own ID, secret and playlist ID (note: these are made up values 😉):
   1. `SPOTIFY_ID=2d849f6e21b1f6d4bfee825e08f247g1f17d27`
   2. `SPOTIFY_SECRET=83fffb4d0fg94gdf86d9ef72bbeaea42`
   3. `SPOTIFY_PLAYLIST_ID=5siifffaf5fMwnd7Od34bu`
4. To use DynamoDB as a storage layer, export:
   1. `AWS_ACCESS_KEY_ID=AKIA4R4MVW6T7CFIEUFJ`
   2. `AWS_SECRET_ACCESS_KEY=6VFyP+TRSa/+XvRWF8VmyFLr5D0Ygo4FJ7s9Hfus`

#### Running the application

1. Clone this repository, `cd` to the repo.
2. On your local machine. run `go run main.go`.

If this is your first time running the application, you'll need to follow the instructions in the terminal to initiate the authorization code flow.
