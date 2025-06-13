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

var _ svc.InquiryV2Api = (*ClientInquiryV2)(nil)

type ClientInquiryV2 struct {
	postbalanceinquiryv2         endpoint.Endpoint
	postexternalaccountinquiryv2 endpoint.Endpoint
	postinternalaccountinquiryv2 endpoint.Endpoint
}

func (c *ClientInquiryV2) PostBalanceInquiryV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.postbalanceinquiryv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientInquiryV2) PostExternalAccountInquiryV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.postexternalaccountinquiryv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientInquiryV2) PostInternalAccountInquiryV2(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.postinternalaccountinquiryv2(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewInquiryV2Client(instancer sd.Instancer, logger log.Logger) *ClientInquiryV2 {
	c := &ClientInquiryV2{}
	var options []grpctransport.ClientOption
	options = append(options)

	c.postbalanceinquiryv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"InquiryV2",
			"PostBalanceInquiryV2",
			encodeBalanceInquiryRequest,
			decodeBalanceInquiryResponse,
			inquiry.BalanceInquiryResponse{},
			options...,
		).Endpoint(), "InquiryV2().PostBalanceInquirV2y"
	})

	c.postexternalaccountinquiryv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"InquiryV2",
			"PostExternalAccountInquiryV2",
			encodeExternalAccountInquiryRequest,
			decodeExternalAccountInquiryResponse,
			inquiry.ExternalAccountInquiryResponse{},
			options...,
		).Endpoint(), "InquirV2().PostExternalAccountInquiryV2"
	})

	c.postinternalaccountinquiryv2 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"InquiryV2",
			"PostInternalAccountInquiryV2",
			encodeInternalAccountInquiryRequest,
			decodeInternalAccountInquiryResponse,
			inquiry.InternalAccountInquiryResponse{},
			options...,
		).Endpoint(), "InquiryV2().PostInternalAccountInquiryV2"
	})

	return c
}
