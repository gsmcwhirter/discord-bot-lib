package stats

import (
	"github.com/gsmcwhirter/go-util/v4/census"
	"github.com/gsmcwhirter/go-util/v4/errors"
)

var (
	RawMessageCount      = census.Int64("raw_request_ct", "Raw request count", "1")
	RawMessagesSentCount = census.Int64("raw_messages_sent_ct", "Raw messages sent count", "1")
	MessagesPostedCount  = census.Int64("messages_posted_ct", "Messages posted count", "1")
	RawEventsCount       = census.Int64("raw_events_ct", "Raw events count", "1")
	OpCodesCount         = census.Int64("opcode_events_ct", "OpCode events count", "1")
)

var (
	TagStatus, _    = census.NewTagKey("status")
	TagEventName, _ = census.NewTagKey("event_name")
	TagOpCode, _    = census.NewTagKey("op_code")
)

var (
	RawMessageCountView = &census.View{
		Name:        "raw_requests",
		TagKeys:     []census.TagKey{},
		Measure:     RawMessageCount,
		Description: "The number of raw messages received",
		Aggregation: census.CountView(),
	}

	RawMessagesSentCountView = &census.View{
		Name:        "raw_messages_sent",
		TagKeys:     []census.TagKey{},
		Measure:     RawMessagesSentCount,
		Description: "The number of raw messages sent",
		Aggregation: census.CountView(),
	}

	MessagesPostedCountView = &census.View{
		Name: "messages_posted",
		TagKeys: []census.TagKey{
			TagStatus,
		},
		Measure:     MessagesPostedCount,
		Description: "The number of messages posted to discord",
		Aggregation: census.CountView(),
	}

	RawEventsCountView = &census.View{
		Name: "raw_events",
		TagKeys: []census.TagKey{
			TagEventName,
		},
		Measure:     RawEventsCount,
		Description: "The number of events processed by the messagehandler",
		Aggregation: census.CountView(),
	}

	OpCodesCountView = &census.View{
		Name: "opcode_events",
		TagKeys: []census.TagKey{
			TagOpCode,
		},
		Measure:     OpCodesCount,
		Description: "The number of opcode events processed by the messagehandler",
		Aggregation: census.CountView(),
	}
)

func Register() error {
	if err := census.RegisterView(RawMessageCountView); err != nil {
		return errors.Wrap(err, "could not register RawMessageCountView")
	}

	if err := census.RegisterView(RawMessagesSentCountView); err != nil {
		return errors.Wrap(err, "could not register RawMessagesSentCountView")
	}

	if err := census.RegisterView(MessagesPostedCountView); err != nil {
		return errors.Wrap(err, "could not register MessagesPostedCountView")
	}

	if err := census.RegisterView(RawEventsCountView); err != nil {
		return errors.Wrap(err, "could not register RawEventsCountView")
	}

	if err := census.RegisterView(OpCodesCountView); err != nil {
		return errors.Wrap(err, "could not register OpCodesCountView")
	}

	return nil
}
