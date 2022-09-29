package errreport

import "context"

// Reporter is the interface that something can use to plug in and monitor panics and errors
type Reporter interface {
	AutoNotify(context.Context)
	Recover(context.Context)
	Notify(context.Context, error)
}

// NopReporter is a Reporter that does nothing (useful for testing)
type NopReporter struct{}

// AutoNotify is a placeholder that does nothing
func (r NopReporter) AutoNotify(ctx context.Context) {}

// Recover is a placeholder that just does a recover
func (r NopReporter) Recover(ctx context.Context) { _ = recover() }

// Notify is a placeholder that does nothing
func (r NopReporter) Notify(ctx context.Context, err error) {}
