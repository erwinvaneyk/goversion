//go:generate goversion fields --goversion "" --pkg version -o version.gen.go
package version

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

const (
	GitTreeStateDirty = "dirty"
	GitTreeStateClean = "clean"
)

var (
	versionInfoMu = &sync.RWMutex{}
	PackageName   = reflect.TypeOf(versionInfo).PkgPath()
	versionInfo   Info
)

// TODO add version parsing and comparing.
type Info struct {
	Version           string `json:"version"`
	BuildDate         string `json:"buildDate"`
	BuildPlatformArch string `json:"buildPlatformArch"`
	BuildPlatformOS   string `json:"buildPlatformOS"`
	BuildBy           string `json:"buildBy"`
	GoVersion         string `json:"goVersion"`
	GitCommit         string `json:"gitCommit"`
	GitTreeState      string `json:"gitTreeState"`
}

func (i Info) IsEmpty() bool {
	infoType := reflect.ValueOf(i)
	for i := 0; i < infoType.NumField(); i++ {
		if !infoType.Field(i).IsZero() {
			return false
		}
	}
	return true
}

func (i Info) String() string {
	return i.ToJson()
}

func (i Info) ToJson() string {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	return string(bs)
}

func (i Info) ToPrettyJson() string {
	bs, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bs)
}

func Set(updatedVersion Info) {
	if updatedVersion.IsEmpty() {
		return
	}
	versionInfoMu.Lock()
	defer versionInfoMu.Unlock()
	versionInfo = updatedVersion
}

func Get() Info {
	versionInfoMu.RLock()
	defer versionInfoMu.RUnlock()
	return versionInfo
}

func (i Info) GenerateLDFlags(pkg string) string {
	var flags string
	if i.BuildBy != "" {
		flags += generateLDFlag(pkg, "buildBy", i.BuildBy) + " "
	}
	if i.BuildDate != "" {
		flags += generateLDFlag(pkg, "buildDate", i.BuildDate) + " "
	}
	if i.BuildPlatformArch != "" {
		flags += generateLDFlag(pkg, "buildPlatformArch", i.BuildPlatformArch) + " "
	}
	if i.BuildPlatformOS != "" {
		flags += generateLDFlag(pkg, "buildPlatformOS", i.BuildPlatformOS) + " "
	}
	if i.GitCommit != "" {
		flags += generateLDFlag(pkg, "gitCommit", i.GitCommit) + " "
	}
	if i.GitTreeState != "" {
		flags += generateLDFlag(pkg, "gitTreeState", i.GitTreeState) + " "
	}
	if i.GoVersion != "" {
		flags += generateLDFlag(pkg, "goVersion", i.GoVersion)
	}
	if i.Version != "" {
		flags += generateLDFlag(pkg, "version", i.Version) + " "
	}
	return flags
}

func generateLDFlag(pkg string, field string, val string) string {
	return fmt.Sprintf("-X \"%s.%s=%s\"", pkg, field, val)
}
