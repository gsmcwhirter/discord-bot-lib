package stats

import (
	"github.com/gsmcwhirter/go-util/v8/errors"
	"github.com/gsmcwhirter/go-util/v8/telemetry"
)

var (
	RawMessageCount      = telemetry.Int64("raw_request_ct", "Raw request count", "1")
	RawMessagesSentCount = telemetry.Int64("raw_messages_sent_ct", "Raw messages sent count", "1")
	MessagesPostedCount  = telemetry.Int64("messages_posted_ct", "Messages posted count", "1")
	RawEventsCount       = telemetry.Int64("raw_events_ct", "Raw events count", "1")
	OpCodesCount         = telemetry.Int64("opcode_events_ct", "OpCode events count", "1")
)

var (
	TagStatus, _    = telemetry.NewTagKey("status")
	TagEventName, _ = telemetry.NewTagKey("event_name")
	TagOpCode, _    = telemetry.NewTagKey("op_code")
)

var (
	RawMessageCountView = &telemetry.View{
		Name:        "raw_requests",
		TagKeys:     []telemetry.TagKey{},
		Measure:     RawMessageCount,
		Description: "The number of raw messages received",
		Aggregation: telemetry.CountView(),
	}

	RawMessagesSentCountView = &telemetry.View{
		Name:        "raw_messages_sent",
		TagKeys:     []telemetry.TagKey{},
		Measure:     RawMessagesSentCount,
		Description: "The number of raw messages sent",
		Aggregation: telemetry.CountView(),
	}

	MessagesPostedCountView = &telemetry.View{
		Name: "messages_posted",
		TagKeys: []telemetry.TagKey{
			TagStatus,
		},
		Measure:     MessagesPostedCount,
		Description: "The number of messages posted to discord",
		Aggregation: telemetry.CountView(),
	}

	RawEventsCountView = &telemetry.View{
		Name: "raw_events",
		TagKeys: []telemetry.TagKey{
			TagEventName,
		},
		Measure:     RawEventsCount,
		Description: "The number of events processed by the messagehandler",
		Aggregation: telemetry.CountView(),
	}

	OpCodesCountView = &telemetry.View{
		Name: "opcode_events",
		TagKeys: []telemetry.TagKey{
			TagOpCode,
		},
		Measure:     OpCodesCount,
		Description: "The number of opcode events processed by the messagehandler",
		Aggregation: telemetry.CountView(),
	}
)

func Register() error {
	if err := telemetry.RegisterView(RawMessageCountView); err != nil {
		return errors.Wrap(err, "could not register RawMessageCountView")
	}

	if err := telemetry.RegisterView(RawMessagesSentCountView); err != nil {
		return errors.Wrap(err, "could not register RawMessagesSentCountView")
	}

	if err := telemetry.RegisterView(MessagesPostedCountView); err != nil {
		return errors.Wrap(err, "could not register MessagesPostedCountView")
	}

	if err := telemetry.RegisterView(RawEventsCountView); err != nil {
		return errors.Wrap(err, "could not register RawEventsCountView")
	}

	if err := telemetry.RegisterView(OpCodesCountView); err != nil {
		return errors.Wrap(err, "could not register OpCodesCountView")
	}

	return nil
}
