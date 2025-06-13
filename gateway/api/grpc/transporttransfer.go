package grpc

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/aburavi/snaputils/proto/transfer"
)

func encodeTransferInterBankRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &transfer.TransferInterBankRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeTransferInterBankRequest", "Masuk encodeTransferInterBankRequest")

	return dreq, nil
}

func decodeTransferInterBankResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*transfer.TransferInterBankResponse)
	return reply, nil
}

func encodeTransferIntraBankRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &transfer.TransferIntraBankRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeTransferIntraBankRequest", "Masuk encodeTransferIntraBankRequest")

	return dreq, nil
}

func decodeTransferIntraBankResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*transfer.TransferIntraBankResponse)
	return reply, nil
}

func encodePaymentHostToHostRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &transfer.TransferPaymentHostoHostRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodePaymentHostToHostRequest", "Masuk encodePaymentHostToHostRequest")

	return dreq, nil
}

func decodePaymentHostToHostResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*transfer.TransferPaymentHostoHostResponse)
	return reply, nil
}

func encodeTransferStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	//var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dreq := &transfer.TransferStatusRequest{}

	drequ := request.(map[string]interface{})
	json.Unmarshal(drequ["body"].([]byte), dreq)

	level.Info(logger).Log("encodeTransferStatusRequest", "Masuk encodeTransferStatusRequest")

	return dreq, nil
}

func decodeTransferStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*transfer.TransferStatusResponse)
	return reply, nil
}
