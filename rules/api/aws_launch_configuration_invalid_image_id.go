package api

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/aws"
)

// AwsLaunchConfigurationInvalidImageIDRule checks whether "aws_instance" has invalid AMI ID
type AwsLaunchConfigurationInvalidImageIDRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	amiIDs        map[string]bool
}

// NewAwsLaunchConfigurationInvalidImageIDRule returns new rule with default attributes
func NewAwsLaunchConfigurationInvalidImageIDRule() *AwsLaunchConfigurationInvalidImageIDRule {
	return &AwsLaunchConfigurationInvalidImageIDRule{
		resourceType:  "aws_launch_configuration",
		attributeName: "image_id",
		amiIDs:        map[string]bool{},
	}
}

// Name returns the rule name
func (r *AwsLaunchConfigurationInvalidImageIDRule) Name() string {
	return "aws_launch_configuration_invalid_image_id"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLaunchConfigurationInvalidImageIDRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLaunchConfigurationInvalidImageIDRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLaunchConfigurationInvalidImageIDRule) Link() string {
	return ""
}

// Metadata returns the metadata about deep checking
func (r *AwsLaunchConfigurationInvalidImageIDRule) Metadata() interface{} {
	return map[string]bool{"deep": true}
}

// Check checks whether "aws_instance" has invalid AMI ID
func (r *AwsLaunchConfigurationInvalidImageIDRule) Check(rr tflint.Runner) error {
	runner := rr.(*aws.Runner)

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		awsClient, err := runner.AwsClient(resource.Body.Attributes)
		if err != nil {
			return err
		}

		var ami string
		err = runner.EvaluateExpr(attribute.Expr, &ami, nil)

		err = runner.EnsureNoError(err, func() error {
			if !r.amiIDs[ami] {
				logger.Debug("Fetch AMI images: %s", ami)
				resp, err := awsClient.EC2.DescribeImages(&ec2.DescribeImagesInput{
					ImageIds: awssdk.StringSlice([]string{ami}),
				})
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok {
						switch aerr.Code() {
						case "InvalidAMIID.Malformed":
							fallthrough
						case "InvalidAMIID.NotFound":
							fallthrough
						case "InvalidAMIID.Unavailable":
							runner.EmitIssue(
								r,
								fmt.Sprintf("\"%s\" is invalid image ID.", ami),
								attribute.Expr.Range(),
							)
							return nil
						}
					}
					err := fmt.Errorf("An error occurred while describing images; %w", err)
					logger.Error("%s", err)
					return err
				}

				if len(resp.Images) != 0 {
					for _, image := range resp.Images {
						r.amiIDs[*image.ImageId] = true
					}
				} else {
					runner.EmitIssue(
						r,
						fmt.Sprintf("\"%s\" is invalid image ID.", ami),
						attribute.Expr.Range(),
					)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
