package lastfmclient

import (
	"../mpdclient"
	"errors"
	"testing"
	"time"
)

func Test(t *testing.T) {
	song := new(mpdclient.Song)
	song.Title = "Fearless"
	song.Artist = "Pink Floyd"
	song.Album = "Meddle"
	client := NewClient("username", "password")
	now := time.Now()
	scrobbles, err := client.ScrobbleTrack(*song, now.Unix())
	if err != nil {
		t.Error(err.Error())
	}
	for _, scrobble := range scrobbles {
		if scrobble.Artist != "Pink Floyd" {
			t.Error(errors.New("Submitted Artist was 'Pink Floyd' but got " + scrobble.Artist))
		}
	}
}
