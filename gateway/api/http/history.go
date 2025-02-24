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

func makerpchistoryHTTPHandler(s svc.Service, r *mux.Router, logger log.Logger, httpOpts []httptransport.ServerOption, tracer trace.Tracer, metrics utils.MetricsMiddleware) {
	{
		ep := api.MakeGetTransactionHistoryListEndpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/transaction-history-list", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeHistoryListResponse,
			httpOpts...,
		)).Methods("POST")
	}

	{
		ep := api.MakeGetTransactionHistoryDetailEndpoint(s, logger, tracer, metrics)
		r.Handle("/api/v1.0/transaction-history-detail", httptransport.NewServer(
			ep,
			decodeAuthRequest,
			encodeHistoryDetailResponse,
			httpOpts...,
		)).Methods("POST")
	}

}
