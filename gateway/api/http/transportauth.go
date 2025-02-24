package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aburavi/snaputils/apisignature"
	"github.com/aburavi/snaputils/proto/authv1"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func decodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var err error
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	//Check header Authorization
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || strings.ToUpper(authHeader[:6]) != "BEARER" {
		return nil, ErrUnAuthorized
	}
	token := authHeader[7:]

	jwtparse, derr := checkAuthorization(token, logger)
	if derr != nil {
		level.Info(logger).Log("parse token failed", fmt.Sprintf("%s", derr.Error()))
		return nil, errors.New("Validation Failed: " + derr.Error())
	}

	secretid := ProcessMapString("secret_id", jwtparse)

	//Check header X-App-Signature
	sigHeader := r.Header.Get("X-App-Signature")
	if sigHeader == "" {
		return nil, ErrForbidden
	}

	//Check Payload Body
	vals := r.Body
	if vals == nil {
		return nil, errors.New("Unknown/Empty Body,...")
	}

	req := map[string]interface{}{}
	header := map[string]string{}

	defer r.Body.Close()
	body, err1 := io.ReadAll(vals)
	if err1 != nil {
		return nil, errors.New("Unknown/Empty Body,...")
	}

	method := r.Method
	dpath := r.URL.Path

	bgen, err := apisignature.GenHMAC512(method, dpath, string(body), secretid, token)
	if err != nil {
		return nil, err
	}

	//Check SignatureHeader and Generate Header
	if bgen != sigHeader {
		return nil, errors.New("Error: Signature no matched,...")
	}

	for name, values := range r.Header {
		for _, value := range values {
			header[name] = value
		}
	}

	req["header"] = header
	req["body"] = body
	req["path"] = dpath

	level.Info(logger).Log("decode payload", fmt.Sprintf("%v", req["body"].(byte)))
	//level.Info(logger).Log("decode params", fmt.Sprintf("%v", req["params"].(map[string][]byte)))
	return req, nil
}

func encodeAuthResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	//var data map[string]interface{}

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
	rsp := response.(*authv1.AuthV1Response)

	level.Info(logger).Log("decode payload", fmt.Sprintf("%v", rsp))
	//level.Info(logger).Log("encode data", fmt.Sprintf("%v", data))

	return json.NewEncoder(w).Encode(rsp)
}
