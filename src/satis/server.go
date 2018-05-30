package satis

import (
	"log"
	"net/http"
	"os"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"github.com/leftofhere/satis-go/src/satis/satisphp"
	"github.com/leftofhere/satis-go/src/satis/satisphp/db"
	"github.com/leftofhere/satis-go/src/satis/satisphp/job"
)

var _ = log.Printf

// Server needs a comment
type Server struct {
	DBPath       string
	AdminUIPath  string
	WebPath      string
	Bind         string
	Name         string
	Homepage     string
	jobProcessor satisphp.SatisJobProcessor
	jobClient    satisphp.SatisClient
}

// Run needs a comment
func (s *Server) Run() error {
	// sync config to db
	if err := s.initDb(); err != nil {
		return err
	}

	// Shared Jobs Channel to queue/process db modifications and generation task
	jobs := make(chan job.SatisJob)

	// Job Processor responsible for interacting with db & static web docs
	gen := &satisphp.StaticWebGenerator{
		DBPath:  s.DBPath,
		WebPath: s.WebPath,
	}

	s.jobProcessor = satisphp.SatisJobProcessor{
		DBPath:    s.DBPath,
		Jobs:      jobs,
		Generator: gen,
	}

	// Client to Job Processor
	jobClient := satisphp.SatisClient{
		DBPath: s.DBPath,
		Jobs:   jobs,
	}

	// route handlers
	resource := &SatisResource{
		Host:           s.Homepage,
		SatisPhpClient: jobClient,
	}

	// Configure Routes
	r := mux.NewRouter()

	r.HandleFunc("/api/repo", resource.addRepo).Methods("POST")
	r.HandleFunc("/api/repo/{id}", resource.saveRepo).Methods("PUT")
	r.HandleFunc("/api/repo/{id}", resource.findRepo).Methods("GET")
	r.HandleFunc("/api/repo", resource.findAllRepos).Methods("GET")
	r.HandleFunc("/api/repo/{id}", resource.deleteRepo).Methods("DELETE")
	r.HandleFunc("/api/generate-web-job", resource.generateStaticWeb).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(s.WebPath)))

	//	r.Handle("/dist/{rest}", http.StripPrefix("/dist/", http.FileServer(http.Dir("./dist/"))))
	// r.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("./dist"))))

	username := os.Getenv("SATIS_GO_USERNAME")
	password := os.Getenv("SATIS_GO_PASSWORD")
	if username != "" && password != "" {
		http.Handle("/", httpauth.SimpleBasicAuth(username, password)(r))
		http.Handle("/admin/", httpauth.SimpleBasicAuth(username, password)(http.StripPrefix("/admin/", http.FileServer(http.Dir(s.AdminUIPath)))))
	} else {
		http.Handle("/", r)
		http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir(s.AdminUIPath))))
	}

	// Start update processor
	go s.jobProcessor.ProcessUpdates()

	// Start HTTP Server
	return http.ListenAndServe(s.Bind, nil)
}

// Sync configured values to satis repository meta data
func (s *Server) initDb() error {
	dbMgr := &db.SatisDBManager{Path: s.DBPath}

	// create empty db if it doesn't exist
	if _, err := os.Stat(s.DBPath + db.DBFile); os.IsNotExist(err) {
		if err := dbMgr.Write(); err != nil {
			return err
		}
	}

	if err := dbMgr.Load(); err != nil {
		return err
	}
	dbMgr.DB.Name = s.Name
	dbMgr.DB.Homepage = s.Homepage
	dbMgr.DB.RequireAll = true
	return dbMgr.Write()
}
