package api

import (
	"context"
	"fmt"

	//"encoding/json"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.opentelemetry.io/otel/trace"

	"gateway/svc"

	"github.com/aburavi/snaputils/utils"
)

// Onboarding Endpoint
func MakeGetTransactionHistoryListEndpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/GetTransactionHistoryList"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Signature().PostAuthSignature(tracerCtx, dreq)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		level.Info(logger).Log("Success PostAuthSignature", fmt.Sprintf("%v", v))
		return v, nil
	}
}

func MakeGetTransactionHistoryDetailEndpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/GetTransactionHistoryDetail"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Signature().PostTrxSignature(tracerCtx, dreq)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		level.Info(logger).Log("Success PostTrxSignature", fmt.Sprintf("%v", v))
		return v, nil
	}
}
