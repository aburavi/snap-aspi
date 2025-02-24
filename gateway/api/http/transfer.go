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

func makerpctransferHTTPHandler(s svc.Service, r *mux.Router, logger log.Logger, httpOpts []httptransport.ServerOption, tracer trace.Tracer, metrics utils.MetricsMiddleware) {
	{
		ep := api.MakeTransferIntraBankEndpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/transfer-intrabank", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeTransferIntraBankResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakeTransferInterBankEndpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/transfer-interbank", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeTransferInterBankResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakePaymentHosttoHostEndpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/debit/payment-host-to-host", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeTransferPaymentResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakeGetTransferStatusEndpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/transfer-status", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeTransferStatusResponse,
			httpOpts...,
		)).Methods("POST")
	}
}
