package eventstore

import (
	"reflect"
	"testing"
	"time"

	logrusTest "github.com/sirupsen/logrus/hooks/test"
)

// BEGIN TEST SETUP - Create an implementation of Event for use in tests
const EventTypeSpoke = "SPOKE"

type SpokeEvent struct {
	BaseEvent
	SpokeMessage string
}

func NewSpokeEvent(streamID string, eventNumber int, spokeMessage string) *SpokeEvent {
	return &SpokeEvent{BaseEvent{StreamID: streamID, EventNumber: eventNumber}, spokeMessage}
}

func (e *SpokeEvent) Init(data *EventData) {
	e.EventNumber = data.EventNumber
	e.StreamID = data.StreamID
	e.SpokeMessage = data.Data
	e.Timestamp = data.Timestamp
}

func (e *SpokeEvent) Data() (*EventData, error) {
	return &EventData{e.StreamID, e.EventNumber, EventTypeSpoke, e.SpokeMessage, time.Time{}}, nil
}

// END TEST SETUP

// Test that we create a factory and recieve an intialised event
func TestEventFactory(t *testing.T) {
	// Setup a factory
	log, _ := logrusTest.NewNullLogger()
	emptyEventCreators := EmptyEventCreatorsMap{
		EventTypeSpoke: func() Event { return &SpokeEvent{} },
	}
	factory := NewEventFactory(log, emptyEventCreators)

	timestamp := time.Now()

	// Create event data
	eventData := &EventData{"abcd", 10, EventTypeSpoke, "Hello, World!", timestamp}

	// Send event data to factory to process
	event := factory.CreateEvent(eventData)

	// Check it's an instance of SpokeEvent and it's has the correct data
	switch e := event.(type) {
	case *SpokeEvent:
		expectedStreamID := "abcd"
		streamID := e.GetStreamID()
		if expectedStreamID != streamID {
			t.Errorf("StreamID of created event invalid, got: %s, want: %s.", streamID, expectedStreamID)
		}

		expectedEventNo := 10
		eventNo := e.GetEventNumber()
		if expectedEventNo != eventNo {
			t.Errorf("EventNumber of created event invalid, got: %d, want: %d.", eventNo, expectedEventNo)
		}

		expectedMessage := "Hello, World!"
		if expectedMessage != e.SpokeMessage {
			t.Errorf("SpokeMessage of created event invalid, got: %s, want: %s.", e.SpokeMessage, expectedMessage)
		}

		if timestamp != e.Timestamp {
			t.Errorf("Timestamp of created event invalid, got: %s, want: %s.", e.Timestamp, timestamp)
		}

	default:
		t.Errorf("The created event is not the correct type, got: %v, want: '*event.SpokeEvent'", reflect.TypeOf(e))
	}
}

// Test that CreateEvent returns nil and logs a message when an unknown event is received
func TestEventFactoryUnkownEvent(t *testing.T) {
	// Setup a factory without any creators
	log, hook := logrusTest.NewNullLogger()
	emptyEventCreators := make(EmptyEventCreatorsMap)
	factory := NewEventFactory(log, emptyEventCreators)

	// Create an unknown event
	eventData := &EventData{"ab", 10, "NOT_REGISTERED_EVENT_TYPE", "DATA_HERE", time.Now()}
	event := factory.CreateEvent(eventData)
	if event != nil {
		t.Error("CreateEvent should return nil when the event type is not valid")
	}

	// Check at least one log message was sent
	if len(hook.Entries) <= 0 {
		t.Error("CreateEvent should log a message when the event type is not valid")
	}

	hook.Reset()
}
