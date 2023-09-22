package ec2

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
	u "toolbox/utils"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

var (
	// used for flags
	id string
)
// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connects to specifed ec2 instance over SSM",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		ec2Connect(id)
	},
}

func ec2Connect(id string) {
	// Load the AWS configuration from the user's environment.
	cfg, err := u.LoadAWSConfig()
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v", err)
		return
	}

	// Create new SSM client from the provided config.
	svc := ssm.NewFromConfig(cfg)

	// Create a context with a timeout that will cancel the request if it takes
	// more than the passed in timeout.
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	res, err := svc.StartSession(ctx, &ssm.StartSessionInput{
		Target: &id,
		})
	if err != nil {
		fmt.Printf("Failed to start SSM session with instance %s: %v\n", id, err)
		return
	}


	// Print the session ID to the console.
	fmt.Printf("Started session with instance %s. Session ID: %s\n", id, *res.SessionId)

	// Start a local shell session.
	cmd := exec.Command("bash", "-i")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Terminate the SSM session when the local shell session is terminated.
	svc.TerminateSession(ctx, &ssm.TerminateSessionInput{
		SessionId: res.SessionId,
	})
	fmt.Printf("Terminated session with instance %s. Session ID: %s\n", id, *res.SessionId)
}



func init() {
	connectCmd.Flags().StringVarP(&id, "instance", "i", "", "instance ID")
	Ec2Cmd.AddCommand(connectCmd)
}