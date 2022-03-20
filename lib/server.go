package lib

import (
	"net/http"
	"sync/atomic"
	"time"
)

var requests uint64
var epoch = time.Unix(0, 0).Format(time.RFC1123)
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}
var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func incrementRequest() {
	atomic.AddUint64(&requests, 1)
}

// StartServer starts up the file server
func StartServer(dir string, port string, cache bool) error {
	go Printer(dir, port)
	fs := http.FileServer(http.Dir(dir))
	if cache {
		http.Handle("/", useCache(fs))
	} else {
		http.Handle("/", noCache(fs))
	}
	err := http.ListenAndServe(port, nil)
	return err
}

func noCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}
		incrementRequest()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func useCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		incrementRequest()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
Track

Get informations from Track Object (Artist, Title, Duration, ...)
$track->getArtist();
$track->getTitle();
$track->getDuration();
Get current Track
$api->getPlaybackInfo()->getCurrentTrack();
Get next Track
$api->getPlaybackInfo()->getNextTrack();
Get previous Track
$api->getPlaybackInfo()->getPreviousTrack();
Get last played Songs
//all last played songs
$api->getLastPlayed();
//filters out songs that were not sent to the server
$api->getLastPlayed(true);
Get Artwork of current or next Track
$api->getTrackArtwork();
$api->getNextTrackArtwork();
Set next Track
$api->setNextTack(123);
Read Track data by filename
$api->readTag("C:\Music\song.mp3");
Playlist

Get data from Playlist Object (Tracks, Track count)
$playlist->getTracks();
$playlist->getCount();
Get Playlist (provides all information, may be slow for some large Playlists)
//all songs
$api->getPlaylist();
//only songs from 5-10
$api->getPlaylist(5, 10);
Get Playlist (doesn't provide all information)
//all songs
$api->getPlaylist2();
//only songs the first 10 songs
$api->getPlaylist2(10);
Microphone

Get Microphone status
$api->getMicrophone();
Enable/Disable Microphone
$api->setMicrophone(true);
$api->setMicrophone(false);
Playback

Get data from Playback Object (Position, Length, State, ...)
$playback->getPosition();
$playback->getLength();
$playback->getState();
Get Playback Object
$api->getPlaybackInfo()->getPlayback();
Encoder

Get data from Encoder Object (Status, Error, ...)
$encoder->getStatus();
$encoder->getError();
Get Encoder Object
$api->getEncoderStatus();
Player

Get data from Player Object (Version, Uptime)
$player->getVersion();
$player->getUptime();
Get Player Object
$api->getStatus();
