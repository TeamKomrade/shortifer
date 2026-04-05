package handler

import (
	"context"
	"database/sql"
	"time"

	pgx "github.com/jackc/pgx/v5"

	http "net/http"
)

type DatabaseHandler struct {
	DatabaseConnString string
}

func (h DatabaseHandler) PingDatabase(res http.ResponseWriter, req *http.Request) {

	conn, err := pgx.Connect(context.Background(), h.DatabaseConnString)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())
	db, err := sql.Open("sqlite", h.DatabaseConnString)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
