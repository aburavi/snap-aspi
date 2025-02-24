package http

import (
	//"github.com/go-kit/endpoint"
	//"github.com/go-kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"

	"gateway/api"
	"gateway/svc"

	utils "github.com/aburavi/snaputils/utils"
)

func makerpctransferHTTPV2Handler(s svc.Service, r *mux.Router, logger log.Logger, httpOpts []httptransport.ServerOption, tracer trace.Tracer, metrics utils.MetricsMiddleware) {
	{
		ep := api.MakeTransferIntraBankV2Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v2.0/transfer-intrabank", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeTransferIntraBankResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakeTransferInterBankV2Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v2.0/transfer-interbank", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeTransferInterBankResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakePaymentHosttoHostV2Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v2.0/debit/payment-host-to-host", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeTransferPaymentResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakeGetTransferStatusV2Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v2.0/transfer-status", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeTransferStatusResponse,
			httpOpts...,
		)).Methods("POST")
	}
}
