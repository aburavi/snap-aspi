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

var _ svc.HistoryApi = (*ClientHistory)(nil)

type ClientHistory struct {
	gettransactionhistorydetail endpoint.Endpoint
	gettransactionhistorylist   endpoint.Endpoint
}

func (c *ClientHistory) GetTransactionHistoryDetail(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.gettransactionhistorydetail(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientHistory) GetTransactionHistoryList(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.gettransactionhistorylist(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewHistoryClient(instancer sd.Instancer, logger log.Logger) *ClientHistory {
	c := &ClientHistory{}
	var options []grpctransport.ClientOption
	options = append(options)

	c.gettransactionhistorydetail = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"History",
			"GetTransactionHistoryDetail",
			encodeHistoryDetailRequest,
			decodeHistoryDetailResponse,
			history.TransactionHistoryDetailResponse{},
			options...,
		).Endpoint(), "History().GetTransactionHistoryDetail"
	})

	c.gettransactionhistorylist = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"History",
			"GetTransactionHistoryList",
			encodeHistoryListRequest,
			decodeHistoryListResponse,
			history.TransactionHistoryListResponse{},
			options...,
		).Endpoint(), "History().GetTransactionHistoryList"
	})

	return c
}
