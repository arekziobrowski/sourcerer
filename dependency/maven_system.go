package dependency

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	dependencyDir             = ".sourcerer-deps"
	mavenDependencyPluginName = "maven-dependency-plugin"
	pomName                   = "pom.xml"
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

func (m *SystemMavenDownloader) Get() error {
	in := filepath.Join(m.workingDirectory, pomName)
	f, err := os.Open(in)
	if os.IsNotExist(err) {
		log.Warnf("There is no %s", in)
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", in)
	}
	defer f.Close()

	log.Infof("Downloading dependencies to %s", filepath.Join(m.workingDirectory, dependencyDir))

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

	err = m.download()
	if err != nil {
		return err
	}

	log.Infof("Finished downloading dependencies to %s", filepath.Join(m.workingDirectory, dependencyDir))
	return nil
}

func (m *SystemMavenDownloader) save(pom *Project) error {
	out := filepath.Join(m.workingDirectory, outFileName)
	bb, err := xml.MarshalIndent(pom, "", "\t")
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %s", out)
	}
	err = ioutil.WriteFile(out, bb, 0777)
	return nil
}

func (m *SystemMavenDownloader) download() error {
	destDir := filepath.Join(m.workingDirectory, dependencyDir)
	// prepare directory
	err := os.MkdirAll(destDir, 0777)
	if err != nil {
		return errors.Wrap(err, "cannot create destination directory")
	}
	cmd := "mvn"
	if runtime.GOOS == "windows" {
		cmd = "mvn.cmd"
	}
	err = run(m.workingDirectory, cmd, "dependency:copy-dependencies", fmt.Sprintf("-DoutputDirectory=%s", destDir), "-f", filepath.Join(m.workingDirectory, outFileName))
	if err != nil {
		return errors.Wrapf(err, "failed to download deps to %s", destDir)
	}
	return nil
}

func run(wd, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = wd
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Errorf("Error occured when running %q: %s", command+" "+strings.Join(args, " "), stderr.String())
		return err
	}
	return nil
}

func overrideCopyDependenciesConfig(xml *Project) {
	overrideConfigForPlugins(xml.Name, xml.GetBuild().GetPlugins().GetPluginSlice())
	overrideConfigForPlugins(xml.Name, xml.GetBuild().GetPluginManagement().GetPlugins().GetPluginSlice())
	overrideConfigForProfiles(xml)
}

func overrideConfigForProfiles(xml *Project) {
	for _, profile := range xml.GetProfiles().GetProfileSlice() {
		overrideConfigForPlugins(xml.Name, profile.GetBuild().GetPluginManagement().GetPlugins().GetPluginSlice())
		overrideConfigForPlugins(xml.Name, profile.GetBuild().GetPlugins().GetPluginSlice())
	}
}

func overrideConfigForPlugins(projectName string, plugins []*Plugin) {
	for _, plugin := range plugins {
		if plugin.ArtifactID == mavenDependencyPluginName {
			log.Infof("Found maven-dependency-plugin for %s", projectName)
			if plugin.Configuration != nil && plugin.Configuration.Entries != nil {
				plugin.Configuration = nil
			}
			clearExecutions(plugin.Executions)
		}
	}
}

func clearExecutions(execution *Executions) {
	if execution == nil {
		return
	}
	for _, e := range execution.Execution {
		e.Configuration = nil
	}
}
