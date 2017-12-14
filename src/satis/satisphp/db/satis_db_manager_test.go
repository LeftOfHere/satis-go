package db

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

var _ = fmt.Print
var _ = log.Print

func ARandomDbMgr() SatisDBManager {
	dbPath := "/tmp/satis-test-data"

	// Make Data Dir
	if err := os.MkdirAll(dbPath, 0744); err != nil {
		log.Fatalf("Unable to create path: %v", err)
	}

	mgr := SatisDBManager{Path: dbPath}
	mgr.DB.Name = "My Repo"
	mgr.DB.Homepage = "http://repo.com"
	mgr.DB.RequireAll = true
	mgr.DB.Repositories = []SatisRepository{
		{Type: "vcs", URL: "http://package.com"},
	}

	mgr.Path = dbPath
	mgr.Write()

	return mgr
}

func TestDbLoad(t *testing.T) {

	// given
	mgr := ARandomDbMgr()
	r := SatisDBManager{Path: mgr.Path}

	// when
	err := r.Load()

	// then
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(r.DB, mgr.DB) {
		t.Error("loaded config doesn't match original")
	}
}
func TestDbWrite(t *testing.T) {
	// given
	r := ARandomDbMgr()
	oldName := r.DB.Name
	// when
	r.DB.Name = "foo"
	modifiedDb := r.DB

	err := r.Write()

	// then
	if err != nil {
		t.Error(err)
	}

	err = r.Load()
	if err != nil {
		t.Error(err)
	}

	if oldName == r.DB.Name {
		t.Errorf("config should have changed: %s / %s", oldName, r.DB.Name)
	}
	if !reflect.DeepEqual(r.DB, modifiedDb) {
		t.Errorf("config didn't persist changes when written: %s / %s", r.DB.Name, modifiedDb.Name)
	}
}
