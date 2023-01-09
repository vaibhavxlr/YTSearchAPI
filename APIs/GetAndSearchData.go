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

func getError(err string) []byte {
	var errorObj DTOs.ErrorResp
	errorObj.ErrorMsg = err
	errorObject, _ := json.Marshal(errorObj)
	return errorObject
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

	// opening a DB connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		msg := "Failed to verify supplied argument"
		fmt.Println(msg)
		w.WriteHeader(404)
		w.Write(getError(msg))
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		msg := "Failed to establish a connection with DB"
		fmt.Println(msg)
		w.WriteHeader(404)
		w.Write(getError(msg))
		return
	}
	sqlStatement := fmt.Sprintf("SELECT * FROM ytvideos ORDER BY publishdate DESC LIMIT %s OFFSET %s", limit, start)

	rows, err := db.Query(sqlStatement)

	if err != nil || rows.Err() != nil {
		msg := fmt.Sprint("Failed to fetch data from DB", err)
		fmt.Println(msg)
		w.WriteHeader(404)
		w.Write(getError(msg))
		return
	}

	results := make([]DTOs.Video, 0)
	for rows.Next() {
		var result DTOs.Video
		err = rows.Scan(&result.VidId, &result.Title,
			&result.Description, &result.PublishDate,
			&result.Thumbnails)
		if err != nil {
			msg := fmt.Sprint("Failed to map data from DB", err)
			fmt.Println(msg)
			w.WriteHeader(404)
			w.Write(getError(msg))
			return
		}
		results = append(results, result)
	}

	val, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(val)
}

func Search(w http.ResponseWriter, r *http.Request) {
	var reqBody DTOs.SearchReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		msg := "Request is not sent in proper format"
		fmt.Println(msg)
		w.WriteHeader(404)
		w.Write(getError(msg))
		return
	}

	// defining connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// opening a DB connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		msg := "Failed to verify supplied argument"
		fmt.Println(msg)
		w.WriteHeader(404)
		w.Write(getError(msg))
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		msg := "Failed to establish a connection with DB"
		fmt.Println(msg)
		w.WriteHeader(404)
		w.Write(getError(msg))
		return
	}
	sqlStatement := fmt.Sprintf(`SELECT * FROM ytvideos WHERE title LIKE %s AND description LIKE %s`, "'%"+reqBody.Title+"%'", "'%"+reqBody.Description+"%'")
	if reqBody.ExactMatch {
		sqlStatement = fmt.Sprintf(`SELECT * FROM ytvideos WHERE title LIKE %s AND description LIKE %s`, "'"+reqBody.Title+"'", "'"+reqBody.Description+"'")
	}

	rows, err := db.Query(sqlStatement)

	if err != nil || rows.Err() != nil {
		msg := fmt.Sprint("Failed to fetch data from DB", err)
		fmt.Println(msg)
		w.WriteHeader(404)
		w.Write(getError(msg))
		return
	}

	results := make([]DTOs.Video, 0)
	for rows.Next() {
		var result DTOs.Video
		err = rows.Scan(&result.VidId, &result.Title,
			&result.Description, &result.PublishDate,
			&result.Thumbnails)
		if err != nil {
			msg := fmt.Sprint("Failed to map data from DB", err)
			fmt.Println(msg)
			w.WriteHeader(404)
			w.Write(getError(msg))
			return
		}
		results = append(results, result)
	}
	if len(results) == 0 {
		msg := "No matching response for the req criteria"
		fmt.Println(msg)
		w.WriteHeader(200)
		w.Write(getError(msg))
		return
	}
	val, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(val)

}
