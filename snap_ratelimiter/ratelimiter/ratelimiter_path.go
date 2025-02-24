package ratelimiter

import (
	"context"
	"log"
	"time"
)

func (r *rateLimitter) SlidingWindowPerAPI(ctx context.Context, userid, api string, now time.Time, max int64) (bool, error) {
	log.Println("sliding window per api")
	// need to be discussed
	// when using key userid:api
	// it is still need calculating key userid?
	// if not so just update key hash
	return r.SlidingWindow(ctx, userid+":"+api, now, max)
}
