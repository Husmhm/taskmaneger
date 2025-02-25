package postgresqluser

import "taskmaneger/repository/postgresql"

type DB struct {
	conn *postgresql.DB
}

func New(conn *postgresql.DB) *DB {
	return &DB{conn: conn}
}
