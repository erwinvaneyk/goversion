package cmd

import (
	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/go-version/extensions"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goversion",
		Short: "Generate linker flags and fields for versioning",
	}

	cmd.AddCommand(NewCmdLDFlags())
	cmd.AddCommand(NewCmdFields())
	cmd.AddCommand(extensions.NewCobraCmd())

	return cmd
}
