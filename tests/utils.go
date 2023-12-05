package tests

import (
	"database/sql"
	"fmt"

	"log"
	"os"
	"os/exec"

	"github.com/go-faker/faker/v4/pkg/slice"
)

func GeneralSuiteSetup() {
	// turn off logging???

}

func TearDownSuite(dbPath string) {
	// remove the db + ignore errors
	os.Remove(dbPath)
	// also remove shm + wals but ignore errors
	os.Remove(fmt.Sprintf("%s-shm", dbPath))
	os.Remove(fmt.Sprintf("%s-wal", dbPath))
}

func SetupTestDatabase(dbPath string) error {
	p, _ := os.Getwd()
	fmt.Printf("%v", p)
	_, err := exec.Command("goose", "-dir", "../../migrations", "sqlite3", dbPath, "up").CombinedOutput()
	if err != nil {
		return err
	}
	return nil

}

func WipeDB(db *sql.DB) error {

	// like other project, wipe the dbs + go again
	if _, err := db.Exec(`PRAGMA foreign_keys = OFF;`); err != nil {
		return err
	}
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return err
	}
	defer rows.Close()
	var tableNames []string
	skipDelete := []string{"sqlite_sequence", "goose_db_version"}
	// create sth to hold the values
	for rows.Next() {
		var (
			name string
		)
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		if !slice.Contains(skipDelete, name) {
			tableNames = append(tableNames, name)
		}
	}

	for i := 0; i < len(tableNames); i++ {
		if _, err := db.Exec(fmt.Sprintf("delete from %s", tableNames[i])); err != nil {
			return err
		}

	}
	if _, err := db.Exec(`PRAGMA foreign_key_check;`); err != nil {
		return err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		return err
	}

	// ok
	return nil
}
