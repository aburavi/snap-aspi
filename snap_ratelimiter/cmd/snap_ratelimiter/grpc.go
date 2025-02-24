package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"time"

	rlimiter "snap_ratelimiter/ratelimiter"

	"github.com/aburavi/snaputils/async"
	"github.com/aburavi/snaputils/proto/ratelimiter"
)

type RatelimiterServer struct {
	ratelimiter.UnimplementedRatelimiterServer
}

func (s *RatelimiterServer) PushSlidingWindow(ctx context.Context, req *ratelimiter.RatelimiterPushSlidingWindowRequest) (*ratelimiter.RatelimiterPushSlidingWindowResponse, error) {
	var future async.Future
	var err error
	data := ratelimiter.RatelimiterPushSlidingWindowResponse{}

	future = async.Exec(func() (interface{}, error) {
		rc := rlimiter.NewRateLimitter()
		push_ratelimiter, err := rc.SlidingWindow(ctx, req.ClientId, time.Now(), req.Max)
		return push_ratelimiter, err
	})

	val, err := future.Await()
	if err != nil {
		fmt.Printf("service hit failed: %s\n", err)
		return nil, err
	}

	if !val.(bool) {
		data.Status = false
		data.Description = "Success"
	} else {
		data.Status = true
		data.Description = "API Over Limit"
	}
	return &data, err
}
