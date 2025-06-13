package svc

import (
	"context"
)

type InquiryV2Api interface {
	PostBalanceInquiryV2(ctx context.Context, params interface{}) (interface{}, error)
	PostExternalAccountInquiryV2(ctx context.Context, params interface{}) (interface{}, error)
	PostInternalAccountInquiryV2(ctx context.Context, params interface{}) (interface{}, error)
}
