package svc

import (
	"context"
)

type InquiryApi interface {
	PostBalanceInquiry(ctx context.Context, params interface{}) (interface{}, error)
	PostExternalAccountInquiry(ctx context.Context, params interface{}) (interface{}, error)
	PostInternalAccountInquiry(ctx context.Context, params interface{}) (interface{}, error)
}
