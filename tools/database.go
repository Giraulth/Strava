package tools

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbName   = "strava"
)

func ConnectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Error connecting to postgres database: %v\n", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error connecting to postgres database: %v\n", err)
		return nil
	}

	return db
}

// GenericInsert inserts data into the specified table with the given columns and values.
func GenericInsert(db *sql.DB, tableName string, tableColumns []string, values ...interface{}) error {
	// Construct the SQL query string with placeholders for each column
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(tableColumns, ", "),
		placeholders(len(values)))

	// Perform the insert operation
	_, err := db.Exec(query, values...)
	return err
}

// Helper function to create the placeholders for the query
func placeholders(n int) string {
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(placeholders, ", ")
}

func GenericCount(db *sql.DB, table string) int {
	query := fmt.Sprintf("SELECT count(*) FROM %s", table)

	row := db.QueryRow(query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		fmt.Println("Error fetching count:", err)
		return 0
	}
	return count
}

func SelectId(db *sql.DB) []int64 {
	activityList := []int64{}
	rows, err := db.Query(
		"SELECT id FROM activity WHERE id not in (select activity_id from kudos);")
	if err != nil {
		fmt.Println("Error executing SELECT query:", err)
		return activityList
	}
	defer rows.Close()

	// Iterate over the rows and scan the int64 values into the array
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return activityList
		}
		activityList = append(activityList, id)
	}
	return activityList
}
