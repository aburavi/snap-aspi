package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/aburavi/snaputils/proto/inquiry"

	"gateway/svc"
)

var _ svc.InquiryApi = (*ClientInquiry)(nil)

type ClientInquiry struct {
	postbalanceinquiry         endpoint.Endpoint
	postexternalaccountinquiry endpoint.Endpoint
	postinternalaccountinquiry endpoint.Endpoint
}

func (c *ClientInquiry) PostBalanceInquiry(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.postbalanceinquiry(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientInquiry) PostExternalAccountInquiry(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.postexternalaccountinquiry(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientInquiry) PostInternalAccountInquiry(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.postinternalaccountinquiry(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewInquiryClient(instancer sd.Instancer, logger log.Logger) *ClientInquiry {
	c := &ClientInquiry{}
	var options []grpctransport.ClientOption
	options = append(options)

	c.postbalanceinquiry = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"Inquiry",
			"PostBalanceInquiry",
			encodeBalanceInquiryRequest,
			decodeBalanceInquiryResponse,
			inquiry.BalanceInquiryResponse{},
			options...,
		).Endpoint(), "Inquiry().PostBalanceInquiry"
	})

	c.postexternalaccountinquiry = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"Inquiry",
			"PostExternalAccountInquiry",
			encodeExternalAccountInquiryRequest,
			decodeExternalAccountInquiryResponse,
			inquiry.ExternalAccountInquiryResponse{},
			options...,
		).Endpoint(), "Inquiry().PostExternalAccountInquiry"
	})

	c.postinternalaccountinquiry = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"Inquiry",
			"PostInternalAccountInquiry",
			encodeInternalAccountInquiryRequest,
			decodeInternalAccountInquiryResponse,
			inquiry.InternalAccountInquiryResponse{},
			options...,
		).Endpoint(), "Inquiry().PostInternalAccountInquiry"
	})

	return c
}
