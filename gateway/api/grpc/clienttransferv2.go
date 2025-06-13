package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/aburavi/snaputils/proto/transfer"

	"gateway/svc"
)

var _ svc.TransferV2Api = (*ClientTransferV2)(nil)

type ClientTransferV2 struct {
	posttransferinterbankv2 endpoint.Endpoint
	posttransferintrabankv2 endpoint.Endpoint
	postpaymenthosttohostv2 endpoint.Endpoint
	posttransferstatusv2    endpoint.Endpoint
}

func (c *ClientTransferV2) PostTransferInterBankV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.posttransferinterbankv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientTransferV2) PostTransferIntraBankV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.posttransferintrabankv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientTransferV2) PostPaymentHostToHostV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.postpaymenthosttohostv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientTransferV2) PostTransferStatusV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.posttransferstatusv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewTransferV2Client(instancer sd.Instancer, logger log.Logger) *ClientTransferV2 {
	c := &ClientTransferV2{}
	var options []grpctransport.ClientOption
	options = append(options)

	c.posttransferinterbankv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"TransferV2",
			"PostTransferInterBankV2",
			encodeTransferInterBankRequest,
			decodeTransferInterBankResponse,
			transfer.TransferInterBankResponse{},
			options...,
		).Endpoint(), "TransferV2().PostBalanceInquiryV2"
	})

	c.posttransferintrabankv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"TransferV2",
			"PostTransferIntraBankV2",
			encodeTransferIntraBankRequest,
			decodeTransferIntraBankResponse,
			transfer.TransferIntraBankResponse{},
			options...,
		).Endpoint(), "TransferV2().PostExternalAccountInquiryV2"
	})

	c.postpaymenthosttohostv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"TransferV2",
			"PostPaymentHostToHostV2",
			encodePaymentHostToHostRequest,
			decodePaymentHostToHostResponse,
			transfer.TransferPaymentHostoHostResponse{},
			options...,
		).Endpoint(), "TransferV2().PostPaymentHostToHostV2"
	})

	c.posttransferstatusv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"TransferV2",
			"PostTransferStatusV2",
			encodeTransferStatusRequest,
			decodeTransferStatusResponse,
			transfer.TransferStatusResponse{},
			options...,
		).Endpoint(), "TransferV2().PostInternalAccountInquiryV2"
	})

	return c
}
