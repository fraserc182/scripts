package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	co "toolbox/config"
)

func LoadAWSConfig() (aws.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(co.AwsRegion), config.WithSharedConfigProfile(co.AwsProfile))
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return cfg, nil
}

func LoadAWSCredentials(config aws.Config) (*sts.GetCallerIdentityOutput, error) {
	// Load AWS credentials
	creds, err := sts.NewFromConfig(config).GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS credentials: %w", err)
	}

	return creds, nil
}
