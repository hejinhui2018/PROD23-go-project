package store

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func (s *Store) SaveCompressedSnapshot(v Snapshot) error {
	f, e := os.Create(filepath.Join(s.Dir(), "snapshot.json.gz"))
	if e != nil {
		return e
	}
	defer f.Close()
	z := gzip.NewWriter(f)
	if e = json.NewEncoder(z).Encode(v); e != nil {
		return e
	}
	return z.Close()
}
func LoadCompressedSnapshot(dir string) (Snapshot, error) {
	var v Snapshot
	f, e := os.Open(filepath.Join(dir, "snapshot.json.gz"))
	if e != nil {
		return v, e
	}
	defer f.Close()
	z, e := gzip.NewReader(f)
	if e != nil {
		return v, e
	}
	defer z.Close()
	e = json.NewDecoder(z).Decode(&v)
	if e == io.EOF {
		return v, nil
	}
	return v, e
}
func SnapshotExists(dir string) bool {
	_, e := os.Stat(filepath.Join(dir, "snapshot.json"))
	return e == nil
}
func SnapshotSize(dir string) int64 {
	st, e := os.Stat(filepath.Join(dir, "snapshot.json"))
	if e != nil {
		return 0
	}
	return st.Size()
}
