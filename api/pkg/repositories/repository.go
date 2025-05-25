package repositories

import (
	"time"

	"github.com/palantir/stacktrace"
)

// IndexParams parameters for indexing a database table
type IndexParams struct {
	Skip           int    `json:"skip"`
	SortBy         string `json:"sort"`
	SortDescending bool   `json:"sort_descending"`
	Query          string `json:"query"`
	Limit          int    `json:"take"`
}

// TimeSeriesData represents a time series data point
type TimeSeriesData struct {
	Timestamp time.Time `json:"timestamp"`
	Count     uint      `json:"count"`
}

const (
	// ErrCodeNotFound is thrown when an entity does not exist in storage
	ErrCodeNotFound = stacktrace.ErrorCode(1000)
)
