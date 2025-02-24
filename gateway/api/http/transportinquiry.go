package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/aburavi/snaputils/proto/inquiry"
)

func encodeBalanceInquiryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	data := &inquiry.BalanceInquiryResponse{}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	if err, ok := response.(error); ok {
		errorEncoder(ctx, err, w)
		//return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rsp := response.([]byte)

	rsperr := json.Unmarshal(rsp, &data)
	if rsperr != nil {
		//return nil
	}

	level.Info(logger).Log("decode payload", fmt.Sprintf("%v", data))
	//level.Info(logger).Log("encode data", fmt.Sprintf("%v", data))

	return json.NewEncoder(w).Encode(data)
}

func encodeAccountInquiryInternalResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	data := &inquiry.InternalAccountInquiryResponse{}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	if err, ok := response.(error); ok {
		errorEncoder(ctx, err, w)
		//return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rsp := response.([]byte)

	rsperr := json.Unmarshal(rsp, &data)
	if rsperr != nil {
		//return nil
	}

	level.Info(logger).Log("decode payload", fmt.Sprintf("%v", data))
	//level.Info(logger).Log("encode data", fmt.Sprintf("%v", data))

	return json.NewEncoder(w).Encode(data)
}

func encodeAccountInquiryExternalResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	data := &inquiry.ExternalAccountInquiryResponse{}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	if err, ok := response.(error); ok {
		errorEncoder(ctx, err, w)
		//return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rsp := response.([]byte)

	rsperr := json.Unmarshal(rsp, &data)
	if rsperr != nil {
		//return nil
	}

	level.Info(logger).Log("decode payload", fmt.Sprintf("%v", data))
	//level.Info(logger).Log("encode data", fmt.Sprintf("%v", data))

	return json.NewEncoder(w).Encode(data)
}
