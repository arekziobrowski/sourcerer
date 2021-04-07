package dependency

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	mavenDependencyPluginName = "maven-dependency-plugin"
	outFileName               = ".sourcerer_pom.xml"
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

	overrideCopyDependenciesConfig(pom)
	err = m.save(pom)
	if err != nil {
		return err
	}
	return nil
}

func (m *SystemMavenDownloader) save(pom *Pom) error {
	out := filepath.Join(m.workingDirectory, outFileName)
	bb, err := xml.MarshalIndent(pom, "", "")
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %s", out)
	}
	err = ioutil.WriteFile(out, bb, 0777)
	return nil
}

func overrideCopyDependenciesConfig(xml *Pom) {
	overrideConfigForPlugins(xml.GetProject().GetBuild().GetPlugins().GetPluginSlice())
	overrideConfigForPlugins(xml.GetProject().GetBuild().GetPluginManagement().GetPlugins().GetPluginSlice())
	overrideConfigForProfiles(xml)
}

func overrideConfigForProfiles(xml *Pom) {
	for _, profile := range xml.GetProject().GetProfiles().GetProfileSlice() {
		overrideConfigForPlugins(profile.GetBuild().GetPluginManagement().GetPlugins().GetPluginSlice())
		overrideConfigForPlugins(profile.GetBuild().GetPlugins().GetPluginSlice())
	}
}

func overrideConfigForPlugins(plugins []*Plugin) {
	for _, plugin := range plugins {
		if plugin.ArtifactID == mavenDependencyPluginName {
			log.Infof("Config: %s", plugin.Configuration)
		}
	}
}
