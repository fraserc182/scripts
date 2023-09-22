package ec2

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"runtime"
	"strings"
	"sync"

	co "toolbox/config"
	u "toolbox/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type InstanceInfo struct {
	ClusterName string
	InstanceID  string
	IPAddress   string
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns all ec2 instances and optionally allows you to connect to one",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		listInstances()
	},
}

func parseInstanceInfo(s string) InstanceInfo {
    // this function takes what is returned from the ec2.DescribeInstances API call
    // and parses it into a struct
	slice := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '(' || r == ')'
	})
	ip := strings.Trim(slice[2], "[]")
	info := InstanceInfo{
		ClusterName: slice[0],
		InstanceID:  slice[1],
		IPAddress:   ip,
	}
	return info
}

// listInstances lists the running EC2 instances in the specified region.
func listInstances() {
	// Set the number of CPUs to use for parallelism.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load the AWS configuration from the user's environment.
	cfg, err := u.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Create an EC2 client using the loaded configuration.
	svc := ec2.NewFromConfig(cfg)

	// Describe the regions matching the specified region name.
	regionsResp, err := svc.DescribeRegions(context.Background(), &ec2.DescribeRegionsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []string{co.AwsRegion},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to describe regions: %v", err)
	}

	// Create a channel to receive the list of instances.
	instancesChan := make(chan string)

	// Create a wait group to synchronize the goroutines that will describe the instances.
	var wg sync.WaitGroup

	// For each region returned by DescribeRegions, start a goroutine to describe the instances in that region.
	for _, region := range regionsResp.Regions {
		// Increment the wait group counter.
		wg.Add(1)

		go func(region string) {
			// Decrement the wait group counter when the goroutine completes.
			defer wg.Done()

			// Describe the running instances in the region.
			resp, err := svc.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
				Filters: []types.Filter{
					{
						Name: aws.String("instance-state-name"),
						Values: []string{
							"running",
						},
					},
				},
			})
			if err != nil {
				log.Printf("Failed to describe instances in region %s: %v", region, err)
				return
			}

			// For each instance returned by DescribeInstances, extract its name, instance ID, and private IP address.
			for _, reservation := range resp.Reservations {
				for _, instance := range reservation.Instances {
					name := "None"
					for _, tag := range instance.Tags {
						if aws.ToString(tag.Key) == "Name" {
							name = url.QueryEscape(aws.ToString(tag.Value))
							break
						}
					}
					instancesChan <- fmt.Sprintf("%s (%s) [%s]",
						name,
						aws.ToString(instance.InstanceId),
						aws.ToString(instance.PrivateIpAddress))
				}
			}
		}(aws.ToString(region.RegionName))
	}

	// Start a goroutine to wait for all the goroutines that describe instances to complete, and then close the channel.
	go func() {
		wg.Wait()
		close(instancesChan)
	}()

	// Read the list of instances from the channel and append them to a slice.
	instances := make([]string, 0)
	for inst := range instancesChan {
		instances = append(instances, inst)
	}

	// Display the list of instances to the user.
	promptInstances(instances)
}

func promptInstances(instances []string) {
	prompt := promptui.Select{
		Label: "Select an instance",
		Items: instances,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt cancelled %v\n", err)
		return
	}

	info := parseInstanceInfo(result)
	ec2Connect(info.InstanceID)
}

func init() {
	Ec2Cmd.AddCommand(listCmd)

}
