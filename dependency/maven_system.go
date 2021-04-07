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
	dependencyDir             = "${user.dir}/.sourcerer-deps"
	mavenDependencyPluginName = "maven-dependency-plugin"
	outFileName               = ".sourcerer-pom.xml"
	outputDirKey              = "outputDirectory"
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

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var pom Project
	if err := xml.Unmarshal(byteValue, &pom); err != nil {
		return errors.Wrapf(err, "failed to unmarshal %s", in)
	}

	overrideCopyDependenciesConfig(&pom)
	err = m.save(&pom)
	if err != nil {
		return err
	}
	return nil
}

func (m *SystemMavenDownloader) save(pom *Project) error {
	out := filepath.Join(m.workingDirectory, outFileName)
	bb, err := xml.MarshalIndent(pom, "", "")
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %s", out)
	}
	err = ioutil.WriteFile(out, bb, 0777)
	return nil
}

func overrideCopyDependenciesConfig(xml *Project) {
	overrideConfigForPlugins(xml.GetBuild().GetPlugins().GetPluginSlice())
	overrideConfigForPlugins(xml.GetBuild().GetPluginManagement().GetPlugins().GetPluginSlice())
	overrideConfigForProfiles(xml)
}

func overrideConfigForProfiles(xml *Project) {
	for _, profile := range xml.GetProfiles().GetProfileSlice() {
		overrideConfigForPlugins(profile.GetBuild().GetPluginManagement().GetPlugins().GetPluginSlice())
		overrideConfigForPlugins(profile.GetBuild().GetPlugins().GetPluginSlice())
	}
}

func overrideConfigForPlugins(plugins []*Plugin) {
	for _, plugin := range plugins {
		if plugin.ArtifactID == mavenDependencyPluginName {
			plugin.Configuration.Entries[outputDirKey] = dependencyDir
		}
	}
}
