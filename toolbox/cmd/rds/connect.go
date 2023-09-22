package rds

import (
	"context"
	"fmt"

	"os"
	"os/exec"
	"syscall"

	co "toolbox/config"
	u "toolbox/utils"

	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// used for flags
	host   string
	dbName string
	user   string
)


func dbConnect(host string, dbName string, region string, user string) error {
    awsConfigObj,_ := u.LoadAWSConfig()
		
	// get RDS authentication token
	// https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/rds/
	authenticationToken, err := auth.BuildAuthToken(
		context.TODO(),
		host+":5432", // Database Endpoint (With Port)
		region,       // AWS Region
		user,         // user to connect with
		awsConfigObj.Credentials,
	)
	if err != nil {
		return fmt.Errorf("failed to create authentication token: %w", err)
	}

	// Go requires full path of executable
	// we use LookPath to get it
	binary, err := exec.LookPath("psql")
	if err != nil {
		return fmt.Errorf("failed to find psql binary: %w", err)
	}
	// pass in the arguments
	formattedString := fmt.Sprintf("host=%s port=5432 dbname=%s user=%s password=%s sslmode=require", host, dbName, user, authenticationToken)
	// slice of strings is required for syscall.Exec
	args := []string{"psql", formattedString}

	// syscall.Exec will end the current Go process and pass to the command it calls
	// we used this as psql is opening an interactive prompt
	// we don't want to keep the Go process running in the background
	fmt.Printf("Connecting to %s/%s as %s\n", host, dbName, user)
	error := syscall.Exec(binary, args, os.Environ())
	// We don't expect this to ever return; if it does something is really wrong
	panic(error)

}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Helper for connecting to RDS DBs with IAM credentials",
	Long:  `Connects to the specified RDS host and DB as either the read-write or read-only user and gives an interactive postgres prompt.`,
	// RunE just returns an error as well, so we can do some error checking
	RunE: func(cmd *cobra.Command, args []string) error {
		// set variables if flag hasn't been passed
		host = viper.GetViper().GetString("hosts."+host)
		user = viper.GetViper().GetString("defaults.user")

		dbConnect(host, dbName, co.AwsRegion, user)
		return nil
	},
}

func init() {
	// set flags and mark as required where neccessary
	connectCmd.Flags().StringVarP(&host, "host", "o", "", "rds endpoint")
	connectCmd.MarkFlagRequired("host")
	connectCmd.Flags().StringVarP(&dbName, "database", "d", "", "database name")
	connectCmd.MarkFlagRequired("database")
	connectCmd.Flags().StringVarP(&user, "username", "u", "read-only", "username")

	// create relationship between viper and flags
	viper.BindPFlag("hosts."+host, connectCmd.Flags().Lookup("host"))
	viper.BindPFlag("defaults.user", connectCmd.Flags().Lookup("user"))

	RdsCmd.AddCommand(connectCmd)
}
