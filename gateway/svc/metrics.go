package svc

import (
)

type MetricsApi interface {
	RequestCount() error
	TotalRequests() error
	ResponseStatus() error
	HttpDuration() error
	RequestLatency() error
}
