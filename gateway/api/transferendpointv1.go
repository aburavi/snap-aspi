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
func MakeTransferIntraBankEndpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/TransferIntraBank"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Transfer().PostTransferIntraBank(tracerCtx, dreq)
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

func MakeTransferInterBankEndpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/TransferInterBank"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Transfer().PostTransferInterBank(tracerCtx, dreq)
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

func MakePaymentHosttoHostEndpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/PostTrxPostPaymentHostToHost"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Transfer().PostPaymentHostToHost(tracerCtx, dreq)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		level.Info(logger).Log("Success PostTrxPostPaymentHostToHost", fmt.Sprintf("%v", v))
		return v, nil
	}

}

func MakeGetTransferStatusEndpoint(s svc.Service, logger log.Logger, tracer trace.Tracer, metrics utils.MetricsMiddleware) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		const op utils.Op = "Http/PostTransferStatus"
		// start tracing
		tracerCtx, span := tracer.Start(ctx, string(op))
		defer span.End()
		metrics.MetricsRequestLatency(string(op))
		metrics.MetricsHttpDuration(string(op))

		dreq := request.(map[string]interface{})
		level.Info(logger).Log("rpcname", dreq["rpcName"].(string))
		v, err := s.Transfer().PostTransferStatus(tracerCtx, dreq)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		level.Info(logger).Log("Success PostTransferStatus", fmt.Sprintf("%v", v))
		return v, nil
	}
}
