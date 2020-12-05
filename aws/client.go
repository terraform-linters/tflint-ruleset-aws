package aws

import (
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	awsbase "github.com/hashicorp/aws-sdk-go-base"
	"github.com/mitchellh/go-homedir"
)

//go:generate go run github.com/golang/mock/mockgen -destination mock/ec2.go -package mock github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API
//go:generate go run github.com/golang/mock/mockgen -destination mock/elasticache.go -package mock github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface ElastiCacheAPI
//go:generate go run github.com/golang/mock/mockgen -destination mock/elb.go -package mock github.com/aws/aws-sdk-go/service/elb/elbiface ELBAPI
//go:generate go run github.com/golang/mock/mockgen -destination mock/elbv2.go -package mock github.com/aws/aws-sdk-go/service/elbv2/elbv2iface ELBV2API
//go:generate go run github.com/golang/mock/mockgen -destination mock/iam.go -package mock github.com/aws/aws-sdk-go/service/iam/iamiface IAMAPI
//go:generate go run github.com/golang/mock/mockgen -destination mock/rds.go -package mock github.com/aws/aws-sdk-go/service/rds/rdsiface RDSAPI
//go:generate go run github.com/golang/mock/mockgen -destination mock/ecs.go -package mock github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI

// Client is a wrapper of the AWS SDK client
// It has interfaces for each services to make testing easier
type Client struct {
	IAM         iamiface.IAMAPI
	EC2         ec2iface.EC2API
	RDS         rdsiface.RDSAPI
	ElastiCache elasticacheiface.ElastiCacheAPI
	ELB         elbiface.ELBAPI
	ELBV2       elbv2iface.ELBV2API
	ECS         ecsiface.ECSAPI
}

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
func NewClient(creds Credentials) (*Client, error) {
	log.Print("[INFO] Initialize AWS Client")

	config, err := getBaseConfig(creds)
	if err != nil {
		return nil, err
	}

	s, err := awsbase.GetSession(config)
	if err != nil {
		return nil, formatBaseConfigError(err)
	}

	return &Client{
		IAM:         iam.New(s),
		EC2:         ec2.New(s),
		RDS:         rds.New(s),
		ElastiCache: elasticache.New(s),
		ELB:         elb.New(s),
		ELBV2:       elbv2.New(s),
		ECS:         ecs.New(s),
	}, nil
}

func getBaseConfig(creds Credentials) (*awsbase.Config, error) {
	expandedCredsFile, err := homedir.Expand(creds.CredsFile)
	if err != nil {
		return nil, err
	}

	return &awsbase.Config{
		AccessKey:             creds.AccessKey,
		AssumeRoleARN:         creds.AssumeRoleARN,
		AssumeRoleExternalID:  creds.AssumeRoleExternalID,
		AssumeRolePolicy:      creds.AssumeRolePolicy,
		AssumeRoleSessionName: creds.AssumeRoleSessionName,
		SecretKey:             creds.SecretKey,
		Profile:               creds.Profile,
		CredsFilename:         expandedCredsFile,
		Region:                creds.Region,
	}, nil
}

// @see https://github.com/hashicorp/aws-sdk-go-base/blob/v0.3.0/session.go#L87
func formatBaseConfigError(err error) error {
	if strings.Contains(err.Error(), "No valid credential sources found for AWS Provider") {
		return errors.New("No valid credential sources found")
	}
	return err
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
