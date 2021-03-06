package cmd

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/erwinvaneyk/cobras"
	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/goversion"
)

type LDFlagsOptions struct {
	Version     goversion.Info
	PackageName string
	PrintLDFlag bool
	Strict      bool
}

func NewCmdLDFlags() *cobra.Command {
	opts := &LDFlagsOptions{
		PackageName: goversion.PackageName,
		PrintLDFlag: true,
	}

	cmd := &cobra.Command{
		Use:   "ldflags",
		Short: "Collect version information from host and print as valid Go ldflags.",
		Example: `
# Print ldflags with a version.
goversion ldflags --version v1.0.4

# Print ldflags with a custom package prefix. This is useful when you have 
# generated a custom version file. See 'goversion generate'.
goversion ldflags --version v1.2.3 --pkg main
		
# Print ldflags in strict mode.
goversion ldflags 					# Should fail, because it is missing the version.
goversion ldflags --version v1.0.1 	# Should succeed.`,
		Run: cobras.Run(opts),
	}

	cmd.Flags().StringVarP(&opts.PackageName, "pkg", "p", opts.PackageName, "The Go package that should be used in the ldflags.")
	cmd.Flags().BoolVar(&opts.PrintLDFlag, "print-ldflag", opts.PrintLDFlag, "If set, the flags will be wrapped with the '-ldflags' flag.")
	cmd.Flags().BoolVar(&opts.Strict, "strict", opts.Strict, "If set, goversion will validate the version info before writing the ldflags.")

	versionInfoVal := reflect.ValueOf(&opts.Version)
	for i := 0; i < versionInfoVal.Elem().NumField(); i++ {
		fieldVal := versionInfoVal.Elem().Field(i)
		fieldType := versionInfoVal.Elem().Type().Field(i)
		cmd.Flags().StringVar(fieldVal.Addr().Interface().(*string), camelCaseToKebabCase(fieldType.Name), fieldVal.String(), fmt.Sprintf("Manually set the '%s' field in the version info.", fieldType.Name))
	}

	return cmd
}

func (o *LDFlagsOptions) Complete(cmd *cobra.Command, args []string) error {
	o.Version = goversion.AugmentFromEnv(o.Version)
	return nil
}

func (o *LDFlagsOptions) Validate() error {
	if o.Strict {
		return goversion.ValidateStrict(o.Version)
	}
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
