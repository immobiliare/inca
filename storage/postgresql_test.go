package storage

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgreSQL_ParseConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			options: map[string]interface{}{
				"host":     "localhost",
				"port":     5432,
				"user":     "postgres",
				"password": "secret",
				"dbname":   "testdb",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			options: map[string]interface{}{
				"port":     5432,
				"user":     "postgres",
				"password": "secret",
				"dbname":   "testdb",
			},
			wantErr: true,
			errMsg:  "provider PostgreSQL: host not defined",
		},
		{
			name: "missing port",
			options: map[string]interface{}{
				"host":     "localhost",
				"user":     "postgres",
				"password": "secret",
				"dbname":   "testdb",
			},
			wantErr: true,
			errMsg:  "provider PostgreSQL: port not defined",
		},
		{
			name: "missing user",
			options: map[string]interface{}{
				"host":     "localhost",
				"port":     5432,
				"password": "secret",
				"dbname":   "testdb",
			},
			wantErr: true,
			errMsg:  "provider PostgreSQL: user not defined",
		},
		{
			name: "missing password",
			options: map[string]interface{}{
				"host":   "localhost",
				"port":   5432,
				"user":   "postgres",
				"dbname": "testdb",
			},
			wantErr: true,
			errMsg:  "provider PostgreSQL: password not defined",
		},
		{
			name: "missing dbname",
			options: map[string]interface{}{
				"host":     "localhost",
				"port":     5432,
				"user":     "postgres",
				"password": "secret",
			},
			wantErr: true,
			errMsg:  "provider PostgreSQL: dbname not defined",
		},
		{
			name: "valid config",
			options: map[string]interface{}{
				"host":     "localhost",
				"port":     5432,
				"user":     "postgres",
				"password": "secret",
				"dbname":   "testdb",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &PostgreSQL{}
			db, err := s.parseConfig(tt.options)

			if tt.wantErr {
				if err == nil {
					t.Error("parseConfig() error = nil, want error")
				}
				if db != nil {
					t.Error("parseConfig() db != nil, want nil")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("parseConfig() error = %v, want %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("parseConfig() error = %v, want nil", err)
				}
				if db == nil {
					t.Error("parseConfig() db = nil, want not nil")
				}
			}
		})
	}
}

func TestPostgreSQL_Connect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		dbSetup func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "successful connection",
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS certificates").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
		{
			name: "ping fails",
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing().WillReturnError(fmt.Errorf("ping failed"))
			},
			wantErr: true,
		},
		{
			name: "create table fails",
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS certificates").
					WillReturnError(fmt.Errorf("create table failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer func() {
				if err := db.Close(); err != nil {
					t.Logf("Failed to close database: %v", err)
				}
			}()

			tt.dbSetup(mock)

			s := &PostgreSQL{db: db}
			err = s.connect()

			if tt.wantErr && err == nil {
				t.Error("connect() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("connect() error = %v, want nil", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgreSQL_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		certName string
		dbSetup  func(mock sqlmock.Sqlmock)
		wantCrt  []byte
		wantKey  []byte
		wantErr  bool
	}{
		{
			name:     "successful get",
			certName: "test-cert",
			dbSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"crt_data", "key_data"}).
					AddRow([]byte("test-crt"), []byte("test-key"))
				mock.ExpectQuery("SELECT crt_data, key_data FROM certificates WHERE name = \\$1").
					WithArgs("test-cert").
					WillReturnRows(rows)
			},
			wantCrt: []byte("test-crt"),
			wantKey: []byte("test-key"),
			wantErr: false,
		},
		{
			name:     "certificate not found",
			certName: "non-existent",
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT crt_data, key_data FROM certificates WHERE name = \\$1").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			wantCrt: nil,
			wantKey: nil,
			wantErr: true,
		},
		{
			name:     "query error",
			certName: "test-cert",
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT crt_data, key_data FROM certificates WHERE name = \\$1").
					WithArgs("test-cert").
					WillReturnError(fmt.Errorf("database error"))
			},
			wantCrt: nil,
			wantKey: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer func() {
				if err := db.Close(); err != nil {
					t.Logf("Failed to close database: %v", err)
				}
			}()

			tt.dbSetup(mock)

			s := &PostgreSQL{db: db}
			gotCrt, gotKey, err := s.Get(tt.certName)

			if tt.wantErr && err == nil {
				t.Error("Get() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Get() error = %v, want nil", err)
			}
			if !bytes.Equal(gotCrt, tt.wantCrt) {
				t.Errorf("Get() gotCrt = %v, want %v", gotCrt, tt.wantCrt)
			}
			if !bytes.Equal(gotKey, tt.wantKey) {
				t.Errorf("Get() gotKey = %v, want %v", gotKey, tt.wantKey)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgreSQL_Put(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		certName string
		crtData  []byte
		keyData  []byte
		dbSetup  func(mock sqlmock.Sqlmock)
		wantErr  bool
	}{
		{
			name:     "successful insert",
			certName: "test-cert",
			crtData:  []byte("test-crt"),
			keyData:  []byte("test-key"),
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO certificates \\(name, crt_data, key_data\\) VALUES \\(\\$1, \\$2, \\$3\\)").
					WithArgs("test-cert", []byte("test-crt"), []byte("test-key")).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:     "successful update",
			certName: "existing-cert",
			crtData:  []byte("updated-crt"),
			keyData:  []byte("updated-key"),
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO certificates \\(name, crt_data, key_data\\) VALUES \\(\\$1, \\$2, \\$3\\)").
					WithArgs("existing-cert", []byte("updated-crt"), []byte("updated-key")).
					WillReturnResult(sqlmock.NewResult(1, 2))
			},
			wantErr: false,
		},
		{
			name:     "database error",
			certName: "test-cert",
			crtData:  []byte("test-crt"),
			keyData:  []byte("test-key"),
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO certificates \\(name, crt_data, key_data\\) VALUES \\(\\$1, \\$2, \\$3\\)").
					WithArgs("test-cert", []byte("test-crt"), []byte("test-key")).
					WillReturnError(fmt.Errorf("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer func() {
				if err := db.Close(); err != nil {
					t.Logf("Failed to close database: %v", err)
				}
			}()

			tt.dbSetup(mock)

			s := &PostgreSQL{db: db}
			err = s.Put(tt.certName, tt.crtData, tt.keyData)

			if tt.wantErr && err == nil {
				t.Error("Put() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Put() error = %v, want nil", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgreSQL_Del(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		certName string
		dbSetup  func(mock sqlmock.Sqlmock)
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "successful delete",
			certName: "test-cert",
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM certificates WHERE name = \\$1").
					WithArgs("test-cert").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:     "certificate not found",
			certName: "non-existent",
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM certificates WHERE name = \\$1").
					WithArgs("non-existent").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
			errMsg:  "certificate not found: non-existent",
		},
		{
			name:     "database error",
			certName: "test-cert",
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM certificates WHERE name = \\$1").
					WithArgs("test-cert").
					WillReturnError(fmt.Errorf("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer func() {
				if err := db.Close(); err != nil {
					t.Logf("Failed to close database: %v", err)
				}
			}()

			tt.dbSetup(mock)

			s := &PostgreSQL{db: db}
			err = s.Del(tt.certName)

			if tt.wantErr {
				if err == nil {
					t.Error("Del() error = nil, want error")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Del() error = %v, want %v", err, tt.errMsg)
				}
			} else if err != nil {
				t.Errorf("Del() error = %v, want nil", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgreSQL_Find(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		filters []string
		dbSetup func(mock sqlmock.Sqlmock)
		want    [][]byte
		wantErr bool
	}{
		{
			name:    "successful find with no filters",
			filters: []string{},
			dbSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "crt_data"}).
					AddRow("test1.example.com", []byte("cert1")).
					AddRow("test2.example.com", []byte("cert2"))
				mock.ExpectQuery("SELECT name, crt_data FROM certificates").
					WillReturnRows(rows)
			},
			want:    [][]byte{[]byte("cert1"), []byte("cert2")},
			wantErr: false,
		},
		{
			name:    "successful find with filter",
			filters: []string{"test1"},
			dbSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "crt_data"}).
					AddRow("test1.example.com", []byte("cert1")).
					AddRow("test2.example.com", []byte("cert2"))
				mock.ExpectQuery("SELECT name, crt_data FROM certificates").
					WillReturnRows(rows)
			},
			want:    [][]byte{[]byte("cert1")},
			wantErr: false,
		},
		{
			name:    "query error",
			filters: []string{},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name, crt_data FROM certificates").
					WillReturnError(fmt.Errorf("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "scan error",
			filters: []string{},
			dbSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name", "crt_data"}).
					AddRow("test1.example.com", nil)
				mock.ExpectQuery("SELECT name, crt_data FROM certificates").
					WillReturnRows(rows)
			},
			want:    [][]byte{nil},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer func() {
				if err := db.Close(); err != nil {
					t.Logf("Failed to close database: %v", err)
				}
			}()

			tt.dbSetup(mock)

			s := &PostgreSQL{db: db}
			got, err := s.Find(tt.filters...)

			if tt.wantErr && err == nil {
				t.Error("Find() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Find() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgreSQL_Config(t *testing.T) {
	t.Parallel()

	options := map[string]interface{}{
		"host":     "localhost",
		"port":     5432,
		"user":     "postgres",
		"password": "secret",
		"dbname":   "testdb",
	}
	wanted := "postgres://postgres:secret@localhost:5432/testdb"

	s := &PostgreSQL{}
	if _, err := s.parseConfig(options); err != nil {
		t.Fatalf("parseConfig() error = %v, want nil", err)
	}

	if config := s.Config(); config["Dsn"] != wanted {
		t.Errorf("Config() = %v, want %v", config, wanted)
	}
}
