package satisphp

import (
	"errors"
	"log"

	"satisphp/api"
	"satisphp/job"
)

var _ = log.Print

// ErrRepoNotFound is a error returned when a repository is not found
var ErrRepoNotFound = errors.New("Repository Not Found")

// SatisClient needs a comment
type SatisClient struct {
	Jobs   chan job.SatisJob
	DBPath string
}

// FindRepo needs a comment
func (s *SatisClient) FindRepo(id string) (api.Repo, error) {
	var repo api.Repo

	repos, err := s.FindAllRepos()
	if err != nil {
		return repo, err
	}

	found := false
	for _, r := range repos {
		if r.Id == id {
			found = true
			repo = r
		}
	}
	if found {
		return repo, nil
	}
	return repo, ErrRepoNotFound
}

// FindAllRepos needs a comment
func (s *SatisClient) FindAllRepos() ([]api.Repo, error) {
	j := job.NewFindAllJob(s.DBPath)

	err := s.performJob(j)

	repos := <-j.ReposResp

	rs := make([]api.Repo, len(repos))
	for i, repo := range repos {
		rs[i] = *api.NewRepo(repo.Type, repo.URL)
	}

	return rs, err
}

// SaveRepo nees a comment
func (s *SatisClient) SaveRepo(repo *api.Repo, generate bool) error {
	// repoEntity := db.SatisRepository{
	// 	Type: repo.Type,
	// 	Url:  repo.Url,
	// }
	j := job.NewSaveRepoJob(s.DBPath, *repo)
	if err := s.performJob(j); err != nil {
		return err
	}
	if generate {
		return s.GenerateSatisWeb()
	} else {
		return nil
	}
}

// DeleteRepo needs a comment
func (s *SatisClient) DeleteRepo(id string, generate bool) error {
	var toDelete api.Repo

	repos, err := s.FindAllRepos()
	if err != nil {
		return err
	}

	found := false
	for _, r := range repos {
		if r.Id == id {
			found = true
			toDelete = r
		}
	}

	if found {
		j := job.NewDeleteRepoJob(s.DBPath, toDelete.Url)
		if err = s.performJob(j); err != nil {
			switch err {
			case job.ErrRepoNotFound:
				return ErrRepoNotFound
			default:
				return err
			}
		}

		if generate {
			return s.GenerateSatisWeb()
		} else {
			return nil
		}
	} else {
		return ErrRepoNotFound
	}
}

// GenerateSatisWeb needs a comment
func (s *SatisClient) GenerateSatisWeb() error {
	j := job.NewGenerateJob()
	return s.performJob(j)
}

// Shutdown needs a comment
func (s *SatisClient) Shutdown() error {
	j := job.NewExitJob()
	return s.performJob(j)
}

func (s *SatisClient) performJob(j job.SatisJob) error {
	s.Jobs <- j

	return <-j.ExitChan()
}
