package common

var (
	Version   string
	BuildTime string
	Commit    string
)

func init() {
	if Version == "" {
		Version = "unknown"
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	if Commit == "" {
		Commit = "unknown"
	}
}
