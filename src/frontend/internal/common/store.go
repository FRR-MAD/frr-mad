package common

var (
	DaemonVersion = "unknown1"
	TUIVersion    = "unknown1"
	GitCommit     = "unknown1"
	BuildDate     = "unknown1"
	RepoURL       = "https://github.com/frr-mad/frr-mad"
)

type AppVersionInfo struct {
	DaemonVersion string
	TUIVersion    string
	GitCommit     string
	BuildDate     string
	RepoURL       string
}

func SetAppVersionInfo(
	daemonVersion string,
	tuiVersion string,
	gitCommit string,
	buildDate string,
	repoURL string,
) {
	DaemonVersion = daemonVersion
	TUIVersion = tuiVersion
	GitCommit = gitCommit
	BuildDate = buildDate
	RepoURL = repoURL
}

func GetAppVersionInfo() AppVersionInfo {
	return AppVersionInfo{
		DaemonVersion: DaemonVersion,
		TUIVersion:    TUIVersion,
		GitCommit:     GitCommit,
		BuildDate:     BuildDate,
		RepoURL:       RepoURL,
	}
}
