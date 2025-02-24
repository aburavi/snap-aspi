package svc

import (
	"context"
)

type HistoryApi interface {
	GetTransactionHistoryDetail(ctx context.Context, params interface{}) (interface{}, error)
	GetTransactionHistoryList(ctx context.Context, params interface{}) (interface{}, error)
}
