package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/sd"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	clientgrpc "gateway/api/grpc"
	apihttp "gateway/api/http"
	"gateway/svc"

	"github.com/aburavi/snaputils/telemetry"
)

func main() {
	// Service Name
	//SERVICE := viper.GetString("SERVICE")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// GRPC configuration
	Hostgw := viper.GetString("GRPC_HOST")
	Portgw := viper.GetString("GRPC_PORT")
	HostAuth := viper.GetString("HOSTAUTH")
	PortAuth := viper.GetString("PORTAUTH")
	HostSignature := viper.GetString("HOSTSIGNATURE")
	PortSignature := viper.GetString("PORTSIGNATURE")
	HostInquiry := viper.GetString("HOSTINQURIY")
	PortInquiry := viper.GetString("PORTINQUIRY")
	HostHistory := viper.GetString("HOSTHISTORY")
	PortHistory := viper.GetString("PORTHISTORY")
	HostTransfer := viper.GetString("HOSTTRANSFER")
	PortTransfer := viper.GetString("PORTTRANSFER")
	HostInquiryV2 := viper.GetString("HOSTV2INQUIRY")
	PortInquiryV2 := viper.GetString("PORTV2INQUIRY")
	HostHistoryV2 := viper.GetString("HOSTV2HISTORY")
	PortHistoryV2 := viper.GetString("PORTV2HISTORY")
	HostTransferV2 := viper.GetString("HOSTV2TRANSFER")
	PortTransferV2 := viper.GetString("PORTV2TRANSFER")

	fs := flag.NewFlagSet("snap-gateway", flag.ExitOnError)
	var (
		//debugAddr      = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		httpAddr              = fs.String(Hostgw, ":"+Portgw, "HTTP listen address")
		rpcsignatureInstance  = fs.String(HostSignature, HostSignature+":"+PortSignature, "Instance of Signature service")
		rpcauthInstance       = fs.String(HostAuth, HostAuth+":"+PortAuth, "Instance of Auth service")
		rpcinquiryInstance    = fs.String(HostInquiry, HostInquiry+":"+PortInquiry, "Instance of Inquiry service")
		rpchistoryInstance    = fs.String(HostHistory, HostHistory+":"+PortHistory, "Instance of History service")
		rpctransferInstance   = fs.String(HostTransfer, HostTransfer+":"+PortTransfer, "Instance of Transfer service")
		rpcinquiryv2Instance  = fs.String(HostInquiryV2, HostInquiryV2+":"+PortInquiryV2, "Instance of InquiryV2 service")
		rpchistoryv2Instance  = fs.String(HostHistoryV2, HostHistoryV2+":"+PortHistoryV2, "Instance of HistoryV2 service")
		rpctransferv2Instance = fs.String(HostTransferV2, HostTransferV2+":"+PortTransferV2, "Instance of TransferV2 service")
	)
	errc := make(chan error)

	// initialize context
	ctx := context.Background()

	// observability
	SERVICENAME := viper.GetString("SERVICENAME")
	OTELCOLLECTORURL := viper.GetString("OTEL-COLLECTOR-URL")
	TRACERNAME := viper.GetString("TRACER_NAME")

	TRACERENABLED := viper.GetBool("TRACER_ENABLED")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := telemetry.InitProvider(ctx, SERVICENAME, OTELCOLLECTORURL)
	if err != nil {
		level.Info(logger).Log("decode payload", fmt.Sprintf("%s: %v", "Failed to initialize opentelemetry provider", err))
		os.Exit(1)
	}
	defer shutdown(ctx)

	var otelTracer trace.Tracer

	otelTracerProvider := noop.NewTracerProvider()
	otelTracer = otelTracerProvider.Tracer("noop-tracer")

	if TRACERENABLED {
		otelTracer = otel.Tracer(TRACERNAME)
	}

	// contructor our service
	var (
		signaturesvc  svc.SignatureApi
		authsvc       svc.AuthApi
		inquirysvc    svc.InquiryApi
		historysvc    svc.HistoryApi
		transfersvc   svc.TransferApi
		inquiryv2svc  svc.InquiryV2Api
		historyv2svc  svc.HistoryV2Api
		transferv2svc svc.TransferV2Api
	)

	{
		authinstance := []string{*rpcauthInstance}
		authsvc = clientgrpc.NewAuthClient(sd.FixedInstancer(authinstance), logger)
	}

	{
		signatureinstance := []string{*rpcsignatureInstance}
		signaturesvc = clientgrpc.NewSignatureClient(sd.FixedInstancer(signatureinstance), logger)
	}

	{
		inquiryinstance := []string{*rpcinquiryInstance}
		inquirysvc = clientgrpc.NewInquiryClient(sd.FixedInstancer(inquiryinstance), logger)
	}

	{
		historyinstance := []string{*rpchistoryInstance}
		historysvc = clientgrpc.NewHistoryClient(sd.FixedInstancer(historyinstance), logger)
	}

	{
		transferinstance := []string{*rpctransferInstance}
		transfersvc = clientgrpc.NewTransferClient(sd.FixedInstancer(transferinstance), logger)
	}

	{
		inquiryv2instance := []string{*rpcinquiryv2Instance}
		inquiryv2svc = clientgrpc.NewInquiryClient(sd.FixedInstancer(inquiryv2instance), logger)
	}

	{
		historyv2instance := []string{*rpchistoryv2Instance}
		historyv2svc = clientgrpc.NewHistoryClient(sd.FixedInstancer(historyv2instance), logger)
	}

	{
		transferv2instance := []string{*rpctransferv2Instance}
		transferv2svc = clientgrpc.NewTransferClient(sd.FixedInstancer(transferv2instance), logger)
	}

	svc := svc.NewService(authsvc, signaturesvc, inquirysvc, historysvc, transfersvc, inquiryv2svc, historyv2svc, transferv2svc)

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		handler := apihttp.NewHTTPHandler(svc, logger, otelTracer)
		errc <- http.ListenAndServe(*httpAddr, handler)
	}()

	err = <-errc
	logger.Log(err)
}
