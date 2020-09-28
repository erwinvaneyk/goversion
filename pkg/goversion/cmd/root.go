package cmd

import (
	"github.com/erwinvaneyk/cobras"
	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/goversion/pkg/extensions"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goversion",
		Short: "Generate linker flags and fields for versioning your Go applications.",
	}

	cmd.AddCommand(NewCmdLDFlags())
	cmd.AddCommand(NewCmdGenerate())
	cmd.AddCommand(extensions.NewCobraCmdWithDefaults())

	return cmd
}

func Execute() {
	cobras.Execute(New())
}
