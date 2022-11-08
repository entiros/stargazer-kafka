package config

import (
	"context"
	awsdk "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
)

func GetAWSConfig(ctx context.Context) (awsdk.Config, error) {

	cfg, err := awscfg.LoadDefaultConfig(ctx)
	if err != nil {
		//log.Logger.Errorf("Failed to read AWS config: %v", err)
		return awsdk.Config{}, err
	}

	return cfg, nil

}
