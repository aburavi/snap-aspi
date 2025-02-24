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

var _ svc.TransferApi = (*ClientTransfer)(nil)

type ClientTransfer struct {
	posttransferinterbank endpoint.Endpoint
	posttransferintrabank endpoint.Endpoint
	posttransferstatus    endpoint.Endpoint
}

func (c *ClientTransfer) PostTransferInterBank(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.posttransferinterbank(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientTransfer) PostTransferIntraBank(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.posttransferintrabank(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientTransfer) PostTransferStatus(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.PostTransferStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewTransferClient(instancer sd.Instancer, logger log.Logger) *ClientTransfer {
	c := &ClientTransfer{}
	var options []grpctransport.ClientOption
	options = append(options)

	c.posttransferinterbank = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"Transfer",
			"PostTransferInterBank",
			encodeTransferInterBankRequest,
			decodeTransferInterBankResponse,
			transfer.TransferInterBankResponse{},
			options...,
		).Endpoint(), "Inquiry().PostBalanceInquiry"
	})

	c.posttransferintrabank = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"Transfer",
			"PostTransferIntraBank",
			encodeTransferIntraBankRequest,
			decodeTransferIntraBankResponse,
			transfer.TransferIntraBankResponse{},
			options...,
		).Endpoint(), "Inquiry().PostExternalAccountInquiry"
	})

	c.posttransferstatus = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"Transfer",
			"PostTransferStatus",
			encodeTransferStatusRequest,
			decodeTransferStatusResponse,
			transfer.TransferStatusResponse{},
			options...,
		).Endpoint(), "Inquiry().PostInternalAccountInquiry"
	})

	return c
}
