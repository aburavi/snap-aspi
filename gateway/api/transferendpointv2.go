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
func MakeTransferIntraBankV2Endpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/TransferIntraBank"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Transferv2().PostTransferIntraBankV2(tracerCtx, dreq)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		level.Info(logger).Log("Success TransferIntraBank", fmt.Sprintf("%v", v))
		return v, nil
	}
}

func MakeTransferInterBankV2Endpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/TransferInterBank"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Transferv2().PostTransferInterBankV2(tracerCtx, dreq)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		level.Info(logger).Log("Success TransferInterBank", fmt.Sprintf("%v", v))
		return v, nil
	}
}

func MakePaymentHosttoHostV2Endpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/PostTrxPaymentHosttoHost"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Transferv2().PostPaymentHostToHostV2(tracerCtx, dreq)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		level.Info(logger).Log("Success PostTrxPaymentHosttoHost", fmt.Sprintf("%v", v))
		return v, nil
	}

}

func MakeGetTransferStatusV2Endpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/PostTransferStatusV2"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Transferv2().PostTransferStatusV2(tracerCtx, dreq)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		level.Info(logger).Log("Success PostTransferStatusV2", fmt.Sprintf("%v", v))
		return v, nil
	}
}
