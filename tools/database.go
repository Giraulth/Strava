package tools

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"sort"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbName   = "strava"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to postgres database: %v\n", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error connecting to postgres database: %v\n", err)
	}

	return db, nil
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

func genericSelect(db *sql.DB, query string) *sql.Rows {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil
	}
	return rows
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

func SelectKudosRanking(db *sql.DB) ([]map[string]interface{}, error) {
    ranking := []map[string]interface{}{}
	kudosRows := genericSelect(db, "SELECT username, count(*) FROM kudos WHERE username != '' GROUP BY username ORDER BY count(*) DESC;")
	activitySeenRows := genericSelect(db,
		"SELECT subquery.username, COUNT(*) AS activity_seen "+
			"FROM Activity "+
			"JOIN ( "+
			"    SELECT username, MIN(activity_id) AS min_activity_id "+
			"    FROM kudos "+
			"    WHERE username != ''"+
			"    GROUP BY username "+
			") AS subquery "+
			"ON Activity.id > subquery.min_activity_id "+
			"GROUP BY subquery.username "+
			"ORDER BY COUNT(*) DESC;")

	defer kudosRows.Close()
	defer activitySeenRows.Close()

	// Get number of kudos per user
	kudosData := make(map[string]int)
	for kudosRows.Next() {
		var username string
		var kudosCount int
		if err := kudosRows.Scan(&username, &kudosCount); err != nil {
			return ranking, err
		}
		kudosData[username] = kudosCount
	}

    for activitySeenRows.Next() {
        var username string
        var activitySeen int
        if err := activitySeenRows.Scan(&username, &activitySeen); err != nil {
            return ranking, err
        }

		ratioFloat := float64(kudosData[username]) / float64(activitySeen) * 100

        entry := map[string]interface{}{
            "username":     username,
            "kudos_count":  kudosData[username],
            "activity_seen_count": activitySeen,
			"ratio": ratioFloat,
        }

        ranking = append(ranking, entry)
    }

    if err := activitySeenRows.Err(); err != nil {
        return ranking, err
    }

    // Sort based on ratio
    sort.Slice(ranking, func(i, j int) bool {
        ratio1 := ranking[i]["ratio"].(float64)
        ratio2 := ranking[j]["ratio"].(float64)
        return ratio1 > ratio2
    })

    for _, entry := range ranking {
        entry["ratio"] = fmt.Sprintf("%.2f%%", entry["ratio"].(float64))
    }

	return ranking, nil
}
func SelectId(db *sql.DB) []int64 {
	activityList := []int64{}
	rows := genericSelect(db, "SELECT id FROM activity WHERE id not in (select activity_id from kudos);")
	defer rows.Close()
	// Iterate over the rows and scan the int64 values into the array
	for rows.Next() {
		var id int64
		fmt.Printf("%d", id)
		if err := rows.Scan(&id); err != nil {
			fmt.Println("Error scanning row:", err)
			return activityList
		}
		activityList = append(activityList, id)
	}
	return activityList
}
