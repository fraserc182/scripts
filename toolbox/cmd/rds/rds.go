package rds

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RdsCmd represents the rds command
var RdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "For working with RDS DBs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rds called")
	},
}

func init() {
}
