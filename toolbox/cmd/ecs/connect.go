package ecs

import (
	"context"
	"fmt"
	"time"

	u "toolbox/utils"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

var (
	// used for flags
	cluster   string
	task      string
	container string
	command   string
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "used to connect to a running container",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ecsConnect(cluster, task, container, command)
	},
}

func ecsConnect(cluster string, task string, container string, command string) {
	// Load the AWS configuration from the user's environment.
	cfg, err := u.LoadAWSConfig()
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v", err)
		return
	}

	svc := ecs.NewFromConfig(cfg)

	// Create a context with a timeout that will cancel the request if it takes
	// more than the passed in timeout.
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	res, err := svc.ExecuteCommand(ctx, &ecs.ExecuteCommandInput{
		Cluster:   &cluster,
		Task:      &task,
		Container: &container,
		Command:   &command,
	})
	_ = res
	if err != nil {
		fmt.Printf("Failed to execute command on container %s: %v\n", container, err)
		return
	}

}

func init() {
	EcsCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringVarP(&cluster, "cluster", "a", "", "cluster name")
	connectCmd.Flags().StringVarP(&task, "task", "t", "", "task id")
	connectCmd.Flags().StringVarP(&container, "container", "b", "", "container name")
	connectCmd.Flags().StringVarP(&command, "command", "c", "/bin/bash", "command to run")

}
