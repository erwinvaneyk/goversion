package extensions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/erwinvaneyk/goversion"
)

func NewCobraCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(goversion.Get().ToPrettyJson())
		},
	}

	return cmd
}
