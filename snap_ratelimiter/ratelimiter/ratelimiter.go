package ratelimiter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"snap_ratelimiter/models"

	"github.com/segmentio/ksuid"

	"github.com/go-redis/redis/v8"

	"github.com/joho/godotenv"
)

type RateLimitter interface {
	ReadStream(ctx context.Context, a *redis.XReadGroupArgs)
	AddStream(ctx context.Context, req models.StreamData) error
	Set(ctx context.Context, obj models.StreamData, expiration time.Duration) error
	Increment(ctx context.Context) error
	FixedWindow(ctx context.Context, userid string, now time.Time, max int64) (bool, error)
	SlidingLog(ctx context.Context, userid string, now time.Time, max int64) (bool, error)
	SlidingWindow(ctx context.Context, userid string, now time.Time, max int64) (bool, error)
}

type rateLimitter struct {
	redisCli *redis.Client
}

func NewRateLimitter() RateLimitter {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	redisaddr := os.Getenv("REDIS_SERVER")
	if redisaddr == "" {
		redisaddr = "localhost"
	}

	redisport := os.Getenv("REDIS_PORT")
	if redisport == "" {
		redisport = "6379"
	}

	redispass := os.Getenv("REDISPASS")
	if redispass == "" {
		redispass = ""
	}

	rediscli := rateLimitter{
		redisCli: redis.NewClient(&redis.Options{
			Addr:     "172.32.233.14:6379",
			Password: redispass,
			DB:       0,
		})}

	ctx := context.Background()
	pong, err := rediscli.redisCli.Ping(ctx).Result()
	log.Println("try to ping", pong, err)
	return &rediscli
}

func (r *rateLimitter) ReadStream(ctx context.Context, a *redis.XReadGroupArgs) {
	r.redisCli.XReadGroup(ctx, a)
}

func (r *rateLimitter) AddStream(ctx context.Context, req models.StreamData) error {
	b, err := json.Marshal(req.IData)
	if err != nil {
		log.Println(err)
	}

	values := map[string]interface{}{
		"flow_request_id":    ksuid.New().String()[:8],
		"network_id":         "192.168.0.1",
		"terminal_id":        "6127",
		"system_trace_audit": req.SystemTraceAudit,
		"uuid":               req.UUID,
		"trx_type":           req.TrxType,
		"data":               string(b),
	}
	// log.Println("map req ", values)
	if err := r.redisCli.XAdd(ctx, &redis.XAddArgs{
		Stream: "tx_bin_stream",
		ID:     "*",
		Values: values,
	}).Err(); err != nil {
		return err
	}
	return nil
}

func (r *rateLimitter) Set(ctx context.Context, obj models.StreamData, expiration time.Duration) error {
	json, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
		return err
	}
	// log.Println(json)

	if err := r.redisCli.Set(ctx, "key1", string(json), expiration).Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (r *rateLimitter) Increment(ctx context.Context) error {
	val, err := r.redisCli.Do(ctx, "INCR", "nur:1623560362042").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
			return err
		}
		fmt.Println(err)
	}
	// fmt.Println(val.(string))
	fmt.Println(val.(int64))

	val2, err2 := r.redisCli.Do(ctx, "EXPIRE", "nur:1623560362042", 30).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
			return err2
		}
		fmt.Println(err2)
	}
	fmt.Println(val2)

	return nil
}

func (r *rateLimitter) FixedWindow(ctx context.Context, userid string, now time.Time, max int64) (bool, error) {
	key := userid + ":" + strconv.Itoa(now.Minute())
	var incr *redis.IntCmd
	_, err := r.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute*2)
		pipe.Incr(ctx, userid+":"+strconv.Itoa(now.Year())+":"+strconv.Itoa(int(now.Month())))
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("key %v val %v \n", key, incr.Val())
	if incr.Val() > max {
		return true, nil
	}
	return false, nil
}

func (r *rateLimitter) SlidingLog(ctx context.Context, userid string, now time.Time, max int64) (bool, error) {
	key := userid + ":z"
	var cmd *redis.IntCmd
	_, err := r.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		prev := now.Add(time.Duration(-1) * time.Minute)
		pipe.ZRemRangeByScore(ctx, key, strconv.Itoa(0), strconv.FormatInt(prev.UnixNano(), 10))
		pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now.UnixNano()), Member: "z:" + strconv.FormatInt(now.UnixNano(), 10)})
		pipe.Expire(ctx, key, time.Minute*5)
		// userid:year:month (not expire)
		cmd = pipe.ZCard(ctx, key)
		pipe.Incr(ctx, userid+":"+strconv.Itoa(now.Year())+":"+strconv.Itoa(int(now.Month())))
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	log.Println(key, cmd.Val())
	if cmd.Val() > max {
		log.Println("break the request")
		return true, nil
	}
	return false, nil
}

func (r *rateLimitter) SlidingWindow(ctx context.Context, userid string, now time.Time, max int64) (bool, error) {
	log.Println("sliding window")
	key := userid + ":hash"
	var mapD *redis.StringStringMapCmd

	_, err := r.redisCli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, key, userid+":"+strconv.FormatInt(now.Unix(), 10), 1)
		pipe.Incr(ctx, userid+":"+strconv.Itoa(now.Year())+":"+strconv.Itoa(int(now.Month())))
		mapD = pipe.HGetAll(ctx, key)
		pipe.Expire(ctx, key, time.Minute*5)
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("value map hash %+v \n", mapD.Val())
	log.Println("time now ", now)
	count := 0
	for keyd, value := range mapD.Val() {
		log.Printf("keymap %s=%s \n", keyd, value)
		d := strings.Split(keyd, ":")
		unix, _ := strconv.ParseInt(d[1], 10, 64)
		t := time.Unix(unix, 0)
		log.Printf(" unix : %v , time : %v \n", unix, t)

		if now.Sub(t).Seconds() > 60 {
			log.Printf("now %v > 60 from %v so delete field \n", now, t)
			r.redisCli.HDel(ctx, key, keyd)
		} else {
			v, _ := strconv.Atoi(value)
			count = count + v
		}
	}

	log.Println("Counting from all set data hash : ", count)

	if count > int(max) {
		log.Println("break the request")
		return true, nil
	}

	return false, nil
}
