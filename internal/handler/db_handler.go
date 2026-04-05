package handler

import (
	"context"

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

	res.WriteHeader(http.StatusOK)
}
