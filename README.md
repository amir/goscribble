goscribble
==========

An MPD Audioscrobbler

Usage
-----
You must be granted a valid API key by last.fm in order to use this.
Go to http://www.last.fm/api/account and create a new application. Copy your `API Key`, and `secret` and paste them in `lastfmclient/lastfmclient.go`, and then:

    go build
    ./goscribble -lastfm_password="yourpassword" -lastfm_username="yourusername"
