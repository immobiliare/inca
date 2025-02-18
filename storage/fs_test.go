package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFS_ID(t *testing.T) {
	t.Parallel()

	s := FS{}
	if got := s.ID(); got != "FS" {
		t.Errorf("FS.ID() = %v, want %v", got, "FS")
	}
}

func TestFS_Tune(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name: "valid path",
			options: map[string]interface{}{
				"path": "/tmp/test",
			},
			want:    "/tmp/test",
			wantErr: false,
		},
		{
			name:    "missing path",
			options: map[string]interface{}{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &FS{}
			err := s.Tune(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Tune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if s.path != tt.want {
				t.Errorf("FS.Tune() path = %v, want %v", s.path, tt.want)
			}
		})
	}
}

func TestFS_Get(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	s := &FS{path: tmpDir}

	testName := "test"
	testDir := filepath.Join(tmpDir, testName)
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	wantCrt := []byte("test-crt")
	wantKey := []byte("test-key")

	if err := os.WriteFile(filepath.Join(testDir, fsCrtName), wantCrt, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testDir, fsKeyName), wantKey, 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		certName string
		wantCrt  []byte
		wantKey  []byte
		wantErr  bool
	}{
		{
			name:     "valid certificate",
			certName: testName,
			wantCrt:  wantCrt,
			wantKey:  wantKey,
			wantErr:  false,
		},
		{
			name:     "non-existent certificate",
			certName: "invalid",
			wantCrt:  nil,
			wantKey:  nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotCrt, gotKey, err := s.Get(tt.certName)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if string(gotCrt) != string(tt.wantCrt) {
					t.Errorf("FS.Get() gotCrt = %v, want %v", string(gotCrt), string(tt.wantCrt))
				}
				if string(gotKey) != string(tt.wantKey) {
					t.Errorf("FS.Get() gotKey = %v, want %v", string(gotKey), string(tt.wantKey))
				}
			}
		})
	}
}

func TestFS_Put(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	s := &FS{path: tmpDir}

	tests := []struct {
		name     string
		certName string
		crtData  []byte
		keyData  []byte
		wantErr  bool
	}{
		{
			name:     "valid certificate",
			certName: "test",
			crtData:  []byte("test-crt"),
			keyData:  []byte("test-key"),
			wantErr:  false,
		},
		{
			name:     "empty certificate",
			certName: "empty",
			crtData:  []byte{},
			keyData:  []byte{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := s.Put(tt.certName, tt.crtData, tt.keyData)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotCrt, gotKey, err := s.Get(tt.certName)
				if err != nil {
					t.Errorf("FS.Get() error = %v", err)
					return
				}
				if string(gotCrt) != string(tt.crtData) {
					t.Errorf("FS.Put() stored crt = %v, want %v", string(gotCrt), string(tt.crtData))
				}
				if string(gotKey) != string(tt.keyData) {
					t.Errorf("FS.Put() stored key = %v, want %v", string(gotKey), string(tt.keyData))
				}
			}
		})
	}
}

func TestFS_Del(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	s := &FS{path: tmpDir}

	testName := "test"
	if err := s.Put(testName, []byte("test-crt"), []byte("test-key")); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		certName string
		wantErr  bool
	}{
		{
			name:     "existing certificate",
			certName: testName,
			wantErr:  false,
		},
		{
			name:     "non-existent certificate",
			certName: "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := s.Del(tt.certName)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Del() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, _, err := s.Get(tt.certName)
				if err == nil {
					t.Error("FS.Del() certificate still exists after deletion")
				}
			}
		})
	}
}

func TestFS_Find(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	s := &FS{path: tmpDir}

	testCerts := map[string][]byte{
		"test1.example.com": []byte("test1-crt"),
		"test2.example.com": []byte("test2-crt"),
		"other.example.com": []byte("other-crt"),
	}

	for name, crt := range testCerts {
		if err := s.Put(name, crt, []byte("key")); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name    string
		filters []string
		want    int
		wantErr bool
	}{
		{
			name:    "no filters",
			filters: []string{},
			want:    3,
			wantErr: false,
		},
		{
			name:    "filter test*",
			filters: []string{"test.*"},
			want:    2,
			wantErr: false,
		},
		{
			name:    "filter non-existent",
			filters: []string{"invalid.*"},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := s.Find(tt.filters...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("FS.Find() = %v certificates, want %v", len(got), tt.want)
			}
		})
	}
}

func TestFS_Config(t *testing.T) {
	t.Parallel()

	want := "/test/path"
	s := &FS{path: want}

	got := s.Config()
	if got["Path"] != want {
		t.Errorf("FS.Config() = %v, want path %v", got["Path"], want)
	}
}
