package job

import (
	"github.com/koshatul/satis-go/src/satis/satisphp/api"
	"github.com/koshatul/satis-go/src/satis/satisphp/db"
)

// NewSaveRepoJob adds or saves a repo tp the repo collection
func NewSaveRepoJob(dbPath string, repo api.Repo) *SaveRepoJob {
	return &SaveRepoJob{
		dbPath:     dbPath,
		repository: repo,
		exitChan:   make(chan error, 1),
	}
}

// SaveRepoJob needs a comment
type SaveRepoJob struct {
	dbPath     string
	repository api.Repo
	exitChan   chan error
}

// ExitChan needs a comment
func (j SaveRepoJob) ExitChan() chan error {
	return j.exitChan
}

// Run needs a comment
func (j SaveRepoJob) Run() error {
	dbMgr := db.SatisDBManager{Path: j.dbPath}

	if err := dbMgr.Load(); err != nil {
		return err
	}
	repos, err := j.doSave(j.repository, dbMgr.DB.Repositories)
	if err != nil {
		return err
	}
	dbMgr.DB.Repositories = repos

	return dbMgr.Write()
}

func (j SaveRepoJob) doSave(repo api.Repo, repos []db.SatisRepository) ([]db.SatisRepository, error) {
	repoEntity := db.SatisRepository{Type: repo.Type, URL: repo.Url}
	found := false
	for i, r := range repos {
		tmp := api.NewRepo(r.Type, r.URL)
		if tmp.Id == repo.Id {
			repos[i] = repoEntity
			found = true
		}
	}
	if !found {
		return append(repos, repoEntity), nil
	}

	return repos, nil
}
