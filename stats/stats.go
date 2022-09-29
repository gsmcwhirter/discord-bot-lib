package stats

import (
	"github.com/gsmcwhirter/go-util/v8/errors"
	"github.com/gsmcwhirter/go-util/v8/telemetry"
)

var (
	// RawMessageCount is the count of all messages received for processing
	RawMessageCount = telemetry.Int64("raw_request_ct", "Raw request count", "1")
	// RawMessagesSentCount is the count of all messages sent by the bot
	RawMessagesSentCount = telemetry.Int64("raw_messages_sent_ct", "Raw messages sent count", "1")
	// MessagesPostedCount is the count of all messages posted
	MessagesPostedCount = telemetry.Int64("messages_posted_ct", "Messages posted count", "1")
	// InteractionResponsesCount is the count of all interaction responses sent
	InteractionResponsesCount = telemetry.Int64("interaction_responses_ct", "Interaction responses count", "1")
	// InteractionAutocompletesCount is the count of all autocomplete responses sent
	InteractionAutocompletesCount = telemetry.Int64("interaction_autocompletes_ct", "Interaction autocompletes count", "1")
	// InteractionDeferralsCount is the count of all interaction deferrals sent
	InteractionDeferralsCount = telemetry.Int64("interaction_deferrals_ct", "Interaction deferrals count", "1")
	// RawEventsCount is the count of all events handled
	RawEventsCount = telemetry.Int64("raw_events_ct", "Raw events count", "1")
	// OpCodesCount is the count of each op-code handled
	OpCodesCount = telemetry.Int64("opcode_events_ct", "OpCode events count", "1")
)

var (
	// TagStatus is the tag for a response status
	TagStatus, _ = telemetry.NewTagKey("status")
	// TagEventName is the tag for an event name
	TagEventName, _ = telemetry.NewTagKey("event_name")
	// TagOpCode is the tag for the op code
	TagOpCode, _ = telemetry.NewTagKey("op_code")
)

var (
	// RawMessageCountView is the view for RaeMessageCount
	RawMessageCountView = &telemetry.View{
		Name:        "raw_requests",
		TagKeys:     []telemetry.TagKey{},
		Measure:     RawMessageCount,
		Description: "The number of raw messages received",
		Aggregation: telemetry.CountView(),
	}

	// RawMessagesSentCountView is the view for RawMessagesSentCount
	RawMessagesSentCountView = &telemetry.View{
		Name:        "raw_messages_sent",
		TagKeys:     []telemetry.TagKey{},
		Measure:     RawMessagesSentCount,
		Description: "The number of raw messages sent",
		Aggregation: telemetry.CountView(),
	}

	// MessagesPostedCountView is the view for MessagesPostedCount
	MessagesPostedCountView = &telemetry.View{
		Name: "messages_posted",
		TagKeys: []telemetry.TagKey{
			TagStatus,
		},
		Measure:     MessagesPostedCount,
		Description: "The number of messages posted to discord",
		Aggregation: telemetry.CountView(),
	}

	// InteractionResponsesCountView is the view for InteractionResponsesCount
	InteractionResponsesCountView = &telemetry.View{
		Name:        "interaction_responses",
		TagKeys:     []telemetry.TagKey{},
		Measure:     InteractionResponsesCount,
		Description: "The number of interaction responses sent",
		Aggregation: telemetry.CountView(),
	}

	// InteractionAutocompletesCountView is the view for InteractionAutocompletesCount
	InteractionAutocompletesCountView = &telemetry.View{
		Name:        "interaction_autocompletes",
		TagKeys:     []telemetry.TagKey{},
		Measure:     InteractionAutocompletesCount,
		Description: "The number of interaction autocompletes sent",
		Aggregation: telemetry.CountView(),
	}

	// InteractionDeferralsCountView is the view for InteractionDeferralsCount
	InteractionDeferralsCountView = &telemetry.View{
		Name:        "interaction_deferrals",
		TagKeys:     []telemetry.TagKey{},
		Measure:     InteractionDeferralsCount,
		Description: "The number of interaction deferrals sent",
		Aggregation: telemetry.CountView(),
	}

	// RawEventsCountView is the view for RawEventsCount
	RawEventsCountView = &telemetry.View{
		Name: "raw_events",
		TagKeys: []telemetry.TagKey{
			TagEventName,
		},
		Measure:     RawEventsCount,
		Description: "The number of events processed by the messagehandler",
		Aggregation: telemetry.CountView(),
	}

	// OpCodesCountView is the view for OpCodesCount
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

// Register registers the various metric views with the telemetry recorder
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

	if err := telemetry.RegisterView(InteractionResponsesCountView); err != nil {
		return errors.Wrap(err, "could not register InteractionResponsesCountView")
	}

	if err := telemetry.RegisterView(InteractionAutocompletesCountView); err != nil {
		return errors.Wrap(err, "could not register InteractionAutocompletesCountView")
	}

	if err := telemetry.RegisterView(InteractionDeferralsCountView); err != nil {
		return errors.Wrap(err, "could not register InteractionDeferralsCountView")
	}

	return nil
}
