package errreport

import "context"

// Reporter is the interface that something can use to plug in and monitor panics and errors
type Reporter interface {
	AutoNotify(context.Context)
	Recover(context.Context)
	Notify(context.Context, error)
}

type NopReporter struct{}

func (r NopReporter) AutoNotify(ctx context.Context)        {}
func (r NopReporter) Recover(ctx context.Context)           { _ = recover() }
func (r NopReporter) Notify(ctx context.Context, err error) {}
