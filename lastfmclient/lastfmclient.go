package lastfmclient

import (
	"../mpdclient"
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ENDPOINT                  = "http://ws.audioscrobbler.com/2.0/"
	APIKEY                    = "APIKEY"
	SECRET                    = "SECRET"
	METHOD_GET_MOBILE_SESSION = "auth.getMobileSession"
	METHOD_SCROBBLE_TRACK     = "track.scrobble"
	METHOD_UPDATE_NOW_PLAYING = "track.updateNowPlaying"
)

type Client struct {
	Username   string
	Password   string
	SessionKey string
	AuthToken  string
}

type MobileSession struct {
	XMLName xml.Name `xml:"lfm"`
	Error   string   `xml:"error"`
	Key     string   `xml:"session>key"`
}

type Scrobble struct {
	Track          string `xml:"track"`
	Artist         string `xml:"artist"`
	Album          string `xml:"album"`
	IgnoredMessage string `xml:"ignoredMessage"`
}

type Scrobbles struct {
	XMLName   xml.Name   `xml:"lfm"`
	Error     string     `xml:"error"`
	Scrobbles []Scrobble `xml:"scrobbles>scrobble"`
}

type NowPlaying struct {
	XMLName    xml.Name `xml:"lfm"`
	Error      string   `xml:"error"`
	NowPlaying Scrobble `xml:"nowplaying"`
}

func NewClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
	}
}

func (c *Client) generateAuthToken() string {
	pH := md5.New()
	io.WriteString(pH, c.Password)
	tH := md5.New()
	io.WriteString(tH, c.Username)
	io.WriteString(tH, fmt.Sprintf("%x", pH.Sum(nil)))

	return fmt.Sprintf("%x", tH.Sum(nil))
}

func (c *Client) getMobileSession() (key string, err error) {
	if c.AuthToken == "" {
		c.AuthToken = c.generateAuthToken()
	}
	h := md5.New()
	io.WriteString(h, "api_key"+APIKEY)
	io.WriteString(h, "authToken"+c.AuthToken)
	io.WriteString(h, "method"+METHOD_GET_MOBILE_SESSION)
	io.WriteString(h, "username"+c.Username)
	io.WriteString(h, SECRET)
	apiSig := fmt.Sprintf("%x", h.Sum(nil))
	resp, err := http.Get(ENDPOINT + "?api_key=" + APIKEY +
		"&authToken=" + c.AuthToken +
		"&method=" + METHOD_GET_MOBILE_SESSION + "&username=" + c.Username +
		"&api_sig=" + apiSig)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	s := MobileSession{}
	err = xml.Unmarshal(body, &s)
	if err != nil {
		return
	}
	if s.Error != "" {
		err = errors.New(s.Error)
		return
	}

	key = string(s.Key)

	return
}

func (c *Client) ScrobbleTrack(song mpdclient.Song, timestamp int64) (scrobbles []Scrobble, err error) {
	if c.SessionKey == "" {
		c.SessionKey, err = c.getMobileSession()
		if err != nil {
			return
		}
	}
	h := md5.New()
	ts := fmt.Sprintf("%s", strconv.FormatInt(timestamp, 10))
	io.WriteString(h, "album[0]"+song.Album)
	io.WriteString(h, "api_key"+APIKEY)
	io.WriteString(h, "artist[0]"+song.Artist)
	io.WriteString(h, "method"+METHOD_SCROBBLE_TRACK)
	io.WriteString(h, "sk"+c.SessionKey)
	io.WriteString(h, "timestamp[0]"+ts)
	io.WriteString(h, "track[0]"+song.Title)
	io.WriteString(h, SECRET)
	apiSig := fmt.Sprintf("%x", h.Sum(nil))
	resp, err := http.PostForm(ENDPOINT, url.Values{
		"album[0]":     {song.Album},
		"api_key":      {APIKEY},
		"artist[0]":    {song.Artist},
		"method":       {METHOD_SCROBBLE_TRACK},
		"sk":           {c.SessionKey},
		"timestamp[0]": {ts},
		"track[0]":     {song.Title},
		"api_sig":      {apiSig},
	})
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	s := Scrobbles{}
	err = xml.Unmarshal(body, &s)
	if err != nil {
		return
	}
	if s.Error != "" {
		err = errors.New(s.Error)
		return
	}

	scrobbles = s.Scrobbles

	return
}

func (c *Client) UpdateNowPlaying(song mpdclient.Song) (scrobble Scrobble, err error) {
	if c.SessionKey == "" {
		c.SessionKey, err = c.getMobileSession()
		if err != nil {
			return
		}
	}
	h := md5.New()
	io.WriteString(h, "album"+song.Album)
	io.WriteString(h, "api_key"+APIKEY)
	io.WriteString(h, "artist"+song.Artist)
	io.WriteString(h, "method"+METHOD_UPDATE_NOW_PLAYING)
	io.WriteString(h, "sk"+c.SessionKey)
	io.WriteString(h, "track"+song.Title)
	io.WriteString(h, SECRET)
	apiSig := fmt.Sprintf("%x", h.Sum(nil))
	resp, err := http.PostForm(ENDPOINT, url.Values{
		"album":   {song.Album},
		"api_key": {APIKEY},
		"artist":  {song.Artist},
		"method":  {METHOD_UPDATE_NOW_PLAYING},
		"sk":      {c.SessionKey},
		"track":   {song.Title},
		"api_sig": {apiSig},
	})
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	s := NowPlaying{}
	err = xml.Unmarshal(body, &s)
	if err != nil {
		return
	}
	if s.Error != "" {
		err = errors.New(s.Error)
		return
	}

	scrobble = s.NowPlaying

	return
}
