package extensions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/go-version"
)

func NewCobraCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Get().ToPrettyJson())
		},
	}



	return cmd
}
