package grpc

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/aburavi/snaputils/proto/inquiry"
)

func encodeBalanceInquiryRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &inquiry.BalanceInquiryRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeBalanceInquiryRequest", "Masuk encodeBalanceInquiryRequest")

	return dreq, nil
}

func decodeBalanceInquiryResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*inquiry.BalanceInquiryResponse)
	return reply, nil
}

func encodeExternalAccountInquiryRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &inquiry.ExternalAccountInquiryRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeAuthRequest", "Masuk encodeExternalAccountInquiryRequest")

	return dreq, nil
}

func decodeExternalAccountInquiryResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*inquiry.ExternalAccountInquiryResponse)
	return reply, nil
}

func encodeInternalAccountInquiryRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &inquiry.InternalAccountInquiryRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeInternalAccountInquiryRequest", "Masuk encodeInternalAccountInquiryRequest")

	return dreq, nil
}

func decodeInternalAccountInquiryResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*inquiry.InternalAccountInquiryResponse)
	return reply, nil
}
