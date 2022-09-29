package stats

import (
	"context"

	"github.com/gsmcwhirter/go-util/v10/errors"
	"github.com/gsmcwhirter/go-util/v10/telemetry"
)

// Known metric names
const (
	RawMessageCount               = "raw_request_ct"
	RawMessagesSentCount          = "raw_messages_sent_ct"
	InteractionResponsesCount     = "interaction_responses_ct"
	InteractionAutocompletesCount = "interaction_autocompletes_ct"
	InteractionDeferralsCount     = "interaction_deferrals_ct"
	MessagesPostedCount           = "messages_posted_ct"
	RawEventsCount                = "raw_events_ct"
	OpCodesCount                  = "opcode_events_ct"
)

// Known metric tag names
const (
	TagStatus    = "status"
	TagEventName = "event_name"
	TagOpCode    = "op_code"
)

// IncCounter increments a counter with the given value
func IncCounter(ctx context.Context, t *telemetry.Telemeter, pkg, name string, v int64, tags ...telemetry.KeyValue) error {
	counter, err := t.Meter(pkg).SyncInt64().Counter(name)
	if err != nil {
		return errors.Wrap(err, "could not create counter")
	}

	counter.Add(ctx, 1, tags...)
	return nil
}
