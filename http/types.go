package http

import (
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
