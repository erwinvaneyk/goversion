// Generated by goversion v0.1.2-SNAPSHOT
package goversion

// The following variables should be filled with goversion ldflags
var (
	version       string
	buildDate     string
	buildArch     string
	buildOS       string
	buildBy       string
	goVersion     string
	gitCommit     string
	gitCommitDate string
	gitBranch     string
	gitTreeState  string
)

func init() {
	Set(Info{
		Version:       version,
		BuildDate:     buildDate,
		BuildArch:     buildArch,
		BuildOS:       buildOS,
		BuildBy:       buildBy,
		GoVersion:     goVersion,
		GitCommit:     gitCommit,
		GitCommitDate: gitCommitDate,
		GitBranch:     gitBranch,
		GitTreeState:  gitTreeState,
	})
}
