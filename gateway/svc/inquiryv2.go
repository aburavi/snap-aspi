package svc

import (
	"context"
)

type InquiryV2Api interface {
	PostBalanceInquiry(ctx context.Context, params interface{}) (interface{}, error)
	PostExternalAccountInquiry(ctx context.Context, params interface{}) (interface{}, error)
	PostInternalAccountInquiry(ctx context.Context, params interface{}) (interface{}, error)
}
