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

const EventsBetweenSnapshots int = 2

func NewSnapshots(db *sql.DB) *Snapshots {
	s := new(Snapshots)
	s.db = db
	s.latestSnapshotPositionStmt = s.prepareGetLatestSnapshotPosition()

	return s
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
		"INSERT INTO `snapshot` (position, status) VALUES (?, 'CREATING_SNAPSHOT')", position,
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
		"UPDATE `snapshot` SET `status` = 'SNAPSHOT_COMPLETE', `location` = ? WHERE `position` = ?", snapshotFilePath, position,
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
