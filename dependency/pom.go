package dependency

import (
	"encoding/xml"
	"io"
)

type Project struct {
	XMLName                xml.Name                `xml:"project"`
	ModelVersion           string                  `xml:"modelVersion,omitempty"`
	Parent                 *Parent                 `xml:"parent,omitempty"`
	GroupID                string                  `xml:"groupId,omitempty"`
	ArtifactID             string                  `xml:"artifactId,omitempty"`
	Version                string                  `xml:"version,omitempty"`
	Packaging              string                  `xml:"packaging,omitempty"`
	Name                   string                  `xml:"name,omitempty"`
	Description            string                  `xml:"description,omitempty"`
	URL                    string                  `xml:"url,omitempty"`
	InceptionYear          string                  `xml:"inceptionYear,omitempty"`
	Organization           *Organization           `xml:"organization,omitempty"`
	Licenses               *Licenses               `xml:"licenses,omitempty"`
	Developers             *Developers             `xml:"developers,omitempty"`
	Contributors           *Contributors           `xml:"contributors,omitempty"`
	MailingLists           *MailingLists           `xml:"mailingLists,omitempty"`
	Prerequisites          *Prerequisites          `xml:"prerequisites,omitempty"`
	Modules                string                  `xml:"modules,omitempty"`
	SCM                    *SCM                    `xml:"scm,omitempty"`
	IssueManagement        *IssueManagement        `xml:"issueManagement,omitempty"`
	CiManagement           *CiManagement           `xml:"ciManagement,omitempty"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement,omitempty"`
	Properties             *Properties             `xml:"properties,omitempty"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement,omitempty"`
	Dependencies           *Dependencies           `xml:"dependencies,omitempty"`
	Repositories           *Repositories           `xml:"repositories,omitempty"`
	PluginRepositories     *PluginRepositories     `xml:"pluginRepositories,omitempty"`
	Build                  *ProjectBuild           `xml:"build,omitempty"`
	Reports                string                  `xml:"reports,omitempty"`
	Reporting              *Reporting              `xml:"reporting,omitempty"`
	Profiles               *Profiles               `xml:"profiles,omitempty"`
	Xmlns                  string                  `xml:"_xmlns,omitempty"`
	XmlnsXsi               string                  `xml:"_xmlns:xsi,omitempty"`
	XsiSchemaLocation      string                  `xml:"_xsi:schemaLocation,omitempty"`
}

func (m *Project) GetBuild() *ProjectBuild {
	if m != nil {
		return m.Build
	}
	return nil
}

func (m *Project) GetProfiles() *Profiles {
	if m != nil {
		return m.Profiles
	}
	return nil
}

type ProjectBuild struct {
	SourceDirectory       string            `xml:"sourceDirectory,omitempty"`
	ScriptSourceDirectory string            `xml:"scriptSourceDirectory,omitempty"`
	TestSourceDirectory   string            `xml:"testSourceDirectory,omitempty"`
	OutputDirectory       string            `xml:"outputDirectory,omitempty"`
	TestOutputDirectory   string            `xml:"testOutputDirectory,omitempty"`
	Extensions            *Extensions       `xml:"extensions,omitempty"`
	DefaultGoal           string            `xml:"defaultGoal,omitempty"`
	Resources             *Resources        `xml:"resources,omitempty"`
	TestResources         *TestResources    `xml:"testResources,omitempty"`
	Directory             string            `xml:"directory,omitempty"`
	FinalName             string            `xml:"finalName,omitempty"`
	Filters               string            `xml:"filters,omitempty"`
	PluginManagement      *PluginManagement `xml:"pluginManagement,omitempty"`
	Plugins               *Plugins          `xml:"plugins,omitempty"`
}

func (m *ProjectBuild) GetPlugins() *Plugins {
	if m != nil {
		return m.Plugins
	}
	return nil
}

func (m *ProjectBuild) GetPluginManagement() *PluginManagement {
	if m != nil {
		return m.PluginManagement
	}
	return nil
}

type Extensions struct {
	Extension []*Parent `xml:"extension,omitempty"`
}

type Parent struct {
	GroupID      string  `xml:"groupId,omitempty"`
	ArtifactID   string  `xml:"artifactId,omitempty"`
	Version      string  `xml:"version,omitempty"`
	Message      *string `xml:"message,omitempty,omitempty"`
	RelativePath *string `xml:"relativePath,omitempty,omitempty"`
}

type PluginManagement struct {
	Plugins *Plugins `xml:"plugins,omitempty"`
}

func (m *PluginManagement) GetPlugins() *Plugins {
	if m != nil {
		return m.Plugins
	}
	return nil
}

type Plugins struct {
	Plugin []*Plugin `xml:"plugin,omitempty"`
}

type Plugin struct {
	GroupID       string        `xml:"groupId,omitempty"`
	ArtifactID    string        `xml:"artifactId,omitempty"`
	Version       string        `xml:"version,omitempty"`
	Extensions    string        `xml:"extensions,omitempty"`
	Executions    *Executions   `xml:"executions,omitempty"`
	Dependencies  *Dependencies `xml:"dependencies,omitempty"`
	Goals         string        `xml:"goals,omitempty"`
	Inherited     string        `xml:"inherited,omitempty"`
	Configuration *Properties   `xml:"configuration,omitempty"`
}

func (m *Plugins) GetPluginSlice() []*Plugin {
	if m != nil {
		return m.Plugin
	}
	return nil
}

type Dependencies struct {
	Dependency []*Dependency `xml:"dependency,omitempty"`
}

type Dependency struct {
	GroupID    string      `xml:"groupId,omitempty"`
	ArtifactID string      `xml:"artifactId,omitempty"`
	Version    string      `xml:"version,omitempty"`
	Type       string      `xml:"type,omitempty"`
	Classifier string      `xml:"classifier,omitempty"`
	Scope      string      `xml:"scope,omitempty"`
	SystemPath string      `xml:"systemPath,omitempty"`
	Exclusions *Exclusions `xml:"exclusions,omitempty"`
	Optional   string      `xml:"optional,omitempty"`
}

type Exclusions struct {
	Exclusion []*Exclusion `xml:"exclusion,omitempty"`
}

type Exclusion struct {
	ArtifactID string `xml:"artifactId,omitempty"`
	GroupID    string `xml:"groupId,omitempty"`
}

type Executions struct {
	Execution []*Execution `xml:"execution,omitempty"`
}

type Execution struct {
	ID            string      `xml:"id,omitempty"`
	Phase         string      `xml:"phase,omitempty"`
	Goals         string      `xml:"goals,omitempty"`
	Inherited     string      `xml:"inherited,omitempty"`
	Configuration *Properties `xml:"configuration,omitempty"`
}

type Resources struct {
	Resource []*Resource `xml:"resource,omitempty"`
}

type Resource struct {
	TargetPath string `xml:"targetPath,omitempty"`
	Filtering  string `xml:"filtering,omitempty"`
	Directory  string `xml:"directory,omitempty"`
	Includes   string `xml:"includes,omitempty"`
	Excludes   string `xml:"excludes,omitempty"`
}

type TestResources struct {
	TestResource []*Resource `xml:"testResource,omitempty"`
}

type CiManagement struct {
	System    string     `xml:"system,omitempty"`
	URL       string     `xml:"url,omitempty"`
	Notifiers *Notifiers `xml:"notifiers,omitempty"`
}

type Notifiers struct {
	Notifier []*Notifier `xml:"notifier,omitempty"`
}

type Notifier struct {
	Type          string      `xml:"type,omitempty"`
	SendOnError   string      `xml:"sendOnError,omitempty"`
	SendOnFailure string      `xml:"sendOnFailure,omitempty"`
	SendOnSuccess string      `xml:"sendOnSuccess,omitempty"`
	SendOnWarning string      `xml:"sendOnWarning,omitempty"`
	Address       string      `xml:"address,omitempty"`
	Configuration *Properties `xml:"configuration,omitempty"`
}

type Properties struct {
	Entries map[string]string
}

func (p *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	type entry struct {
		XMLName xml.Name
		Key     string `xml:"name,attr"`
		Value   string `xml:",chardata"`
	}
	p.Entries = map[string]string{}
	for {
		var e entry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		p.Entries[e.XMLName.Local] = e.Value
	}
	return nil
}

func (p *Properties) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	tokens := []xml.Token{start}

	for key, value := range p.Entries {
		t := xml.StartElement{Name: xml.Name{Local: key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{Name: t.Name})
	}

	tokens = append(tokens, xml.EndElement{Name: start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	return e.Flush()
}

type Contributors struct {
	Contributor []*Contributor `xml:"contributor,omitempty"`
}

type Contributor struct {
	Name            string      `xml:"name,omitempty"`
	Email           string      `xml:"email,omitempty"`
	URL             string      `xml:"url,omitempty"`
	Organization    string      `xml:"organization,omitempty"`
	OrganizationURL string      `xml:"organizationUrl,omitempty"`
	Roles           string      `xml:"roles,omitempty"`
	Timezone        string      `xml:"timezone,omitempty"`
	Properties      *Properties `xml:"properties,omitempty"`
	ID              *string     `xml:"id,omitempty,omitempty"`
}

type DependencyManagement struct {
	Dependencies *Dependencies `xml:"dependencies,omitempty"`
}

type Developers struct {
	Developer []*Contributor `xml:"developer,omitempty"`
}

type DistributionManagement struct {
	Repository         *Repository `xml:"repository,omitempty"`
	SnapshotRepository *Repository `xml:"snapshotRepository,omitempty"`
	Site               *Site       `xml:"site,omitempty"`
	DownloadURL        string      `xml:"downloadUrl,omitempty"`
	Relocation         *Parent     `xml:"relocation,omitempty"`
	Status             string      `xml:"status,omitempty"`
}

type Repository struct {
	UniqueVersion *string   `xml:"uniqueVersion,omitempty,omitempty"`
	Releases      *Releases `xml:"releases,omitempty"`
	Snapshots     *Releases `xml:"snapshots,omitempty"`
	ID            string    `xml:"id,omitempty"`
	Name          string    `xml:"name,omitempty"`
	URL           string    `xml:"url,omitempty"`
	Layout        string    `xml:"layout,omitempty"`
}

type Releases struct {
	Enabled        string `xml:"enabled,omitempty"`
	UpdatePolicy   string `xml:"updatePolicy,omitempty"`
	ChecksumPolicy string `xml:"checksumPolicy,omitempty"`
}

type Site struct {
	ID   string `xml:"id,omitempty"`
	Name string `xml:"name,omitempty"`
	URL  string `xml:"url,omitempty"`
}

type IssueManagement struct {
	System string `xml:"system,omitempty"`
	URL    string `xml:"url,omitempty"`
}

type Licenses struct {
	License []*License `xml:"license,omitempty"`
}

type License struct {
	Name         string `xml:"name,omitempty"`
	URL          string `xml:"url,omitempty"`
	Distribution string `xml:"distribution,omitempty"`
	Comments     string `xml:"comments,omitempty"`
}

type MailingLists struct {
	MailingList []*MailingList `xml:"mailingList,omitempty"`
}

type MailingList struct {
	Name          string `xml:"name,omitempty"`
	Subscribe     string `xml:"subscribe,omitempty"`
	Unsubscribe   string `xml:"unsubscribe,omitempty"`
	Post          string `xml:"post,omitempty"`
	Archive       string `xml:"archive,omitempty"`
	OtherArchives string `xml:"otherArchives,omitempty"`
}

type Organization struct {
	Name string `xml:"name,omitempty"`
	URL  string `xml:"url,omitempty"`
}

type PluginRepositories struct {
	PluginRepository []*Repository `xml:"pluginRepository,omitempty"`
}

type Prerequisites struct {
	Maven string `xml:"maven,omitempty"`
}

type Profiles struct {
	Profile []*Profile `xml:"profile,omitempty"`
}

func (m *Profiles) GetProfileSlice() []*Profile {
	if m != nil {
		return m.Profile
	}
	return nil
}

type Profile struct {
	ID                     string                  `xml:"id,omitempty"`
	Activation             *Activation             `xml:"activation,omitempty"`
	Build                  *ProfileBuild           `xml:"build,omitempty"`
	Modules                string                  `xml:"modules,omitempty"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement,omitempty"`
	Properties             *Properties             `xml:"properties,omitempty"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement,omitempty"`
	Dependencies           *Dependencies           `xml:"dependencies,omitempty"`
	Repositories           *Repositories           `xml:"repositories,omitempty"`
	PluginRepositories     *PluginRepositories     `xml:"pluginRepositories,omitempty"`
	Reports                string                  `xml:"reports,omitempty"`
	Reporting              *Reporting              `xml:"reporting,omitempty"`
}

func (m *Profile) GetBuild() *ProfileBuild {
	if m != nil {
		return m.Build
	}
	return nil
}

type Activation struct {
	ActiveByDefault string    `xml:"activeByDefault,omitempty"`
	JDK             string    `xml:"jdk,omitempty"`
	OS              *OS       `xml:"os,omitempty"`
	Property        *Property `xml:"property,omitempty"`
	File            *File     `xml:"file,omitempty"`
}

type File struct {
	Missing string `xml:"missing,omitempty"`
	Exists  string `xml:"exists,omitempty"`
}

type OS struct {
	Name    string `xml:"name,omitempty"`
	Family  string `xml:"family,omitempty"`
	Arch    string `xml:"arch,omitempty"`
	Version string `xml:"version,omitempty"`
}

type Property struct {
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}

type ProfileBuild struct {
	DefaultGoal      string            `xml:"defaultGoal,omitempty"`
	Resources        *Resources        `xml:"resources,omitempty"`
	TestResources    *TestResources    `xml:"testResources,omitempty"`
	Directory        string            `xml:"directory,omitempty"`
	FinalName        string            `xml:"finalName,omitempty"`
	Filters          string            `xml:"filters,omitempty"`
	PluginManagement *PluginManagement `xml:"pluginManagement,omitempty"`
	Plugins          *Plugins          `xml:"plugins,omitempty"`
}

func (m *ProfileBuild) GetPluginManagement() *PluginManagement {
	if m != nil {
		return m.PluginManagement
	}
	return nil
}

func (m *ProfileBuild) GetPlugins() *Plugins {
	if m != nil {
		return m.Plugins
	}
	return nil
}

type Reporting struct {
	ExcludeDefaults string            `xml:"excludeDefaults,omitempty"`
	OutputDirectory string            `xml:"outputDirectory,omitempty"`
	Plugins         *ReportingPlugins `xml:"plugins,omitempty"`
}

type ReportingPlugins struct {
	Plugin []*ReportingPlugin `xml:"plugin,omitempty"`
}

type ReportingPlugin struct {
	GroupID       string      `xml:"groupId,omitempty"`
	ArtifactID    string      `xml:"artifactId,omitempty"`
	Version       string      `xml:"version,omitempty"`
	ReportSets    *ReportSets `xml:"reportSets,omitempty"`
	Inherited     string      `xml:"inherited,omitempty"`
	Configuration *Properties `xml:"configuration,omitempty"`
}

type ReportSets struct {
	ReportSet []*ReportSet `xml:"reportSet,omitempty"`
}

type ReportSet struct {
	ID            string      `xml:"id,omitempty"`
	Reports       string      `xml:"reports,omitempty"`
	Inherited     string      `xml:"inherited,omitempty"`
	Configuration *Properties `xml:"configuration,omitempty"`
}

type Repositories struct {
	Repository []*Repository `xml:"repository,omitempty"`
}

type SCM struct {
	Connection          string `xml:"connection,omitempty"`
	DeveloperConnection string `xml:"developerConnection,omitempty"`
	Tag                 string `xml:"tag,omitempty"`
	URL                 string `xml:"url,omitempty"`
}
