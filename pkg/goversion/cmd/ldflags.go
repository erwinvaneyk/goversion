package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/erwinvaneyk/cobras"
	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/goversion"
)

type LDFlagsOptions struct {
	Version     goversion.Info
	PackageName string
	PrintLDFlag bool
}

func NewCmdLDFlags() *cobra.Command {
	opts := &LDFlagsOptions{
		PackageName: goversion.PackageName,
		PrintLDFlag: true,
	}

	cmd := &cobra.Command{
		Use:   "ldflags",
		Short: "Collect version information from host and print as valid Go ldflags.",
		Run:   cobras.Run(opts),
	}

	cmd.Flags().StringVar(&opts.PackageName, "pkg", opts.PackageName, "The Go package that should be used in the ldflags.")
	cmd.Flags().BoolVar(&opts.PrintLDFlag, "print-ldflag", opts.PrintLDFlag, "If set, the flags will be wrapped with the '-ldflags' flag")
	// TODO strict mode

	versionInfoVal := reflect.ValueOf(&opts.Version)
	for i := 0; i < versionInfoVal.Elem().NumField(); i++ {
		fieldVal := versionInfoVal.Elem().Field(i)
		fieldType := versionInfoVal.Elem().Type().Field(i)
		cmd.Flags().StringVar(fieldVal.Addr().Interface().(*string), camelCaseToKebabCase(fieldType.Name), fieldVal.String(), fmt.Sprintf("Manually set the '%s' field in the version info.", fieldType.Name))
	}

	return cmd
}

func (o *LDFlagsOptions) Complete(cmd *cobra.Command, args []string) error {
	// Infer build by
	if o.Version.BuildBy == "" {
		out, err := exec.Command("git", "config", "user.name").CombinedOutput()
		if err == nil {
			o.Version.BuildBy = strings.TrimSpace(string(out))
		}

		// TODO
		// out, err = exec.Command("git", "config", "user.email").CombinedOutput()
		// if err == nil {
		// 	o.Version.BuildBy += fmt.Sprintf(" (%s)", strings.TrimSpace(string(out)))
		// }

		o.Version.BuildBy = strings.TrimSpace(o.Version.BuildBy)
	}

	// Infer the build date
	if o.Version.BuildDate == "" {
		o.Version.BuildDate = time.Now().UTC().Format(time.RFC3339)
	}

	// Infer build platform OS
	if o.Version.BuildOS == "" {
		out, err := exec.Command("uname").CombinedOutput()
		if err == nil {
			o.Version.BuildOS = strings.TrimSpace(string(out))
		}
	}

	// Infer build platform architecture
	if o.Version.BuildArch == "" {
		out, err := exec.Command("uname", "-m").CombinedOutput()
		if err == nil {
			o.Version.BuildArch = strings.TrimSpace(string(out))
		}
	}

	// Infer the git commit
	if o.Version.GitCommit == "" {
		out, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
		if err == nil {
			o.Version.GitCommit = strings.TrimSpace(string(out))
		}
	}

	// Infer the git branch
	if o.Version.GitBranch == "" {
		out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
		if err == nil {
			o.Version.GitBranch = strings.TrimSpace(string(out))
		}
	}

	// Infer the git commit date
	if o.Version.GitCommitDate == "" && o.Version.GitCommit != "" {
		out, err := exec.Command("git", "show", "-s", "--format=%ci", o.Version.GitCommit).CombinedOutput()
		if err == nil {
			o.Version.GitCommitDate = strings.TrimSpace(string(out))
		}
	}

	// Infer git status
	if o.Version.GitTreeState == "" {
		out, err := exec.Command("git", "diff", "--quiet").CombinedOutput()
		if len(out) == 0 {
			if err == nil {
				o.Version.GitTreeState = goversion.GitTreeStateClean
			} else {
				o.Version.GitTreeState = goversion.GitTreeStateDirty
			}
		}
	}

	// Infer go version
	if o.Version.GoVersion == "" {
		out, err := exec.Command("go", "version").CombinedOutput()
		if err == nil {
			o.Version.GoVersion = strings.Split(strings.TrimSpace(string(out)), " ")[2]
		}
	}

	return nil
}

func (o *LDFlagsOptions) Validate() error {
	return nil
}

func (o *LDFlagsOptions) Run(ctx context.Context) error {
	if o.PrintLDFlag {
		fmt.Print("-ldflags '")
	}
	fmt.Printf("%s", o.Version.ToLDFlags(o.PackageName))
	if o.PrintLDFlag {
		fmt.Print("'")
	}

	return nil
}

// camelCaseToKebabCase converts a camelCaseVariableName or PascalCaseVariableName to its kebab-case-variable-name equivalent.
func camelCaseToKebabCase(pascalCaseName string) string {
	if pascalCaseName == "" {
		return ""
	}
	kebabCaseName := []rune(strings.ToLower(pascalCaseName[0:1]))
	var prevCharWasUppercase bool
	for _, c := range []rune(pascalCaseName[1:]) {
		if c >= 'A' && c <= 'Z' {
			if !prevCharWasUppercase {
				kebabCaseName = append(kebabCaseName, '-')
				prevCharWasUppercase = true
			}
			kebabCaseName = append(kebabCaseName, c+'a'-'A')
		} else {
			prevCharWasUppercase = false
			kebabCaseName = append(kebabCaseName, c)
		}
	}
	return string(kebabCaseName)
}
