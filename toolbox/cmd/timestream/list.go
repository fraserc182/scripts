package timestream

import (
	"context"
	"fmt"
	"time"
	u "toolbox/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists timestream databases in the given region",
	Long: `Pass the region as an argument to the command. For example: -r eu-west-1.

Timestream is available in the following regions:
- us-east-1
- us-east-2
- us-west-1
- us-west-2
- eu-west-1
- eu-central-1
- ap-southeast-2
- ap-northeast-1`,
	Run: func(cmd *cobra.Command, args []string) {
		listTimestreamDatabase()
	},
}

func listTimestreamDatabase() {
	query := "SHOW DATABASES"

	// Load the AWS configuration from the user's environment.
	cfg, err := u.LoadAWSConfig()
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v", err)
		return
	}

	// create new timestream client from config
	svc := timestreamquery.NewFromConfig(cfg)

	// Create a context with a timeout that will cancel the request if it takes
	// more than the passed in timeout.
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	// create a new query input
	input := &timestreamquery.QueryInput{
		QueryString: aws.String(query),
	}

	// query the enriched-gps table
	output, err := svc.Query(ctx, input)
	if err != nil {
		fmt.Printf("Failed to execute query: %v\n", err)
		return
	}

	for _, row := range output.Rows {
		for _, data := range row.Data {
			fmt.Printf("%v\n", *data.ScalarValue)
		}
	}
	
}

func init() {
	TimestreamCmd.AddCommand(listCmd)
}
