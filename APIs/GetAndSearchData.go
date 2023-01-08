package APIs

import (
	"YTSearchAPI/DTOs"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func SetContentTypeAsJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func GetVideos(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	limit := r.URL.Query().Get("limit")

	// checking if passed start and limit are valid
	startInt, err1 := strconv.Atoi(start)
	limitInt, err2 := strconv.Atoi(limit)
	if err1 != nil || startInt > 100 || startInt < 0 {
		start = "0"
	}
	if err2 != nil || limitInt > 100 || limitInt < 0 {
		limit = "10"
	}

	// defining connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// we open and close connection for every tranche of videos
	// opening a DB connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to verify supplied argument")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to establish a connection with DB")
	}
	sqlStatement := fmt.Sprintf("SELECT * FROM ytvideos ORDER BY publishdate DESC LIMIT %s OFFSET %s", limit, start)

	rows, err := db.Query(sqlStatement)
	
	if err != nil || rows.Err() != nil {
		fmt.Println("Failed to make an entry in DB", err)
		var errorObj DTOs.ErrorResp
		errorObj.ErrorMsg = err.Error()
		w.WriteHeader(http.StatusNotFound)
		erroObject, _ := json.Marshal(errorObj)
		w.Write(erroObject)
	}

	results := make([]DTOs.Video, 0)
	for rows.Next() {
		var result DTOs.Video
		err = rows.Scan(&result.VidId, &result.Title,
			&result.Description, &result.PublishDate,
			&result.Thumbnails)
		if err != nil {

		}
		results = append(results, result)
	}

	val, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(val)
}

func Search(w http.ResponseWriter, r *http.Request) {

}
