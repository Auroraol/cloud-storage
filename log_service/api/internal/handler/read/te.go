package read

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

package main

import (
"database/sql"
"encoding/json"
"fmt"
"log"
"net/http"
"strconv"
"strings"
"time"

"github.com/gorilla/mux"
_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var logger = log.New(os.Stdout, "logfile-service: ", log.LstdFlags)

func init() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		logger.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		logger.Fatal(err)
	}
}

type Logfile struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Host          string    `json:"host"`
	Path          string    `json:"path"`
	CreateTime    time.Time `json:"create_time"`
	Comment       string    `json:"comment"`
	MonitorChoice int       `json:"monitor_choice"`
}

type Handler struct{}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pk, _ := strconv.Atoi(vars["id"])
	responseData := h.query(pk)
	h.writeResponse(w, responseData)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	responseData := h.add(r)
	h.writeResponse(w, responseData)
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pk, _ := strconv.Atoi(vars["id"])
	responseData := h.update(r, pk)
	h.writeResponse(w, responseData)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pk, _ := strconv.Atoi(vars["id"])
	responseData := h.del(pk)
	h.writeResponse(w, responseData)
}

func (h *Handler) query(pk int) map[string]interface{} {
	rows, err := db.Query("SELECT id, name, host, path, create_time, comment, monitor_choice FROM logfile WHERE id=?", pk)
	if err != nil {
		logger.Println(err)
		return map[string]interface{}{"code": 500, "msg": "Query failed"}
	}
	defer rows.Close()

	var logfiles []Logfile
	for rows.Next() {
		var logfile Logfile
		if err := rows.Scan(&logfile.ID, &logfile.Name, &logfile.Host, &logfile.Path, &logfile.CreateTime, &logfile.Comment, &logfile.MonitorChoice); err != nil {
			logger.Println(err)
			return map[string]interface{}{"code": 500, "msg": "Query failed"}
		}
		logfiles = append(logfiles, logfile)
	}

	return map[string]interface{}{"code": 200, "msg": "Query Successful", "data": logfiles}
}

func (h *Handler) add(r *http.Request) map[string]interface{} {
	name := r.FormValue("name")
	path := r.FormValue("path")
	comment := r.FormValue("comment")
	host := r.FormValue("host")
	monitorChoice, _ := strconv.Atoi(r.FormValue("monitor_choice"))

	if err := h.validateArguments(name, path, comment, host, monitorChoice, 0); err != nil {
		return err
	}

	createTime := time.Now().Format("2006-01-02 15:04:05")
	result, err := db.Exec("INSERT INTO logfile (name, host, path, create_time, comment, monitor_choice) VALUES (?, ?, ?, ?, ?, ?)", name, host, path, createTime, comment, monitorChoice)
	if err != nil {
		logger.Println(err)
		return map[string]interface{}{"code": 500, "msg": "Add failed"}
	}

	lastInsertID, _ := result.LastInsertId()
	return map[string]interface{}{"code": 200, "msg": "Add successful", "data": map[string]int64{"id": lastInsertID}}
}

func (h *Handler) update(r *http.Request, pk int) map[string]interface{} {
	name := r.FormValue("name")
	path := r.FormValue("path")
	comment := r.FormValue("comment")
	host := r.FormValue("host")
	monitorChoice, _ := strconv.Atoi(r.FormValue("monitor_choice"))

	if err := h.validateArguments(name, path, comment, host, monitorChoice, pk); err != nil {
		return err
	}

	_, err := db.Exec("UPDATE logfile SET name=?, host=?, path=?, comment=?, monitor_choice=? WHERE id=?", name, host, path, comment, monitorChoice, pk)
	if err != nil {
		logger.Println(err)
		return map[string]interface{}{"code": 500, "msg": "Update failed"}
	}

	return map[string]interface{}{"code": 200, "msg": "Update successful", "data": map[string]int{"id": pk}}
}

func (h *Handler) del(pk int) map[string]interface{} {
	_, err := db.Exec("DELETE FROM logfile WHERE id=?", pk)
	if err != nil {
		logger.Println(err)
		return map[string]interface{}{"code": 500, "msg": "Delete failed"}
	}

	return map[string]interface{}{"code": 200, "msg": "Delete successful"}
}

func (h *Handler) validateArguments(name, path, comment, host string, monitorChoice, pk int) map[string]interface{} {
	if path == "" {
		return map[string]interface{}{"code": 400, "msg": "Bad POST data", "error": map[string]string{"path": "Required"}}
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM logfile WHERE name=? AND id!=?", name, pk).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		logger.Println(err)
		return map[string]interface{}{"code": 500, "msg": "Validation failed"}
	}
	if count > 0 {
		return map[string]interface{}{"code": 400, "msg": "Bad POST data", "error": map[string]string{"path": "Already existed"}}
	}

	for _, field := range []struct {
		value string
		key   string
	}{
		{name, "name"},
		{host, "host"},
		{comment, "comment"},
	} {
		if field.value == "" {
			return map[string]interface{}{"code": 400, "msg": "Bad POST data", "error": map[string]string{field.key: "Required"}}
		}
	}

	if monitorChoice != 0 && monitorChoice != -1 {
		return map[string]interface{}{"code": 400, "msg": "Bad POST data", "error": map[string]string{"monitor_choice": "Invalid"}}
	}

	return nil
}

func (h *Handler) writeResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

