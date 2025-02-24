package http

import (
	"net/http"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"

	"gateway/svc"

	"github.com/aburavi/snaputils/utils"
)

func NewHTTPHandler(s svc.Service, logger log.Logger, tracer trace.Tracer) http.Handler {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r := mux.NewRouter()

	metrics := utils.NewMetrics(RequestCount, RequestLatency, HttpDuration)

	makerpcauthHTTPHandler(s, r, logger, options, tracer, metrics)
	makerpcsignatureHTTPHandler(s, r, logger, options, tracer, metrics)
	makerpcinquiryHTTPHandler(s, r, logger, options, tracer, metrics)
	makerpctransferHTTPHandler(s, r, logger, options, tracer, metrics)
	makerpchistoryHTTPHandler(s, r, logger, options, tracer, metrics)
	MakeMetricsHandler(r, logger, options, tracer)

	return r
}
