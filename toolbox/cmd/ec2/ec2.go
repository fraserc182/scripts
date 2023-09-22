package ec2

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Ec2Cmd represents the ec2 command
var Ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Commands to interact with ec2 instances",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ec2 called")
	},
}

func init() {
}
