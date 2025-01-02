package storage

import (
	"database/sql"
	"fmt"

	"github.com/immobiliare/inca/pki"
	"github.com/immobiliare/inca/util"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQL struct {
	Storage
	dsn string
	db  *sql.DB
}

func (s PostgreSQL) ID() string {
	return "PostgreSQL"
}

func (s *PostgreSQL) Tune(options map[string]interface{}) (err error) {
	s.db, err = s.parseConfig(options)
	if err != nil {
		return err
	}

	if err = s.connect(); err != nil {
		return err
	}

	return nil
}

func (s *PostgreSQL) parseConfig(options map[string]interface{}) (*sql.DB, error) {
	host, ok := options["host"]
	if !ok {
		return nil, fmt.Errorf("provider %s: host not defined", s.ID())
	}

	port, ok := options["port"]
	if !ok {
		return nil, fmt.Errorf("provider %s: port not defined", s.ID())
	}

	user, ok := options["user"]
	if !ok {
		return nil, fmt.Errorf("provider %s: user not defined", s.ID())
	}

	password, ok := options["password"]
	if !ok {
		return nil, fmt.Errorf("provider %s: password not defined", s.ID())
	}

	dbname, ok := options["dbname"]
	if !ok {
		return nil, fmt.Errorf("provider %s: dbname not defined", s.ID())
	}

	s.dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		user.(string), password.(string), host.(string), port.(int), dbname.(string))

	return sql.Open("pgx", s.dsn)
}

func (s *PostgreSQL) connect() (err error) {
	err = s.db.Ping()
	if err != nil {
		return
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS certificates (
			name VARCHAR(255) PRIMARY KEY,
			crt_data BYTEA NOT NULL,
			key_data BYTEA NOT NULL
		)
	`)
	return
}

func (s *PostgreSQL) Get(name string) ([]byte, []byte, error) {
	var crtData, keyData []byte
	err := s.db.QueryRow(
		"SELECT crt_data, key_data FROM certificates WHERE name = $1", name).
		Scan(&crtData, &keyData)
	if err != nil {
		return nil, nil, err
	}
	return crtData, keyData, nil
}

func (s *PostgreSQL) Put(name string, crtData, keyData []byte) error {
	_, err := s.db.Exec(`
		INSERT INTO certificates (name, crt_data, key_data)
		VALUES ($1, $2, $3)
		ON CONFLICT (name) DO UPDATE
		SET crt_data = EXCLUDED.crt_data,
			key_data = EXCLUDED.key_data
	`, name, crtData, keyData)
	return err
}

func (s *PostgreSQL) Del(name string) error {
	result, err := s.db.Exec("DELETE FROM certificates WHERE name = $1", name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("certificate not found: %s", name)
	}
	return nil
}

func (s *PostgreSQL) Find(filters ...string) ([][]byte, error) {
	rows, err := s.db.Query("SELECT name, crt_data FROM certificates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results [][]byte
	for rows.Next() {
		var name string
		var crtData []byte
		if err := rows.Scan(&name, &crtData); err != nil {
			return nil, err
		}

		if pki.IsValidCN(name) && util.RegexesMatch(name, filters...) {
			results = append(results, crtData)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *PostgreSQL) Config() map[string]string {
	return map[string]string{
		"Dsn": s.dsn,
	}
}
