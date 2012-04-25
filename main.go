package main

import (
	"./lastfmclient"
	"./mpdclient"
	"flag"
	"log"
	"time"
)

var (
	currentSong    *mpdclient.Song
	mpdClient      *mpdclient.Client
	lastfmClient   *lastfmclient.Client
	In             = make(chan CurrentTrack, 100)
	mpdPassword    = flag.String("mpd_password", "", "MPD password")
	lastfmUsername = flag.String("lastfm_username", "username", "Last.fm username")
	lastfmPassword = flag.String("lastfm_password", "password", "Last.fm password")
	mpdAddress     = flag.String("mpd_address", "localhost:6600", "MPD service address")
)

type CurrentTrack struct {
	Song      mpdclient.Song
	StartTime int64
}

func listen() {
	check := time.NewTicker(time.Second)
	for {
		select {
		case <-check.C:
			playedLongEnough()
		case track := <-In:
			scrobbles, err := lastfmClient.ScrobbleTrack(track.Song, track.StartTime)
			if err != nil {
				log.Print(err)
				In <- track
			}
			for _, scrobble := range scrobbles {
				if scrobble.IgnoredMessage == "" {
					log.Printf("%s from %s by %s scrobbled\n", scrobble.Track, scrobble.Album, scrobble.Artist)
				} else {
					log.Printf("%s from %s by %s ignored: %s\n", scrobble.Track, scrobble.Album, scrobble.Artist, scrobble.IgnoredMessage)
				}
			}
		}
	}
}

func playedLongEnough() {
	status, err := mpdClient.Status()
	if err != nil {
		log.Print(err)
	}
	if status.ElapsedTime > 240 || (status.TotalTime > 30 && status.ElapsedTime > status.TotalTime/2) {
		song, err := mpdClient.CurrentSong()
		if err != nil {
			log.Print(err)
		} else {
			if currentSong.Title != song.Title && currentSong.Artist != song.Artist {
				currentSong = song
				currentTrack := new(CurrentTrack)
				now := time.Now()
				currentTrack.Song = *song
				currentTrack.StartTime = now.Unix() - int64(status.ElapsedTime)
				In <- *currentTrack
			}
			if status.ElapsedTime == status.TotalTime {
				currentSong.Album = ""
				currentSong.Artist = ""
				currentSong.Title = ""
			}
		}
	}
}

func init() {
	flag.Parse()
	var err error
	mpdClient, err = mpdclient.Dial(*mpdAddress, *mpdPassword)
	if err != nil {
		log.Fatal(err)
	}
	currentSong = new(mpdclient.Song)
	lastfmClient = lastfmclient.NewClient(*lastfmUsername, *lastfmPassword)
}

func main() {
	listen()
}
