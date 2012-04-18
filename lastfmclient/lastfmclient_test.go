package lastfmclient

import (
	"../mpdclient"
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
	err := client.ScrobbleTrack(*song, now.Unix())
	if err != nil {
		t.Error(err.Error())
	}
}
