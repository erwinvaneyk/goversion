# Go-version

Simplify versioning of Go applications.

The project consists out of two parts:
1. The version library `github.com/erwinvaneyk/goversion` for providing version
   info and corresponding operations.
2. The _optional_ `goversion` CLI to simplify generating ldflags and versioning
   fields.
   
## Installation

TODO go get

TODO section

To install the optional `goversion` CLI, use one of the following options:

### Using a pre-built release

TODO

### Using Go tools



### Using Go install (deprecated)
```bash
GO11MODULE=off go install github.com/erwinvaneyk/goversion/cmd/goversion 
```

Note: ensure that you have added your `GOBIN` directory to your `PATH`.

## Usage

There are two ways to make use of `goversion` in your project. Both assume 
that you have added the package to your module or GOPATH. If you haven't:
 
```bash
go get github.com/erwinvaneyk/goversion
```

### Using the imported package.
To use `goversion`, import and use the package somewhere in your application.  
For example:

```go
package main

import (
	"fmt"

	"github.com/erwinvaneyk/goversion"
)

func main() {
	fmt.Println(goversion.Get())
}
```


```bash
# With goversion:
go run $(goversion ldflags --version v1.0.1) ./simple

# Or, manually:
go run -ldflags ' \
		-X "github.com/erwinvaneyk/goversion.version=v1.0.1" \
		-X "github.com/erwinvaneyk/goversion.gitCommit=$(git rev-parse HEAD)" \
		-X "github.com/erwinvaneyk/goversion.buildDate=$(date)"' \
	    ./simple
```

### Using generated fields
Using the package does require a long package name to be added, and is not that
extensible. So, if you have `goversion` installed, you could also generate the
ldflag fields in your main package:

```bash
goversion generate -o /path/to/your/main/package/version.gen.go
``` 

Or use the go tooling for the generation process, add the following:

```go
//go:generate goversion generate -o version.gen.go
package main

// ...
```

And run

```bash
go generate
```

Both options will generate a file with the versioninfo fields in the main 
package: 

```go
// Generated by goversion
package main

import github.com/erwinvaneyk/goversion

// The following variables should be filled with goversion ldflags 
var (
	version   string
	gitCommit string
	buildDate string
	goVersion string
)

func init() {
	goversion.Set(goversion.Info{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: goVersion,
	})
}
```

Using the generated code, you can now set the ldflags using the shorter 
`main` instead of the full package name:
 
```bash
# With goversion:
go run $(goversion ldflags --pkg main --version v1.0.4) ./generated 

# Or, manually:
go run -ldflags ' \
    		-X "main.version=v1.0.3" \
    		-X "main.gitCommit=$(git rev-parse HEAD)" \
    		-X "main.buildDate=$(date)"' \
    		./generated
```

See the [examples](./examples) for complete, functioning examples.

### Using Cobra
This project contains a default command to include into your 
[Cobra](https://github.com/spf13/cobra) based CLI. Just add the following
line to the setup of your root command:   

```go
import (
	goversionext "github.com/erwinvaneyk/goversion/pkg/extensions"
)

func init() {
    // ...
	cmd.AddCommand(goversionext.NewCobraCmd())
    // ...
}
```

### Reproducible Builds

To make your builds, checksums, and signatures reproducible, you will need to 
make the following modifications when generating the ldflags:
- Manually set the `--build-date` to a specific date and time at which the build
  should be done.