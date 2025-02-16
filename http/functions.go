package http

import (
	"context"
	"os"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jellydator/ttlcache/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	banned = ttlcache.New[string, *atomic.Int64](
		ttlcache.WithTTL[string, *atomic.Int64](30*time.Minute),
		ttlcache.WithCapacity[string, *atomic.Int64](1000))
	banned.OnInsertion(func(ctx context.Context, i *ttlcache.Item[string, *atomic.Int64]) {
		count := i.Value()
		count.Add(1)
		log.Warn().Msgf("client ip added to banned %s at %d", i.Key(), count.Load())
	})
	for _, s := range []string{".php", ".env", "wp-login"} {
		var c atomic.Int64
		c.Store(1)
		banned.Set(s, &c, ttlcache.DefaultTTL)
	}

	go banned.Start()

	blackHole = ttlcache.New[string, time.Time](
		ttlcache.WithTTL[string, time.Time](30*time.Minute),
		ttlcache.WithCapacity[string, time.Time](1000))
	blackHole.OnInsertion(func(ctx context.Context, i *ttlcache.Item[string, time.Time]) {
		log.Warn().Msgf("client ip added to blackhole %s at %s", i.Key(), i.Value().Format(time.RFC3339))
	})
	go blackHole.Start()
}

func hardCloseConnection(c *gin.Context) {
	conn, _, err := c.Writer.Hijack()
	if err != nil {
		return
	}
	err = conn.Close()
	if err != nil {
		return
	}
}

func logSwitch(data *ginHands) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		{
			log.Warn().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	case data.StatusCode >= 500:
		{
			log.Error().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	default:
		level, err := zerolog.ParseLevel(os.Getenv("GLOBAL_LOG_LEVEL"))
		if err != nil {
			level = zerolog.InfoLevel
		}
		log.WithLevel(level).Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
	}
}
