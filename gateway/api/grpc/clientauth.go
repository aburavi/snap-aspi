package grpc

import (
	"context"

	"github.com/aburavi/snaputils/proto/authv1"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gateway/svc"
)

var _ svc.AuthApi = (*ClientAuth)(nil)

type ClientAuth struct {
	getauthv1        endpoint.Endpoint
	getrefreshauthv1 endpoint.Endpoint
}

func (c *ClientAuth) GetAuthV1(ctx context.Context, params interface{}) (interface{}, error) {
	// Creating a new metadata object
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.getauthv1(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientAuth) GetRefreshAuthV1(ctx context.Context, params interface{}) (interface{}, error) {
	// Creating a new metadata object
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.getrefreshauthv1(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewAuthClient(instancer sd.Instancer, logger log.Logger) *ClientAuth {
	c := &ClientAuth{}
	var options []grpctransport.ClientOption
	options = append(options)

	c.getauthv1 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"AuthV1",
			"GetAuthV1",
			encodeAuthV1Request,
			decodeAuthV1Response,
			authv1.AuthV1Response{},
			options...,
		).Endpoint(), "AuthV1().GetAuthV1"
	})

	c.getrefreshauthv1 = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"AuthV1",
			"GetRefreshAuthV1",
			encodeAuthV1Request,
			decodeRefreshAuthV1Response,
			authv1.RefreshAuthV1Response{},
			options...,
		).Endpoint(), "AuthV1().GetRefreshAuthV1"
	})

	return c
}
