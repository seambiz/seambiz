package sdb

import (
	"bytes"
	"database/sql"
	"html/template"

	"github.com/rs/zerolog/log"
)

// OpenDatabaseDSN open DSN and returns Connection Pool. Does not open a Connection. Panics, if DSN is invalid.
func OpenDatabaseDSN(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}

	return db
}

// Connect connects to DB or returns the error.
func Connect(user string, password string, host string, port int, schema string, dsn string) (*sql.DB, error) {
	var err error
	var connString bytes.Buffer

	para := map[string]interface{}{}
	para["User"] = user
	para["Pass"] = password
	para["Host"] = host
	para["Port"] = port
	para["Schema"] = schema

	tmpl, err := template.New("dbconn").Option("missingkey=zero").Parse(dsn)
	if err != nil {
		log.Error().Err(err).Msg("tmpl parse")
		return nil, err
	}

	err = tmpl.Execute(&connString, para)
	if err != nil {
		log.Error().Err(err).Msg("tmpl execute")
		return nil, err
	}

	log.Debug().Str("dsn", connString.String()).Msg("connect to db")
	db, err := sql.Open("mysql", connString.String())
	if err != nil {
		log.Error().Err(err).Msg("mysql connect")
		return nil, err
	}

	return db, nil
}

// OpenDatabase open DSN and returns Connection Pool. Does not open a Connection.
func OpenDatabase(user string, pass string, host string, schema string) *sql.DB {
	var dsn string

	dsn = user
	if pass != "" {
		dsn += ":" + pass
	}
	dsn += "@"
	if host != "" {
		dsn += "tcp("
		dsn += host
		dsn += ")"
	}
	dsn += "/"
	dsn += schema
	dsn += "?parseTime=true"

	return OpenDatabaseDSN(dsn)
}
