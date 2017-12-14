package job

import (
	"github.com/koshatul/satis-go/src/satis/satisphp/db"
)

// NewDeleteRepoJob removes a repo from the repo collection
func NewDeleteRepoJob(dbPath string, repo string) *DeleteRepoJob {
	return &DeleteRepoJob{
		dbPath:     dbPath,
		repository: repo,
		exitChan:   make(chan error, 1),
	}
}

// DeleteRepoJob needs a comment
type DeleteRepoJob struct {
	dbPath     string
	repository string
	exitChan   chan error
}

// ExitChan needs a comment
func (j DeleteRepoJob) ExitChan() chan error {
	return j.exitChan
}

// Run needs a comment
func (j DeleteRepoJob) Run() error {
	dbMgr := db.SatisDBManager{Path: j.dbPath}

	if err := dbMgr.Load(); err != nil {
		return err
	}
	repos, err := j.doDelete(j.repository, dbMgr.DB.Repositories)
	if err != nil {
		return err
	}
	dbMgr.DB.Repositories = repos

	return dbMgr.Write()
}
func (j DeleteRepoJob) doDelete(repo string, repos []db.SatisRepository) ([]db.SatisRepository, error) {
	var err error = nil
	found := false

	rs := make([]db.SatisRepository, 0, len(repos))
	for _, r := range repos {
		if r.URL == repo {
			found = true
		} else {
			rs = append(rs, r)
		}
	}
	if !found {
		err = ErrRepoNotFound
	}
	return rs, err
}
