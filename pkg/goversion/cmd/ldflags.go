package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/erwinvaneyk/cobras"
	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/goversion"
)

type LDFlagsOptions struct {
	Version     goversion.Info
	PackageName string
}

func NewCmdLDFlags() *cobra.Command {
	opts := &LDFlagsOptions{
		PackageName: goversion.PackageName,
	}

	cmd := &cobra.Command{
		Use: "ldflags",
		Run: cobras.Run(opts),
	}

	cmd.Flags().StringVar(&opts.PackageName, "pkg", opts.PackageName, "The Go package that should be used in the ldflags.")

	cmd.Flags().StringVar(&opts.Version.Version, "version", opts.Version.Version, "")
	cmd.Flags().StringVar(&opts.Version.GitCommit, "git-commit", opts.Version.GitCommit, "")
	cmd.Flags().StringVar(&opts.Version.GitCommit, "git-tree-state", opts.Version.GitTreeState, "")
	cmd.Flags().StringVar(&opts.Version.GitCommit, "go-version", opts.Version.GoVersion, "")
	cmd.Flags().StringVar(&opts.Version.BuildDate, "build-date", opts.Version.BuildDate, "")
	cmd.Flags().StringVar(&opts.Version.BuildDate, "build-by", opts.Version.BuildBy, "")
	cmd.Flags().StringVar(&opts.Version.BuildDate, "build-arch", opts.Version.BuildArch, "")
	cmd.Flags().StringVar(&opts.Version.BuildDate, "build-os", opts.Version.BuildOS, "")

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
			o.Version.GoVersion = strings.TrimSpace(string(out))
		}
	}

	return nil
}

func (o *LDFlagsOptions) Validate() error {
	return nil
}

func (o *LDFlagsOptions) Run(ctx context.Context) error {
	fmt.Printf("-ldflags '%s'", o.Version.ToLDFlags(o.PackageName))

	return nil
}