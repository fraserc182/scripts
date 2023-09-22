package timestream

import (
	"fmt"

	"github.com/spf13/cobra"
)

// timestreamCmd represents the timestream command
var TimestreamCmd = &cobra.Command{
	Use:   "timestream",
	Short: "A brief description of your command",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("timestream called")
	},
}

func init() {
}
