package db

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var _ = log.Print

const (
	// DBFile needs a comment
	DBFile = "/db.json"
	// StagingFile needs a comment
	StagingFile = "/stage.json"
)

// SatisDBManager needs a comment
type SatisDBManager struct {
	Path string
	DB   SatisDB
}

// Load needs a comment
func (c *SatisDBManager) Load() error {
	content, err := ioutil.ReadFile(c.Path + DBFile)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(content, &c.DB); err != nil {
		return err
	}
	return nil
}

// Write needs a comment
func (c *SatisDBManager) Write() error {
	return c.doWrite(c.Path + DBFile)
}

// WriteStaging needs a comment
func (c *SatisDBManager) WriteStaging() error {
	return c.doWrite(c.Path + StagingFile)
}

func (c *SatisDBManager) doWrite(path string) error {
	b, err := json.MarshalIndent(c.DB, "", "    ") // pretty print
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path, b, 0644); err != nil {
		return err
	}
	return nil
}

// SaveRepo needs a comment
func (c *SatisDBManager) SaveRepo(repo SatisRepository) error {
	return nil
}
