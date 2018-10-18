// Package eventstore provides an event store that is driven by a SQL database.
//
// The event store pulls every event, starting from the start of the event
// stream and applies them to the provided projections.
package eventstore
