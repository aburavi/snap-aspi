package http

import (
	//"github.com/go-kit/endpoint"
	//"github.com/go-kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"

	"gateway/api"

	utils "github.com/aburavi/snaputils/utils"

	"gateway/svc"
)

func makerpcinquiryHTTPV2Handler(s svc.Service, r *mux.Router, logger log.Logger, httpOpts []httptransport.ServerOption, tracer trace.Tracer, metrics utils.MetricsMiddleware) {
	{
		ep := api.MakeGetBalanceInquiryV2Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v2.0/balance-inquiry", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeBalanceInquiryResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakeGetAccountInquiryInternalV2Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v2.0/account-inquiry-internal", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeAccountInquiryInternalResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakeGetAccountInquiryExternalV2Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v2.0/account-inquiry-external", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeAccountInquiryExternalResponse,
			httpOpts...,
		)).Methods("POST")
	}

}
