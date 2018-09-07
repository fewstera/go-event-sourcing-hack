package eventsourcing

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"os"
)

type Snapshots struct {
	db                         *sql.DB
	latestSnapshotPosition     int
	latestSnapshotPositionStmt *sql.Stmt
}

type SnapshotState struct {
	Position   int
	Projection *Projection
}

const EventsBetweenSnapshots int = 10
const SnapshotStatusCreating string = "CREATING_SNAPSHOT"
const SnapshotStatusComplete string = "SNAPSHOT_COMPLETE"

func NewSnapshots(db *sql.DB) *Snapshots {
	s := new(Snapshots)
	s.db = db
	s.latestSnapshotPositionStmt = s.prepareGetLatestSnapshotPosition()

	return s
}

func (s *Snapshots) GetStateFromLatestSnapshot() *SnapshotState {
	var position int
	var location string

	err := s.db.QueryRow("SELECT position, location FROM snapshot WHERE `status` = ?  ORDER BY `position` DESC LIMIT 1", SnapshotStatusComplete).Scan(&position, &location)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("Error fetching latest snapshot: %v\n", err)
		} else {
			fmt.Println("No snapshot found, starting from scratch.")
		}
		return &SnapshotState{0, nil}
	}

	snapshotFile, err := os.Open(location)
	if err != nil {
		fmt.Printf("Error reading snapshot file: %v\n", err)
		return &SnapshotState{0, nil}
	}

	var projection Projection
	dec := gob.NewDecoder(snapshotFile)
	err = dec.Decode(&projection)
	if err != nil {
		fmt.Printf("Error loading snapshot into projection: %v\n", err)
		return &SnapshotState{0, nil}
	}

	return &SnapshotState{position, &projection}
}

func (s *Snapshots) prepareGetLatestSnapshotPosition() *sql.Stmt {
	stmt, err := s.db.Prepare("SELECT position FROM snapshot ORDER BY position DESC LIMIT 1")
	if err != nil {
		panic(err)
	}
	return stmt
}

func (s *Snapshots) takeNewSnapshotIfNeeded(projection *Projection, position int) {
	// Only check if there is any new snapshots if our current position is greater than last snapshot + distance between snapshots
	if position >= s.latestSnapshotPosition+EventsBetweenSnapshots {
		fmt.Println("Querying the current snapshot status")
		err := s.latestSnapshotPositionStmt.QueryRow().Scan(&s.latestSnapshotPosition)
		if err != nil {
			if err == sql.ErrNoRows {
				s.latestSnapshotPosition = 0
			} else {
				fmt.Printf("Error fetching latest snapshot position: %v", err)
				return
			}
		}

		// If we still need to take a snapshot after pulling the latest snapshot version then take one.
		if position >= s.latestSnapshotPosition+EventsBetweenSnapshots {
			s.createNewSnapshot(projection, position)
		}
	}
}

func (s *Snapshots) createNewSnapshot(projection *Projection, position int) {
	_, err := s.db.Exec(
		"INSERT INTO `snapshot` (position, status) VALUES (?, ?)", position, SnapshotStatusCreating,
	)

	if err != nil {
		fmt.Printf("Error SQL: %v\n", err)
		return
	}

	snapshotFilePath, err := s.writeSnapshotToDisk(projection, position)
	if err != nil {
		fmt.Printf("Error writing snapshot: %v\n", err)
		return
	}

	_, err = s.db.Exec(
		"UPDATE `snapshot` SET `status` = ?, `location` = ? WHERE `position` = ?", SnapshotStatusComplete, snapshotFilePath, position,
	)
	if err != nil {
		fmt.Printf("Error setting snapshot success status in db: %v\n", err)
		return
	}

	s.latestSnapshotPosition = position
	fmt.Printf("Wrote snapshot: %v\n", snapshotFilePath)
}

func (s *Snapshots) writeSnapshotToDisk(projection *Projection, position int) (string, error) {
	snapshotFolder := "snapshots"
	snapshotFilePath := fmt.Sprintf("%s/%d.bin", snapshotFolder, position)

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
	err = enc.Encode(projection)
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
