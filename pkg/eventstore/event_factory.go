package eventstore

import (
	"github.com/sirupsen/logrus"
)

// EmptyEventCreatorsMap is a mapping of event types as strings to funcs that return
// empty instances of that event.
type EmptyEventCreatorsMap map[string]func() Event

// EventFactory is used to create instances of events from event data.
//
// You must instantiate the factory using the NewFactory func.
//
// You create new events using the CreateEvent() method.
type EventFactory struct {
	logger *logrus.Logger
	// emptyEventCreators is a mapping of event types (e.g. INTEL_CREATED) to funcs that return
	// new empty instances of those events (e.g. IntelCreatedEvent{}).
	emptyEventCreators EmptyEventCreatorsMap
}

// NewEventFactory returns a new instance of a EventFactory.
//
// A logrus Logger must be provided so that it can be logged when an unrecognised event is sent to
// the CreateEvent method.
//
// You also need to provide emptyEventCreators, a mapping of event type strings (e.g. INTEL_CREATED)
// to funcs that return a new empty instance of that event.
//
// Example:
//   emptyEventCreators := event.EmptyEventCreatorsMap{
//     "INTEL_CREATED": func() event.Event { return &IntelCreatedEvent{} },
//   }
//   ef := event.NewEventFactory(logrus, emptyEventCreators)
//
func NewEventFactory(logger *logrus.Logger, emptyEventCreators EmptyEventCreatorsMap) *EventFactory {
	f := new(EventFactory)
	f.logger = logger
	f.emptyEventCreators = emptyEventCreators
	return f
}

// CreateEvent creates an instance of an event. The type of instance that is created is
// determined by the EventType field of the data provided. For example, if the EventType
// is INTEL_CREATED then an instance of IntelCreatedEvent will be created.
//
// To determine how event types are mapped to event instances, the emptyEventCreators map
// that was provided when instantiating the factory is used. CreateEvent creates an empty
// event using the mapped to func and then calls the Init method on that instance to instantiate
// the instance with the event data.
func (f *EventFactory) CreateEvent(data *EventData) Event {
	emptyEventCreator, found := f.emptyEventCreators[data.EventType]
	if !found {
		f.logger.Infof("Event with unknown event type received '%v'", data.EventType)
		return nil
	}
	event := emptyEventCreator()
	event.Init(data)
	return event
}
