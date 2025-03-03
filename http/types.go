package http

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type ginHands struct {
	LogLevel   string
	SerName    string
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	MsgStr     string
}

var banned *ttlcache.Cache[string, *atomic.Int64]
var blackHole *ttlcache.Cache[string, time.Time]

type DebugCookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Quoted bool `json:"quoted"`// indicates whether the Value was originally quoted

	Path       string    `json:"path"`// optional
	Domain     string    `json:"domain"`// optional
	Expires    time.Time `json:"expires"`// optional
	RawExpires string `json:"raw_expires"`// for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge      int `json:"max_age"`
	Secure      bool `json:"secure"`
	HttpOnly    bool `json:"http_only"`
	SameSite    http.SameSite `json:"same_site"`
	Partitioned bool `json:"partitioned"`
	Raw         string   `json:"raw"`
	Unparsed    []string `json:"unparsed"`// Raw text of unparsed attribute-value pairs
}
