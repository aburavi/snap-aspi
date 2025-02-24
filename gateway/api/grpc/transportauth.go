package grpc

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aburavi/snaputils/proto/authv1"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func encodeAuthV1Request(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &authv1.AuthV1Request{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeAuthV1Request", "Masuk encodeAuthV1Request")

	return dreq, nil
}

func decodeAuthV1Response(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*authv1.AuthV1Response)
	return reply, nil
}

func encodeRefreshAuthV1Request(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &authv1.RefreshAuthV1Request{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeRefreshAuthV1Request", "Masuk encodeRefreshAuthV1Request")

	return dreq, nil
}

func decodeRefreshAuthV1Response(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*authv1.RefreshAuthV1Response)
	return reply, nil
}
