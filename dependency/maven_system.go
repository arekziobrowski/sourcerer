package dependency

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type SystemMavenDownloader struct {
	workingDirectory string
}

func NewSystemMavenDownloader(wd string) *SystemMavenDownloader {
	return &SystemMavenDownloader{
		workingDirectory: wd,
	}
}

func (m *SystemMavenDownloader) Get(src string) error {
	in := filepath.Join(m.workingDirectory, src)
	f, err := os.Open(in)
	if os.IsNotExist(err) {
		log.Warnf("There is no %s", in)
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", in)
	}
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)

	var pom *Pom
	if err := xml.Unmarshal(byteValue, pom); err != nil {
		return errors.Wrapf(err, "failed to unmarshal %s", in)
	}

	return nil
}
