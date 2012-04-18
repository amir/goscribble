package mpdclient

import (
	"log"
	"testing"
)

func Test(t *testing.T) {
	client, err := Dial("localhost:6600", "")
	if err != nil {
		t.Error(err.Error())
	}
	song, err := client.CurrentSong()
	if err != nil {
		t.Error(err.Error())
	}

	status, err := client.Status()
	if err != nil {
		t.Error(err.Error())
	}
	log.Printf("%s - %s [%d/%d]", song.Artist, song.Title, status.TotalTime, status.ElapsedTime)
}
