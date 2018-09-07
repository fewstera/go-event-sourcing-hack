package eventsourcing

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

type EventStore struct {
	db                          *sql.DB
	projection                  *Projection
	currentPosition             int
	numberOfPollsSinceLastEvent int
	latestSnapshotPosition      int
	fetchMoreStmt               *sql.Stmt
	latestSnapshotPositionStmt  *sql.Stmt
}

const EventStoreStreamCategory string = "USER"
const EventsBetweenSnapshots int = 1

func NewEventStore(db *sql.DB, projection *Projection) *EventStore {
	eventStore := new(EventStore)
	eventStore.db = db
	eventStore.projection = projection

	eventStore.fetchMoreStmt = eventStore.prepareFetchMoreStatement()
	eventStore.latestSnapshotPositionStmt = eventStore.prepareGetLatestSnapshotPosition()

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

func (eventStore *EventStore) prepareGetLatestSnapshotPosition() *sql.Stmt {
	stmt, err := eventStore.db.Prepare("SELECT position FROM snapshot ORDER BY position DESC LIMIT 1")
	if err != nil {
		panic(err)
	}
	return stmt
}

func (eventStore *EventStore) fetchMoreRecentEvents() {
	limit := 5
	startPosition := eventStore.currentPosition

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
			fmt.Println("Got new event")
			err = eventStore.projection.Apply(event)
			if err != nil {
				eventStore.handleFetchMoreError(err)
				return
			}
		}

		eventStore.currentPosition = position
		eventStore.numberOfPollsSinceLastEvent = 0
	}

	if err != nil {
		eventStore.handleFetchMoreError(err)
		return
	}

	if startPosition == eventStore.currentPosition {
		eventStore.numberOfPollsSinceLastEvent++
	}

	if eventStore.numberOfPollsSinceLastEvent > 1 {
		eventStore.checkIfNewSnapshotNeeded()
	}

	time.Sleep(time.Duration(200) * time.Millisecond)
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
		fmt.Printf("Unkown event %v\n", eventType)
		return nil, nil
	}

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

func (eventStore *EventStore) checkIfNewSnapshotNeeded() {
	if eventStore.currentPosition >= eventStore.latestSnapshotPosition+EventsBetweenSnapshots {
		fmt.Println("Querying the current snapshot status")
		err := eventStore.latestSnapshotPositionStmt.QueryRow().Scan(&eventStore.latestSnapshotPosition)
		if err != nil {
			if err == sql.ErrNoRows {
				eventStore.latestSnapshotPosition = 0
			} else {
				fmt.Printf("Error fetching latest snapshot position: %v", err)
				return
			}
		}

		if eventStore.currentPosition >= eventStore.latestSnapshotPosition+EventsBetweenSnapshots {
			eventStore.createNewSnapshot()
		}
	}
}

func (eventStore *EventStore) createNewSnapshot() {
	_, err := eventStore.db.Exec(
		"INSERT INTO `snapshot` (position, status) VALUES (?, 'CREATING_SNAPSHOT')", eventStore.currentPosition,
	)

	if err != nil {
		fmt.Printf("Error SQL: %v\n", err)
	}

	eventStore.latestSnapshotPosition = eventStore.currentPosition
	snapshotFilePath, err := eventStore.writeSnapshotToDisk()
	if err != nil {
		fmt.Printf("Error writing snapshot: %v\n", err)
	}

	_, err = eventStore.db.Exec(
		"UPDATE `snapshot` SET `status` = 'SNAPSHOT_COMPLETE', `location` = ? WHERE `position` = ?", snapshotFilePath, eventStore.currentPosition,
	)
	if err != nil {
		fmt.Printf("Error setting snapshot success status in db: %v\n", err)
	}

	fmt.Printf("Wrote snapshot: %v\n", snapshotFilePath)
}

func (eventStore *EventStore) writeSnapshotToDisk() (string, error) {
	snapshotFolder := "snapshots"
	snapshotFilePath := fmt.Sprintf("%s/%d.bin", snapshotFolder, eventStore.currentPosition)

	// mkdir -p snapshots/
	err := os.MkdirAll(snapshotFolder, 0755)
	if err != nil {
		fmt.Printf("Error verifying or creating snapshot folder: %v\n", err)
		return "", err
	}

	snapshotFile, err := os.Create(snapshotFilePath)
	if err != nil {
		return "", err
	}

	defer snapshotFile.Close()
	enc := gob.NewEncoder(snapshotFile)
	err = enc.Encode(eventStore.projection)
	if err != nil {
		fmt.Printf("encode error: %v\n", err)
		rmErr := os.Remove(snapshotFilePath)
		if rmErr != nil {
			fmt.Printf("Failed to delete failed snapshot: %v\n", rmErr)
		}
		return "", err
	}

	return snapshotFilePath, nil
}
