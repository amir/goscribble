package mpdclient

type Song struct {
	Album  string
	Artist string
	Title  string
}

func readCurrentSong(values Values) *Song {
	song := new(Song)
	song.Album = values["Album"]
	song.Artist = values["Artist"]
	song.Title = values["Title"]

	return song
}
