package trackstore

type TrackStore interface {
	AddTrack(trackId string) error
	AddTracks(trackIds []string) error
	HasTrack(trackId string) bool
}
