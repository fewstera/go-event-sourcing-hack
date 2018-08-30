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
		event.GetStreamCategory(), event.GetStreamId(), event.GetEventNumber(), event.GetEventType(), string(dataBytes),
	)

	if err != nil {
		return err
	}

	fmt.Println("Wrote event to store")
	return nil
}

func (eventStore *EventStore) prepareFetchMoreStatement() *sql.Stmt {
	stmt, err := eventStore.db.Prepare("SELECT position, stream_category, event_number, event_type, data FROM event WHERE position > ? LIMIT ?")
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
		var streamCategory string
		var eventNumber int
		var eventType string
		var data string

		err := rows.Scan(&position, &streamCategory, &eventNumber, &eventType, &data)
		if err != nil {
			eventStore.handleFetchMoreError(err)
			return
		}

		event, err := eventStore.parseDbEvent(position, streamCategory, eventNumber, eventType, data)
		if err != nil {
			eventStore.handleFetchMoreError(err)
			return
		}
		err = eventStore.projection.Apply(event)
		if err != nil {
			eventStore.handleFetchMoreError(err)
			return
		}

		eventStore.currentPosition = position
	}

	if err != nil {
		eventStore.handleFetchMoreError(err)
		return
	}

	eventStore.fetchMoreRecentEvents()
}

func (eventStore *EventStore) parseDbEvent(position int, streamCategory string, eventNumber int, eventType string, data string) (Event, error) {
	var event Event = nil
	return event, nil
}

func (eventStore *EventStore) handleFetchMoreError(err error) {
	retryAfter := 1000
	fmt.Printf("ERROR: %v\n", err)
	fmt.Printf("Retrying after %v milliseconds.\n", retryAfter)
	time.Sleep(time.Duration(retryAfter) * time.Millisecond)
	eventStore.fetchMoreRecentEvents()
}
