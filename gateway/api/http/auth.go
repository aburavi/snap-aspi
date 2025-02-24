package http

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"

	"gateway/api"

	"github.com/aburavi/snaputils/utils"

	"gateway/svc"
)

func makerpcauthHTTPHandler(s svc.Service, r *mux.Router, logger log.Logger, httpOpts []httptransport.ServerOption, tracer trace.Tracer, metrics utils.MetricsMiddleware) {
	{
		ep := api.MakeGetAuthV1Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/access-token/b2b", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeAuthResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakeGetRefreshAuthV1Endpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/access-token/refresh", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeAuthResponse,
			httpOpts...,
		)).Methods("POST")
	}

}
