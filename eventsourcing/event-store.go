package eventsourcing

import (
	"fmt"
	"time"

	"database/sql"
)

type EventStore struct {
	db              *sql.DB
	projection      *Projection
	currentPosition int
	fetchMoreStmt   *sql.Stmt
}

const EventStoreStreamCategory string = "USER"

func NewEventStore(db *sql.DB, projection *Projection) *EventStore {
	eventStore := new(EventStore)
	eventStore.db = db
	eventStore.projection = projection
	eventStore.currentPosition = 0

	eventStore.fetchMoreStmt = eventStore.prepareFetchMoreStatement()
	go eventStore.fetchMoreRecentEvents()

	return eventStore
}

func (eventStore *EventStore) SaveEvent(event Event) error {
	dataBytes, err := event.GetData()
	if err != nil {
		return err
	}

	_, err = eventStore.db.Exec(
		"INSERT INTO `event` (stream_category, stream_id, event_number, event_type, data) VALUES (?, ?, ?, ?, ?)",
		EventStoreStreamCategory, event.GetStreamId(), event.GetEventNumber(), event.GetEventType(), string(dataBytes),
	)

	if err != nil {
		return err
	}

	fmt.Println("Wrote event to store")
	return nil
}

func (eventStore *EventStore) prepareFetchMoreStatement() *sql.Stmt {
	sql := fmt.Sprintf("SELECT position, event_type, stream_id, event_number, data FROM event WHERE stream_category = '%v' AND position > ? LIMIT ?", EventStoreStreamCategory)
	stmt, err := eventStore.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	return stmt
}

func (eventStore *EventStore) fetchMoreRecentEvents() {
	limit := 5

	rows, err := eventStore.fetchMoreStmt.Query(eventStore.currentPosition, limit)
	if err != nil {
		eventStore.handleFetchMoreError(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var position int
		var eventType string
		var streamId string
		var eventNumber int
		var data []byte

		err := rows.Scan(&position, &eventType, &streamId, &eventNumber, &data)
		if err != nil {
			eventStore.handleFetchMoreError(err)
			return
		}

		event, err := eventStore.parseDbEvent(eventType, streamId, eventNumber, data)
		if err != nil {
			eventStore.handleFetchMoreError(err)
			return
		}

		if event != nil {
			err = eventStore.projection.Apply(event)
			if err != nil {
				eventStore.handleFetchMoreError(err)
				return
			}
		}

		eventStore.currentPosition = position
	}

	if err != nil {
		eventStore.handleFetchMoreError(err)
		return
	}

	eventStore.fetchMoreRecentEvents()
}

func (eventStore *EventStore) parseDbEvent(eventType string, streamId string, eventNumber int, data []byte) (Event, error) {
	var event Event
	switch eventType {
	case EventTypeUserCreated:
		event = new(UserCreatedEvent)
	case EventTypeUserGotOlder:
		event = new(UserGotOlderEvent)
	case EventTypeUserNameChanged:
		event = new(UserNameChangedEvent)
	default:
		fmt.Println("Unkown event %v\n", eventType)
		return nil, nil
	}

	fmt.Println("Event %v\n", eventType)

	err := event.InitFromDbEvent(streamId, eventNumber, data)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (eventStore *EventStore) handleFetchMoreError(err error) {
	retryAfter := 1000
	fmt.Printf("ERROR: %v\n", err)
	fmt.Printf("Retrying after %v milliseconds.\n", retryAfter)
	time.Sleep(time.Duration(retryAfter) * time.Millisecond)
	eventStore.fetchMoreRecentEvents()
}
