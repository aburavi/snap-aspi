package svc

import (
	"context"
)

type HistoryV2Api interface {
	GetTransactionHistoryDetail(ctx context.Context, params interface{}) (interface{}, error)
	GetTransactionHistoryList(ctx context.Context, params interface{}) (interface{}, error)
}
