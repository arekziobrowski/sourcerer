package source

type MavenDownloader struct {
	workingDirectory string
}

func NewMavenDownloader(wd string) *MavenDownloader {
	return &MavenDownloader{
		workingDirectory: wd,
	}
}

func (g *MavenDownloader) Get(src string) error {
	panic("implement me")
}
