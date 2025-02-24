package grpc

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/aburavi/snaputils/proto/history"
)

func encodeHistoryDetailRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &history.TransactionHistoryDetailRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeHistoryDetailRequest", "Masuk encodeHistoryDetailRequest")

	return dreq, nil
}

func decodeHistoryDetailResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*history.TransactionHistoryDetailResponse)
	return reply, nil
}

func encodeHistoryListRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &history.TransactionHistoryListRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeHistoryListRequest", "Masuk encodeHistoryListRequest")

	return dreq, nil
}

func decodeHistoryListResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*history.TransactionHistoryListResponse)
	return reply, nil
}
