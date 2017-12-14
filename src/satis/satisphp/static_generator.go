package satisphp

import (
	"log"
	"os/exec"

	"github.com/koshatul/satis-go/src/satis/satisphp/db"
)

var _ = log.Print

// Generator needs a comment
type Generator interface {
	Generate() error
}

// StaticWebGenerator needs a comment
type StaticWebGenerator struct {
	DBPath  string
	WebPath string
}

// Generate needs a comment
func (s *StaticWebGenerator) Generate() error {
	log.Print("Generating...")
	out, err := exec.
		Command("satis", "--no-interaction", "build", s.DBPath+db.StagingFile, s.WebPath).
		CombinedOutput()
	if err != nil {
		log.Printf("Satis Generation Error: %s", string(out[:]))
	}
	return err
}
