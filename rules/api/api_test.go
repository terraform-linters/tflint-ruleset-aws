package api

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/golang/mock/gomock"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	awsruleset "github.com/terraform-linters/tflint-ruleset-aws/aws"
	"github.com/terraform-linters/tflint-ruleset-aws/aws/mock"
)

func Test_APIListData(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Response []*ec2.SecurityGroup
		Expected helper.Issues
	}{
		{
			Name: "security group is invalid",
			Content: `
resource "aws_alb" "balancer" {
    security_groups = [
        "sg-1234abcd",
        "sg-abcd1234",
    ]
}`,
			Response: []*ec2.SecurityGroup{
				{
					GroupId: aws.String("sg-12345678"),
				},
				{
					GroupId: aws.String("sg-abcdefgh"),
				},
			},
			Expected: helper.Issues{
				{
					Rule:    NewAwsALBInvalidSecurityGroupRule(),
					Message: "\"sg-1234abcd\" is invalid security group.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 9},
						End:      hcl.Pos{Line: 4, Column: 22},
					},
				},
				{
					Rule:    NewAwsALBInvalidSecurityGroupRule(),
					Message: "\"sg-abcd1234\" is invalid security group.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 9},
						End:      hcl.Pos{Line: 5, Column: 22},
					},
				},
			},
		},
		{
			Name: "security group is valid",
			Content: `
resource "aws_alb" "balancer" {
    security_groups = [
        "sg-1234abcd",
        "sg-abcd1234",
    ]
}`,
			Response: []*ec2.SecurityGroup{
				{
					GroupId: aws.String("sg-1234abcd"),
				},
				{
					GroupId: aws.String("sg-abcd1234"),
				},
			},
			Expected: helper.Issues{},
		},
		{
			Name: "use list variables",
			Content: `
variable "security_groups" {
    default = ["sg-1234abcd", "sg-abcd1234"]
}

resource "aws_alb" "balancer" {
    security_groups = "${var.security_groups}"
}`,
			Response: []*ec2.SecurityGroup{
				{
					GroupId: aws.String("sg-12345678"),
				},
				{
					GroupId: aws.String("sg-abcdefgh"),
				},
			},
			Expected: helper.Issues{
				{
					Rule:    NewAwsALBInvalidSecurityGroupRule(),
					Message: "\"sg-1234abcd\" is invalid security group.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 23},
						End:      hcl.Pos{Line: 7, Column: 47},
					},
				},
				{
					Rule:    NewAwsALBInvalidSecurityGroupRule(),
					Message: "\"sg-abcd1234\" is invalid security group.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 23},
						End:      hcl.Pos{Line: 7, Column: 47},
					},
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range cases {
		runner := NewTestRunner(t, map[string]string{"resource.tf": tc.Content})

		ec2mock := mock.NewMockEC2API(ctrl)
		ec2mock.EXPECT().DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{}).Return(&ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: tc.Response,
		}, nil)
		runner.AwsClient.EC2 = ec2mock

		rule := NewAwsALBInvalidSecurityGroupRule()
		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Runner.(*helper.Runner).Issues)
	}
}

func Test_APIStringData(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Response []*rds.DBSubnetGroup
		Expected helper.Issues
	}{
		{
			Name: "db_subnet_group_name is invalid",
			Content: `
resource "aws_db_instance" "mysql" {
    db_subnet_group_name = "app-server"
}`,
			Response: []*rds.DBSubnetGroup{
				{
					DBSubnetGroupName: aws.String("app-server1"),
				},
				{
					DBSubnetGroupName: aws.String("app-server2"),
				},
			},
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceInvalidDBSubnetGroupRule(),
					Message: "\"app-server\" is invalid DB subnet group name.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 28},
						End:      hcl.Pos{Line: 3, Column: 40},
					},
				},
			},
		},
		{
			Name: "db_subnet_group_name is valid",
			Content: `
resource "aws_db_instance" "mysql" {
    db_subnet_group_name = "app-server"
}`,
			Response: []*rds.DBSubnetGroup{
				{
					DBSubnetGroupName: aws.String("app-server1"),
				},
				{
					DBSubnetGroupName: aws.String("app-server2"),
				},
				{
					DBSubnetGroupName: aws.String("app-server"),
				},
			},
			Expected: helper.Issues{},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range cases {
		runner := NewTestRunner(t, map[string]string{"resource.tf": tc.Content})

		rdsmock := mock.NewMockRDSAPI(ctrl)
		rdsmock.EXPECT().DescribeDBSubnetGroups(&rds.DescribeDBSubnetGroupsInput{}).Return(&rds.DescribeDBSubnetGroupsOutput{
			DBSubnetGroups: tc.Response,
		}, nil)
		runner.AwsClient.RDS = rdsmock

		rule := NewAwsDBInstanceInvalidDBSubnetGroupRule()
		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Runner.(*helper.Runner).Issues)
	}
}

func Test_API_error(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Response error
		Error    tflint.Error
	}{
		{
			Name: "API error",
			Content: `
resource "aws_alb" "balancer" {
    security_groups = [
        "sg-1234abcd",
        "sg-abcd1234",
    ]
}`,
			Response: errors.New("MissingRegion: could not find region configuration"),
			Error: tflint.Error{
				Code:    tflint.ExternalAPIError,
				Level:   tflint.ErrorLevel,
				Message: "An error occurred while invoking DescribeSecurityGroups; MissingRegion: could not find region configuration",
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rule := NewAwsALBInvalidSecurityGroupRule()

	for _, tc := range cases {
		runner := NewTestRunner(t, map[string]string{"resource.tf": tc.Content})

		ec2mock := mock.NewMockEC2API(ctrl)
		ec2mock.EXPECT().DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{}).Return(nil, tc.Response)
		runner.AwsClient.EC2 = ec2mock

		err := rule.Check(runner)
		AssertAppError(t, tc.Error, err)
	}
}

func NewTestRunner(t *testing.T, files map[string]string) *awsruleset.Runner {
	return &awsruleset.Runner{
		Runner:    helper.TestRunner(t, files),
		AwsClient: &awsruleset.Client{},
	}
}

// AssertAppError is an assertion helper for comparing tflint.Error
func AssertAppError(t *testing.T, expected tflint.Error, got error) {
	if appErr, ok := got.(*tflint.Error); ok {
		if appErr == nil {
			t.Fatalf("expected err is `%s`, but nothing occurred", expected.Error())
		}
		if appErr.Code != expected.Code {
			t.Fatalf("expected error code is `%s`, but get `%s`", expected.Code, appErr.Code)
		}
		if appErr.Level != expected.Level {
			t.Fatalf("expected error level is `%s`, but get `%s`", expected.Level, appErr.Level)
		}
		if appErr.Error() != expected.Error() {
			t.Fatalf("expected error is `%s`, but get `%s`", expected.Error(), appErr.Error())
		}
	} else {
		t.Fatalf("unexpected error occurred: %s", got)
	}
}
