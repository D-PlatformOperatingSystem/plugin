/*Package commands implement dapp client commands*/
package commands

import (
	"github.com/spf13/cobra"
)

/*
 *
 */

// Cmd storage client command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "storage command",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
	//add sub command
	)
	return cmd
}
