package api

import (
	"fmt"
	"hash/crc32"
)

// Repo needs a comment
type Repo struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

// NewRepo needs a comment
func NewRepo(t string, u string) *Repo {
	crc := crc32.NewIEEE()
	crc.Write([]byte(u))
	v := crc.Sum32()

	return &Repo{
		Id:   fmt.Sprint(v),
		Type: t,
		Url:  u,
	}
}
