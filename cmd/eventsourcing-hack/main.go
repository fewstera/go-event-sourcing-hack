package main

import (
	"database/sql"
	"fmt"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
	"github.com/fewstera/go-event-sourcing-hack/pkg/server"
	"github.com/fewstera/go-event-sourcing-hack/pkg/user"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	db := initDb("root:password@tcp(db:3306)/events")
	projection := user.NewProjection()
	ef := eventstore.NewEventFactory(log, user.EmptyEventCreators())
	es := eventstore.NewDBEventStore(db, ef, log, []eventstore.Projection{projection})
	ch := user.NewCommandHandler(es, projection)

	s := server.NewServer(log, projection, ch)

	err := es.StartPolling()
	if err != nil {
		panic(fmt.Sprintf("Error starting polling: %v", err.Error()))
	}

	if err = s.Start(); err != nil {
		panic(err)
	}
}

func initDb(connectionString string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true", connectionString))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err.Error()))
	}
	return db
}
