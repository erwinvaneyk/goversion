package cmd

import (
	"github.com/erwinvaneyk/cobras"
	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/goversion/pkg/extensions"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goversion",
		Short: "Generate linker flags and fields for versioning",
	}

	cmd.AddCommand(NewCmdLDFlags())
	cmd.AddCommand(NewCmdFields())
	cmd.AddCommand(extensions.NewCobraCmd())

	return cmd
}

func Execute() {
	cobras.Execute(New())
}
