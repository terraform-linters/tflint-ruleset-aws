package api

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-ruleset-aws/aws"
)

type mockClient struct {
	aws.AwsClient

	describeSecurityGroups func() (map[string]bool, error)
	describeDBSubnetGroups func() (map[string]bool, error)
	describeImages         func(*ec2.DescribeImagesInput) (map[string]bool, error)
}

func (m *mockClient) DescribeSecurityGroups() (map[string]bool, error) {
	return m.describeSecurityGroups()
}

func (m *mockClient) DescribeDBSubnetGroups() (map[string]bool, error) {
	return m.describeDBSubnetGroups()
}

func (m *mockClient) DescribeImages(in *ec2.DescribeImagesInput) (map[string]bool, error) {
	return m.describeImages(in)
}

func Test_APIListData(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Client   *mockClient
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
			Client: &mockClient{
				describeSecurityGroups: func() (map[string]bool, error) {
					return map[string]bool{
						"sg-12345678": true,
						"sg-abcdefgh": true,
					}, nil
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
			Client: &mockClient{
				describeSecurityGroups: func() (map[string]bool, error) {
					return map[string]bool{
						"sg-1234abcd": true,
						"sg-abcd1234": true,
					}, nil
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
			Client: &mockClient{
				describeSecurityGroups: func() (map[string]bool, error) {
					return map[string]bool{
						"sg-12345678": true,
						"sg-abcdefgh": true,
					}, nil
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

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := NewTestRunner(t, map[string]string{"resource.tf": tc.Content}, tc.Client)

			rule := NewAwsALBInvalidSecurityGroupRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Runner.(*helper.Runner).Issues)
		})
	}
}

func Test_APIStringData(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Client   *mockClient
		Expected helper.Issues
	}{
		{
			Name: "db_subnet_group_name is invalid",
			Content: `
resource "aws_db_instance" "mysql" {
    db_subnet_group_name = "app-server"
}`,
			Client: &mockClient{
				describeDBSubnetGroups: func() (map[string]bool, error) {
					return map[string]bool{
						"app-server1": true,
						"app-server2": true,
					}, nil
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
			Client: &mockClient{
				describeDBSubnetGroups: func() (map[string]bool, error) {
					return map[string]bool{
						"app-server1": true,
						"app-server2": true,
						"app-server":  true,
					}, nil
				},
			},
			Expected: helper.Issues{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := NewTestRunner(t, map[string]string{"resource.tf": tc.Content}, tc.Client)

			rule := NewAwsDBInstanceInvalidDBSubnetGroupRule()
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Runner.(*helper.Runner).Issues)
		})
	}
}

func Test_API_error(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Response error
		Client   *mockClient
		Error    error
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
			Client: &mockClient{
				describeSecurityGroups: func() (map[string]bool, error) {
					return nil, errors.New("MissingRegion: could not find region configuration")
				},
			},
			Error: errors.New("An error occurred while invoking DescribeSecurityGroups; MissingRegion: could not find region configuration"),
		},
	}

	rule := NewAwsALBInvalidSecurityGroupRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := NewTestRunner(t, map[string]string{"resource.tf": tc.Content}, tc.Client)

			err := rule.Check(runner)
			if err == nil {
				t.Fatal("an error is expected, but does not happen")
			}
			if err.Error() != tc.Error.Error() {
				t.Fatalf("`%s` is expected, but got `%s`", tc.Error.Error(), err.Error())
			}
		})
	}
}

func NewTestRunner(t *testing.T, files map[string]string, client *mockClient) *aws.Runner {
	return &aws.Runner{
		Runner: helper.TestRunner(t, files),
		AwsClients: map[string]aws.Client{
			"aws": client,
		},
	}
}
