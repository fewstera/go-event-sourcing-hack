package eventstore

import (
	"testing"
	"time"
)

// Test that when we use the embedable BaseEvent the provided methods actually work
func TestBaseEvent(t *testing.T) {
	event := &BaseEvent{"12345", 10, time.Now().UTC()}

	eventNo := event.GetEventNumber()
	expectedEventNo := 10
	if eventNo != expectedEventNo {
		t.Errorf("GetEventNumber returns incorrect value, got: %d, want: %d.", eventNo, expectedEventNo)
	}

	streamID := event.GetStreamID()
	expectedStreamID := "12345"
	if streamID != expectedStreamID {
		t.Errorf("GetStreamID returns incorrect ID, got: %s, want: %s.", streamID, expectedStreamID)
	}
}
