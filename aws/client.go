package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	awsbase "github.com/hashicorp/aws-sdk-go-base/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
)

// AwsClient is a wrapper of the AWS SDK client.
// This is the real implementation that satisfies the interface.
type AwsClient struct {
	IAM         *iam.Client
	EC2         *ec2.Client
	RDS         *rds.Client
	ElastiCache *elasticache.Client
	ELB         *elasticloadbalancing.Client
	ELBV2       *elasticloadbalancingv2.Client
	ECS         *ecs.Client
}

var _ Client = (*AwsClient)(nil)

// Credentials is credentials for AWS used in deep check mode
type Credentials struct {
	AccessKey             string
	SecretKey             string
	Profile               string
	CredsFile             string
	AssumeRoleARN         string
	AssumeRoleExternalID  string
	AssumeRolePolicy      string
	AssumeRoleSessionName string
	Region                string
}

// NewClient returns a new Client with configured session
func NewClient(creds Credentials) (Client, error) {
	logger.Info("Initialize AWS Client")

	config, err := getBaseConfig(creds)
	if err != nil {
		return nil, err
	}

	_, awsConfig, diags := awsbase.GetAwsConfig(context.Background(), config)
	for _, diag := range diags.Errors() {
		err = errors.Join(err, fmt.Errorf("%s; %s", diag.Summary(), diag.Detail()))
	}
	if err != nil {
		return nil, err
	}

	return &AwsClient{
		IAM:         iam.NewFromConfig(awsConfig),
		EC2:         ec2.NewFromConfig(awsConfig),
		RDS:         rds.NewFromConfig(awsConfig),
		ElastiCache: elasticache.NewFromConfig(awsConfig),
		ELB:         elasticloadbalancing.NewFromConfig(awsConfig),
		ELBV2:       elasticloadbalancingv2.NewFromConfig(awsConfig),
		ECS:         ecs.NewFromConfig(awsConfig),
	}, nil
}

func getBaseConfig(creds Credentials) (*awsbase.Config, error) {
	expandedCredsFile, err := homedir.Expand(creds.CredsFile)
	if err != nil {
		return nil, err
	}

	config := &awsbase.Config{
		AccessKey:              creds.AccessKey,
		SecretKey:              creds.SecretKey,
		Profile:                creds.Profile,
		Region:                 creds.Region,
		CallerName:             "tflint-ruleset-aws",
		CallerDocumentationURL: "https://github.com/terraform-linters/tflint-ruleset-aws/blob/master/docs/deep_checking.md",
	}

	if creds.AssumeRoleARN != "" || creds.AssumeRoleExternalID != "" || creds.AssumeRolePolicy != "" || creds.AssumeRoleSessionName != "" {
		config.AssumeRole = append(config.AssumeRole, awsbase.AssumeRole{
			RoleARN:     creds.AssumeRoleARN,
			ExternalID:  creds.AssumeRoleExternalID,
			Policy:      creds.AssumeRolePolicy,
			SessionName: creds.AssumeRoleSessionName,
		})
	}
	if expandedCredsFile != "" {
		config.SharedCredentialsFiles = []string{expandedCredsFile}
	}

	return config, nil
}

// Merge returns a merged credentials
func (c Credentials) Merge(other Credentials) Credentials {
	if other.AccessKey != "" {
		c.AccessKey = other.AccessKey
	}
	if other.SecretKey != "" {
		c.SecretKey = other.SecretKey
	}
	if other.Profile != "" {
		c.Profile = other.Profile
	}
	if other.CredsFile != "" {
		c.CredsFile = other.CredsFile
	}
	if other.Region != "" {
		c.Region = other.Region
	}
	if other.AssumeRoleARN != "" {
		c.AssumeRoleARN = other.AssumeRoleARN
	}
	if other.AssumeRoleSessionName != "" {
		c.AssumeRoleSessionName = other.AssumeRoleSessionName
	}
	if other.AssumeRoleExternalID != "" {
		c.AssumeRoleExternalID = other.AssumeRoleExternalID
	}
	if other.AssumeRolePolicy != "" {
		c.AssumeRolePolicy = other.AssumeRolePolicy
	}
	return c
}
