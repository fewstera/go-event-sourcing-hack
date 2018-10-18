package eventstore

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// A Projection is for deriving current state from the stream of events.
// As the event store receives events the Apply method will be called with each
// individual event in a sequential order.
type Projection interface {
	// Apply will be called when an event is received
	Apply(evt Event)
}

// DBEventStore is an event store that is driven by a SQL database and applies events to a list
// of projections.
//
// When starting up, the events will be read from position 0 in chunks of numberOfEventsToFetchPerQuery
// until it reaches the end of the event stream.
//
// The event store constantly polls the database asking for new events. When new events are received
// the event store will apply the events to all projections it knows about.
type DBEventStore struct {
	db              *sql.DB
	projections     []Projection
	eventFactory    *EventFactory
	log             *logrus.Logger
	stopPollingChan chan struct{}

	currentPosition int
	fetchMoreStmt   *sql.Stmt
}

// EventStore interface to define event store actions
type EventStore interface {
	SaveEvent(evt Event) error
}

const (
	numberOfEventsToFetchPerQuery int    = 100
	eventStoreStreamCategory      string = "INTEL"
)

// NewDBEventStore instantiates a new DBEventStore using the provided db and list of projections.
//
// After the instance has been initialised, the event polling process is started inside a new go routing.
func NewDBEventStore(db *sql.DB, eventFactory *EventFactory, log *logrus.Logger, projections []Projection) *DBEventStore {
	e := new(DBEventStore)
	e.db = db
	e.projections = projections
	e.eventFactory = eventFactory
	e.log = log

	e.fetchMoreStmt = e.prepareFetchMoreStatement()
	return e
}

// StartPolling makes the event store begin polling events from the database. It
// will continue until StopPollling is called.
func (e *DBEventStore) StartPolling() error {
	if e.stopPollingChan != nil {
		return errors.New("event store is already polling")
	}

	e.stopPollingChan = make(chan struct{})
	go e.fetchMoreRecentEvents()

	return nil
}

// StopPolling stops the event store from polling the database
func (e *DBEventStore) StopPolling() error {
	if e.stopPollingChan == nil {
		return errors.New("event store is not polling")
	}

	// Send a message on the stop channel
	e.stopPollingChan <- struct{}{}

	return nil
}

// SaveEvent saves the provided event to the database. If there is a issue persisting the event
// then an error is returned, otherwise nil is returned.
func (e *DBEventStore) SaveEvent(evt Event) error {
	eventData, err := evt.Data()
	if err != nil {
		return err
	}

	_, err = e.db.Exec(
		"INSERT INTO `event` (stream_category, stream_id, event_number, event_type, data) VALUES (?, ?, ?, ?, ?)",
		eventStoreStreamCategory, eventData.StreamID, eventData.EventNumber, eventData.EventType, eventData.Data,
	)
	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == 1062 {
				err = &EventNumberConflictError{evt.GetStreamID(), evt.GetEventNumber()}
			}
		}
		e.log.Warnf("Error saving event: %v", err.Error())
		return err
	}

	return nil
}

// prepareFetchMoreStatement prepares the SQL statement for fetching events. Preparing the statement
// should speed up querying.
func (e *DBEventStore) prepareFetchMoreStatement() *sql.Stmt {
	query := fmt.Sprintf("SELECT timestamp, position, event_type, stream_id, event_number, data FROM event WHERE stream_category = '%v' AND position > ? ORDER BY position LIMIT %v", eventStoreStreamCategory, numberOfEventsToFetchPerQuery)
	stmt, err := e.db.Prepare(query)
	if err != nil {
		panic(err)
	}
	return stmt
}

// fetchMoreRecentEvents is used to new events and apply them to the projections. This method is ran in an
// infinite loop, so should be called in a go routine.
func (e *DBEventStore) fetchMoreRecentEvents() {
	rows, err := e.fetchMoreStmt.Query(e.currentPosition)
	if err != nil {
		go e.handleFetchMoreError(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var timestamp time.Time
		var position int
		var eventType string
		var streamID string
		var eventNumber int
		var data string

		err := rows.Scan(&timestamp, &position, &eventType, &streamID, &eventNumber, &data)
		if err != nil {
			go e.handleFetchMoreError(err)
			return
		}

		evt := e.eventFactory.CreateEvent(&EventData{streamID, eventNumber, eventType, data, timestamp})
		if evt != nil {
			// Apply the event to each projection
			for _, projection := range e.projections {
				projection.Apply(evt)
			}
		}

		e.currentPosition = position
	}

	time.Sleep(time.Duration(200) * time.Millisecond)
	go e.continueIfNotStopped()
}

func (e *DBEventStore) continueIfNotStopped() {
	select {
	// If a message has been recieved on the stop polling channel stop.
	case <-e.stopPollingChan:
		e.log.Info("Stopping polling as requested.")
		e.stopPollingChan = nil
	// Otherwise fetch more
	default:
		go e.fetchMoreRecentEvents()
	}
}

// handleFetchMoreError handles errors when fetching more events. It logs the error
// and waits 1 second before trying again.
func (e *DBEventStore) handleFetchMoreError(err error) {
	retryAfter := 1000
	e.log.Warnf("Error in fetchMoreRecentEvents: %v\n", err)
	e.log.Infof("Retrying after %v milliseconds.\n", retryAfter)
	time.Sleep(time.Duration(retryAfter) * time.Millisecond)
	go e.continueIfNotStopped()
}
