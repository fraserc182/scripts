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

var (
	// used for flags
	host string
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Queries timestream to ensure there is data present",
	Long: `Runs the following queries against the target Timestream database:

	A successful query will return a count of the number of records in the table.
	If the count is 0, then there is no data present in the table.

	SELECT BIN(time, 1m) AS timestamp, COUNT(*) as results
	FROM "HOST"."enriched-gps" 
	WHERE time between ago(30m) and ago(2m)
	GROUP BY BIN(time, 1m)
	order by timestamp desc

	SELECT BIN(time, 1m) AS timestamp, COUNT(*) as results
	FROM "HOST"."actuals" 
	WHERE time between ago(30m) and ago(2m)
	GROUP BY BIN(time, 1m)
	order by timestamp desc`,
	Run: func(cmd *cobra.Command, args []string) {
		queryTimestream(host)
	},
}

func queryTimestream(host string) {
	queryEnrichedGps := fmt.Sprintf("SELECT BIN(time, 1m) AS timestamp, COUNT(*) as results\nFROM \"%s\".\"enriched-gps\" \nWHERE time between ago(30m) and ago(2m)\nGROUP BY BIN(time, 1m)\norder by timestamp desc", host)
	queryActuals := fmt.Sprintf("SELECT BIN(time, 1m) AS timestamp, COUNT(*) as results\nFROM \"%s\".\"actuals\" \nWHERE time between ago(30m) and ago(2m)\nGROUP BY BIN(time, 1m)\norder by timestamp desc", host)
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

	// create channels for the query results
	enrichedGPSResultCh := make(chan *timestreamquery.QueryOutput)
	actualsResultCh := make(chan *timestreamquery.QueryOutput)

	// create a goroutine for each query
	go func() {
		// create a new query input
		input := &timestreamquery.QueryInput{
			QueryString: aws.String(queryEnrichedGps),
		}

		// query the enriched-gps table
		output, err := svc.Query(ctx, input)
		if err != nil {
			fmt.Printf("Failed to execute query: %v\n", err)
			return
		}

		// send the result through the channel
		enrichedGPSResultCh <- output
	}()

	go func() {
		// create a new query input
		input := &timestreamquery.QueryInput{
			QueryString: aws.String(queryActuals),
		}

		// query the actuals table
		output, err := svc.Query(ctx, input)
		if err != nil {
			fmt.Printf("Failed to execute query: %v\n", err)
			return
		}

		// send the result through the channel
		actualsResultCh <- output
	}()

	// create a goroutine to receive the results from the channels and print them out
	go func() {
		// receive the enriched-gps query result
		enrichedGPSResult := <-enrichedGPSResultCh

		// print out the enriched-gps query results
		fmt.Println("Enriched GPS Timestamp\tResults")
		for _, row := range enrichedGPSResult.Rows {
			fmt.Printf("%v\t%s\n", *row.Data[0].ScalarValue, *row.Data[1].ScalarValue)
		}

		// receive the actuals query result
		actualsResult := <-actualsResultCh

		// print out the actuals query results
		fmt.Println("Actuals Timestamp\tResults")
		for _, row := range actualsResult.Rows {
			fmt.Printf("%v\t%s\n", *row.Data[0].ScalarValue, *row.Data[1].ScalarValue)
		}
	}()

	// wait for the goroutines to finish
	time.Sleep(5 * time.Second)
}

func init() {
	queryCmd.Flags().StringVarP(&host, "host", "o", "", "timestream database to query")
	TimestreamCmd.AddCommand(queryCmd)
}
