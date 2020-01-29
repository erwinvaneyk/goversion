package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/go-version"
)

type LDFlagsOptions struct {
	Version     version.Info
	PackageName string
}

func NewCmdLDFlags() *cobra.Command {
	opts := &LDFlagsOptions{
		PackageName: version.PackageName,
	}

	cmd := &cobra.Command{
		Use: "ldflags",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := opts.Complete(args)
			if err != nil {
				return err
			}
			return opts.Run()
		},
	}

	cmd.Flags().StringVar(&opts.PackageName, "pkg", opts.PackageName, "")
	cmd.Flags().StringVar(&opts.Version.Version, "version", opts.Version.Version, "")
	cmd.Flags().StringVar(&opts.Version.GitCommit, "git-commit", opts.Version.GitCommit, "")
	cmd.Flags().StringVar(&opts.Version.BuildDate, "build-date", opts.Version.BuildDate, "")

	return cmd
}

func (o *LDFlagsOptions) Complete(args []string) error {
	// Infer build by
	if o.Version.BuildBy == "" {
		out, err := exec.Command("git", "config", "user.name").CombinedOutput()
		if err == nil {
			o.Version.BuildBy = strings.TrimSpace(string(out))
		}

		out, err = exec.Command("git", "config", "user.email").CombinedOutput()
		if err == nil {
			o.Version.BuildBy += fmt.Sprintf(" (%s)", strings.TrimSpace(string(out)))
		}

		o.Version.BuildBy = strings.TrimSpace(o.Version.BuildBy)
	}

	// Infer the build date
	if o.Version.BuildDate == "" {
		o.Version.BuildDate = time.Now().UTC().Format(time.RFC3339)
	}

	// Infer build platform OS
	if o.Version.BuildPlatformOS == "" {
		out, err := exec.Command("uname").CombinedOutput()
		if err == nil {
			o.Version.BuildPlatformOS = strings.TrimSpace(string(out))
		}
	}

	// Infer build platform architecture
	if o.Version.BuildPlatformArch == "" {
		out, err := exec.Command("uname", "-m").CombinedOutput()
		if err == nil {
			o.Version.BuildPlatformArch = strings.TrimSpace(string(out))
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
				o.Version.GitTreeState = version.GitTreeStateClean
			} else {
				o.Version.GitTreeState = version.GitTreeStateDirty
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

func (o *LDFlagsOptions) Run() error {
	fmt.Printf("-ldflags '%s'", o.Version.GenerateLDFlags(o.PackageName))

	return nil
}