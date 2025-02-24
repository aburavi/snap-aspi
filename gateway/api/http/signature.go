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

func makerpcsignatureHTTPHandler(s svc.Service, r *mux.Router, logger log.Logger, httpOpts []httptransport.ServerOption, tracer trace.Tracer, metrics utils.MetricsMiddleware) {
	{
		ep := api.MakePostAuthSignatureEndpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/utilities/signature-auth", httptransport.NewServer(
			ep,
			decodeNoTokenRequest,
			encodeSignatureResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakePostTrxSignatureEndpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/utilities/signature-service", httptransport.NewServer(
			ep,
			decodeNoTokenRequest,
			encodeSignatureResponse,
			httpOpts...,
		)).Methods("POST")
	}

}
