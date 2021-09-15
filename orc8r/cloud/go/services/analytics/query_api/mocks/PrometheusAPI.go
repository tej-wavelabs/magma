// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	model "github.com/prometheus/common/model"
	mock "github.com/stretchr/testify/mock"
)

// PrometheusAPI is an autogenerated mock type for the PrometheusAPI type
type PrometheusAPI struct {
	mock.Mock
}

// Query provides a mock function with given fields: ctx, query, ts
func (_m *PrometheusAPI) Query(ctx context.Context, query string, ts time.Time) (model.Value, v1.Warnings, error) {
	ret := _m.Called(ctx, query, ts)

	var r0 model.Value
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) model.Value); ok {
		r0 = rf(ctx, query, ts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Value)
		}
	}

	var r1 v1.Warnings
	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time) v1.Warnings); ok {
		r1 = rf(ctx, query, ts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(v1.Warnings)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, time.Time) error); ok {
		r2 = rf(ctx, query, ts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryRange provides a mock function with given fields: ctx, query, r
func (_m *PrometheusAPI) QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, v1.Warnings, error) {
	ret := _m.Called(ctx, query, r)

	var r0 model.Value
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.Range) model.Value); ok {
		r0 = rf(ctx, query, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Value)
		}
	}

	var r1 v1.Warnings
	if rf, ok := ret.Get(1).(func(context.Context, string, v1.Range) v1.Warnings); ok {
		r1 = rf(ctx, query, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(v1.Warnings)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, v1.Range) error); ok {
		r2 = rf(ctx, query, r)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
