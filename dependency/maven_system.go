package dependency

type MavenDownloader struct {
	workingDirectory string
}

func NewSystemMavenDownloader(wd string) *MavenDownloader {
	return &MavenDownloader{
		workingDirectory: wd,
	}
}

func (g *MavenDownloader) Get(src string) error {
	panic("implement me")
}
