package api

import (
	"fmt"
	"log"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/aws"
)

// AwsLaunchConfigurationInvalidImageIDRule checks whether "aws_instance" has invalid AMI ID
type AwsLaunchConfigurationInvalidImageIDRule struct {
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
func (r *AwsLaunchConfigurationInvalidImageIDRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLaunchConfigurationInvalidImageIDRule) Link() string {
	return ""
}

// Check checks whether "aws_instance" has invalid AMI ID
func (r *AwsLaunchConfigurationInvalidImageIDRule) Check(rr tflint.Runner) error {
	runner := rr.(*aws.Runner)

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var ami string
		err := runner.EvaluateExpr(attribute.Expr, &ami, nil)

		return runner.EnsureNoError(err, func() error {
			if !r.amiIDs[ami] {
				log.Printf("[DEBUG] Fetch AMI images: %s", ami)
				resp, err := runner.AwsClient.EC2.DescribeImages(&ec2.DescribeImagesInput{
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
							runner.EmitIssueOnExpr(
								r,
								fmt.Sprintf("\"%s\" is invalid image ID.", ami),
								attribute.Expr,
							)
							return nil
						}
					}
					err := &tflint.Error{
						Code:    tflint.ExternalAPIError,
						Level:   tflint.ErrorLevel,
						Message: "An error occurred while describing images",
						Cause:   err,
					}
					log.Printf("[ERROR] %s", err)
					return err
				}

				if len(resp.Images) != 0 {
					for _, image := range resp.Images {
						r.amiIDs[*image.ImageId] = true
					}
				} else {
					runner.EmitIssueOnExpr(
						r,
						fmt.Sprintf("\"%s\" is invalid image ID.", ami),
						attribute.Expr,
					)
				}
			}
			return nil
		})
	})
}
