package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gateway/svc"

	"github.com/aburavi/snaputils/proto/history"
)

var _ svc.HistoryV2Api = (*ClientHistoryV2)(nil)

type ClientHistoryV2 struct {
	gettransactionhistorydetailv2 endpoint.Endpoint
	gettransactionhistorylistv2   endpoint.Endpoint
}

func (c *ClientHistoryV2) GetTransactionHistoryDetailV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.gettransactionhistorydetailv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientHistoryV2) GetTransactionHistoryListV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.gettransactionhistorylistv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewHistoryV2Client(instancer sd.Instancer, logger log.Logger) *ClientHistoryV2 {
	c := &ClientHistoryV2{}
	var options []grpctransport.ClientOption
	options = append(options)

	c.gettransactionhistorydetailv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"HistoryV2",
			"GetTransactionHistoryDetailV2",
			encodeHistoryDetailRequest,
			decodeHistoryDetailResponse,
			history.TransactionHistoryDetailResponse{},
			options...,
		).Endpoint(), "History().GetTransactionHistoryDetailV2"
	})

	c.gettransactionhistorylistv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"HistoryV2",
			"GetTransactionHistoryListV2",
			encodeHistoryListRequest,
			decodeHistoryListResponse,
			history.TransactionHistoryListResponse{},
			options...,
		).Endpoint(), "History().GetTransactionHistoryListV2"
	})

	return c
}
