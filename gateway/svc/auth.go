package svc

import (
	"context"
)

type AuthApi interface {
	GetAuthV1(ctx context.Context, params interface{}) (interface{}, error)
	GetRefreshAuthV1(ctx context.Context, params interface{}) (interface{}, error)
}
