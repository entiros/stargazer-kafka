package msk

import (
	"context"
	"fmt"
	awsdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/entiros/stargazer-kafka/internal/aws"

	"strings"
)

func GetMSKIamBootstrap(ctx context.Context, clusterArn string) ([]string, error) {

	client, err := getMSKClient(ctx)
	if err != nil {
		return []string{}, err
	}

	a, err := client.GetBootstrapBrokers(
		ctx,
		&kafka.GetBootstrapBrokersInput{ClusterArn: awsdk.String(clusterArn)},
	)
	if err != nil {
		//log.Logger.Errorf("Failed to get bootstrap brokers: %v", err)
		return []string{}, fmt.Errorf("failed to get bootstrap brokers from MSK: %v", err)
	}

	return strings.Split(*a.BootstrapBrokerStringPublicSaslIam, ","), nil
}

func getMSKClient(ctx context.Context) (*kafka.Client, error) {

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := aws.GetAWSConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := kafka.NewFromConfig(cfg)
	return client, nil
}
