package client

import (
	"log"
	"net/http"

	"satisphp/api"
)

var _ = log.Print

// SatisClient needs a comment
type SatisClient struct {
	Host string
}

// AddRepo adds a repo to the SatisClient
func (c *SatisClient) AddRepo(repo *api.Repo) (*api.Repo, error) {
	r := &api.Repo{}
	url := c.Host + "/api/repo"

	req, err := makeRequest("POST", url, repo)
	if err != nil {
		return r, err
	}
	err = processResponseEntity(req, &r, http.StatusCreated)
	return r, err
}

// SaveRepo saves updates to an existing repo on SatisClient
func (c *SatisClient) SaveRepo(repo *api.Repo) (*api.Repo, error) {
	r := &api.Repo{}
	url := c.Host + "/api/repo/" + repo.Id

	req, err := makeRequest("PUT", url, repo)
	if err != nil {
		return r, err
	}
	err = processResponseEntity(req, &r, http.StatusOK)
	return r, err
}

// FindRepo searches for a repo
func (c *SatisClient) FindRepo(id string) (*api.Repo, error) {
	var repo api.Repo
	url := c.Host + "/api/repo/" + id

	req, err := makeRequest("GET", url, nil)
	if err != nil {
		return &repo, err
	}
	err = processResponseEntity(req, &repo, http.StatusOK)
	return &repo, err
}

// FindAllRepos returns all repositories
func (c *SatisClient) FindAllRepos() ([]api.Repo, error) {
	var repos []api.Repo
	url := c.Host + "/api/repo"

	req, err := makeRequest("GET", url, nil)
	if err != nil {
		return repos, err
	}
	err = processResponseEntity(req, &repos, http.StatusOK)
	return repos, err
}

// DeleteRepo deletes a repo
func (c *SatisClient) DeleteRepo(id string) error {
	url := c.Host + "/api/repo/" + id

	req, err := makeRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	return processResponseEntity(req, nil, http.StatusNoContent)
}

// GenerateStaticWeb needs a comment
func (c *SatisClient) GenerateStaticWeb() error {
	url := c.Host + "/api/generate-web-job"

	req, err := makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	return processResponseEntity(req, nil, http.StatusCreated)
}
