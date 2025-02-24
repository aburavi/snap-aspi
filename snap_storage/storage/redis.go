package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	//"strings"
	"errors"
	"time"

	"snap_storage/models"

	"github.com/segmentio/ksuid"

	"github.com/go-redis/redis/v8"

	"github.com/joho/godotenv"
)

type Storage interface {
	ReadStream(ctx context.Context, a *redis.XReadGroupArgs)
	AddStream(ctx context.Context, req models.StreamData) error
	SetExternalId(ctx context.Context, userid string, refid string, datenow int64, expiration time.Duration) (string, error)
	GetExternalId(ctx context.Context, userid string, refid string, datenow int64) (string, error)
	SetRefferenceNo(ctx context.Context, userid string, refid string, datenow int64, expiration time.Duration) (string, error)
	SetTrxId(ctx context.Context, userid string, orirefid string, refid string, expiration time.Duration) (string, error)
	GetTrxId(ctx context.Context, userid string, orirefid string) (string, error)
}

type storage struct {
	redisCli *redis.Client
}

func NewStorage() Storage {
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

	rediscli := storage{
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

func (r *storage) ReadStream(ctx context.Context, a *redis.XReadGroupArgs) {
	r.redisCli.XReadGroup(ctx, a)
}

func (r *storage) AddStream(ctx context.Context, req models.StreamData) error {
	b, err := json.Marshal(req.IData)
	if err != nil {
		log.Println(err)
	}

	values := map[string]interface{}{
		"flow_request_id":    ksuid.New().String()[:8],
		"network_id":         "192.168.0.1",
		"terminal_id":        "6128",
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

func (r *storage) SetExternalId(ctx context.Context, userid string, refid string, datenow int64, expiration time.Duration) (string, error) {
	key := userid + ":" + strconv.FormatInt(datenow, 10) + ":" + refid
	todaykey := userid + ":" + strconv.FormatInt(datenow, 10) + "*"
	allkey := userid + "*" + refid
	fmt.Printf("key: %s\n", key)
	fmt.Printf("todaykey: %s\n", todaykey)
	fmt.Printf("allkey: %s\n", allkey)
	checktoday := true
	//checkall := true
	//extid := ""
	//check today
	iter := r.redisCli.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("keys today: ", iter.Val())
		checktoday = false
		break
	}
	if err := iter.Err(); err != nil {
		log.Println(err)
		return "failed", err
	}

	//iter2 := r.redisCli.Scan(ctx, 0, allkey, 0).Iterator()
	//for iter2.Next(ctx) {
	//	fmt.Println("keys all", iter2.Val())
	//	checkall = false
	//	break
	//}
	//if err := iter2.Err(); err != nil {
	//	log.Println(err)
	//	return "failed", err
	//}

	if checktoday {
		fmt.Printf("New inserted externalid this day\n")
		if err := r.redisCli.Set(ctx, key, refid, expiration).Err(); err != nil {
			log.Println(err)
			return "failed", err
		}
		return "insert", nil
		//} else if extid == refid {
		//	return "allow", nil
	} else {
		return "conflict", errors.New("Conflict Unique externalid today\n")
	}
}

func (r *storage) GetExternalId(ctx context.Context, userid string, refid string, datenow int64) (string, error) {
	key := userid + ":" + strconv.FormatInt(datenow, 10) + ":" + refid

	op := r.redisCli.Get(ctx, key)
	if err := op.Err(); err != nil {
		fmt.Printf("unable to GET data. error: %v", err)
		return "failed", err
	}
	res, err := op.Result()
	if err != nil {
		fmt.Printf("unable to GET data. error: %v", err)
		return "failed", err
	}

	return res, nil

}

func (r *storage) SetRefferenceNo(ctx context.Context, userid string, refid string, datenow int64, expiration time.Duration) (string, error) {
	checkall := true
	key := userid + ":" + refid
	iter := r.redisCli.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("keys all", iter.Val())
		checkall = false
		break
	}
	if err := iter.Err(); err != nil {
		log.Println(err)
		return "failed", err
	}

	if !checkall {
		fmt.Printf("Duplicate RefferenceNo data")
		return "failed", errors.New("Duplicate RefferenceNo data")
	}

	if err := r.redisCli.Set(ctx, key, datenow, expiration).Err(); err != nil {
		log.Println("unable to SET reffno data. error: %v", err)
		return "failed", err
	}

	return "success", nil

}

func (r *storage) SetTrxId(ctx context.Context, userid string, orirefid string, refid string, expiration time.Duration) (string, error) {
	key := userid + ":" + orirefid

	if err := r.redisCli.Set(ctx, key, refid, expiration).Err(); err != nil {
		log.Println("unable to SET reffno data. error: %v", err)
		return "failed", err
	}

	return "success", nil

}

func (r *storage) GetTrxId(ctx context.Context, userid string, orirefid string) (string, error) {
	key := userid + ":" + orirefid

	op := r.redisCli.Get(ctx, key)
	if err := op.Err(); err != nil {
		fmt.Printf("unable to GET data. error: %v", err)
		return "failed", err
	}
	res, err := op.Result()
	if err != nil {
		fmt.Printf("unable to GET data. error: %v", err)
		return "failed", err
	}

	return res, nil

}
