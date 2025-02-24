package grpc

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/aburavi/snaputils/proto/signature"
)

func encodeSignatureAuthRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &signature.AuthSignatureRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeSignatureAuthRequest", "Masuk encodeSignatureAuthRequest")

	return dreq, nil
}

func decodeSignatureAuthResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*signature.SignatureResponse)
	return reply, nil
}

func encodeTrxSignatureRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &signature.TrxSignatureRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeTrxSignatureRequest", "Masuk encodeTrxSignatureRequest")

	return dreq, nil
}

func decodeTrxSignatureResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*signature.SignatureResponse)
	return reply, nil
}
