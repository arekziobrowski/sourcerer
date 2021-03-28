package dependency

type SystemMavenDownloader struct {
	workingDirectory string
}

func NewSystemMavenDownloader(wd string) *SystemMavenDownloader {
	return &SystemMavenDownloader{
		workingDirectory: wd,
	}
}

func (g *SystemMavenDownloader) Get(src string) error {
	panic("implement me")
}
