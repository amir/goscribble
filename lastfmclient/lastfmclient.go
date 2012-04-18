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
)

type Client struct {
	Username   string
	Password   string
	SessionKey string
	AuthToken  string
}

type Session struct {
	XMLName xml.Name `xml:"lfm"`
	Error   string   `xml:"error"`
	Key     string   `xml:"session>key"`
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

	s := Session{}
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

func (c *Client) ScrobbleTrack(song mpdclient.Song, timestamp int64) error {
	if c.SessionKey == "" {
		var err error
		c.SessionKey, err = c.getMobileSession()
		if err != nil {
			return err
		}
	}
	h := md5.New()
	ts := fmt.Sprintf("%s", strconv.FormatInt(timestamp, 10))
	io.WriteString(h, "api_key"+APIKEY)
	io.WriteString(h, "artist[0]"+song.Artist)
	io.WriteString(h, "method"+METHOD_SCROBBLE_TRACK)
	io.WriteString(h, "sk"+c.SessionKey)
	io.WriteString(h, "timestamp[0]"+ts)
	io.WriteString(h, "track[0]"+song.Title)
	io.WriteString(h, SECRET)
	apiSig := fmt.Sprintf("%x", h.Sum(nil))
	resp, err := http.PostForm(ENDPOINT, url.Values{
		"api_key":      {APIKEY},
		"artist[0]":    {song.Artist},
		"method":       {METHOD_SCROBBLE_TRACK},
		"sk":           {c.SessionKey},
		"timestamp[0]": {ts},
		"track[0]":     {song.Title},
		"api_sig":      {apiSig},
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
