//go:generate goversion generate --goversion "" --pkg goversion -o version.gen.go
package goversion

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/ghodss/yaml"
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

// Info contains all the version-related information.
//
// TODO add version parsing and comparing.
type Info struct {
	// Version is the semantic version of the application.
	Version string `json:"version"`

	// BuildDate contains the date and time when the binary was built.
	BuildDate string `json:"buildDate"`

	// BuildArch is the system architecture that was used to build the binary.
	BuildArch string `json:"buildArch"`

	// BuildOS is the operating system that was used to build the binary.
	BuildOS string `json:"buildOS"`

	// BuildBy is a free-form field that contains info about who or what was responsible for the build.
	BuildBy string `json:"buildBy"`

	// GoVersion the go version that was used to build the binary.
	GoVersion string `json:"goVersion"`

	// GitCommit is the HEAD commit at the moment of building
	GitCommit string `json:"gitCommit"`

	// GitTreeState indicates whether there where uncommitted changes when the binary was built.
	//
	// If there uncommitted changes, this field will be "dirty". Otherwise, if
	// there are no uncommitted changes, this field will be "clean".
	GitTreeState string `json:"gitTreeState"`
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
	return i.ToJSON()
}

func (i Info) ToJSON() string {
	bs, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	return string(bs)
}

func (i Info) ToYAML() string {
	bs, err := yaml.Marshal(i)
	if err != nil {
		panic(err)
	}

	return string(bs)
}

func (i Info) ToPrettyJSON() string {
	bs, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bs)
}

func (i Info) ToLDFlags(pkg string) string {
	var flags string
	if i.BuildBy != "" {
		flags += generateLDFlag(pkg, "buildBy", i.BuildBy) + " "
	}
	if i.BuildDate != "" {
		flags += generateLDFlag(pkg, "buildDate", i.BuildDate) + " "
	}
	if i.BuildArch != "" {
		flags += generateLDFlag(pkg, "buildArch", i.BuildArch) + " "
	}
	if i.BuildOS != "" {
		flags += generateLDFlag(pkg, "buildOS", i.BuildOS) + " "
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

func generateLDFlag(pkg string, field string, val string) string {
	return fmt.Sprintf("-X \"%s.%s=%s\"", pkg, field, val)
}
