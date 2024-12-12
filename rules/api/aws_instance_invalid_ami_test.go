package api

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsInstanceInvalidAMI_invalid(t *testing.T) {
	content := `
resource "aws_instance" "invalid" {
	ami = "ami-1234abcd"
}`
	client := &mockClient{
		describeImages: func(in *ec2.DescribeImagesInput) (map[string]bool, error) {
			if !(len(in.ImageIds) == 1 && in.ImageIds[0] == "ami-1234abcd") {
				t.Fatal("passed ImageIds should be ami-1234abcd")
			}
			return map[string]bool{}, nil
		},
	}
	runner := NewTestRunner(t, map[string]string{"instances.tf": content}, client)

	rule := NewAwsInstanceInvalidAMIRule()
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	expected := helper.Issues{
		{
			Rule:    NewAwsInstanceInvalidAMIRule(),
			Message: "\"ami-1234abcd\" is invalid AMI ID.",
			Range: hcl.Range{
				Filename: "instances.tf",
				Start:    hcl.Pos{Line: 3, Column: 8},
				End:      hcl.Pos{Line: 3, Column: 22},
			},
		},
	}
	helper.AssertIssues(t, expected, runner.Runner.(*helper.Runner).Issues)
}

func Test_AwsInstanceInvalidAMI_valid(t *testing.T) {
	content := `
resource "aws_instance" "valid" {
	ami = "ami-9ad76sd1"
}`
	client := &mockClient{
		describeImages: func(in *ec2.DescribeImagesInput) (map[string]bool, error) {
			if !(len(in.ImageIds) == 1 && in.ImageIds[0] == "ami-9ad76sd1") {
				t.Fatal("passed ImageIds should be ami-9ad76sd1")
			}
			return map[string]bool{"ami-9ad76sd1": true}, nil
		},
	}
	runner := NewTestRunner(t, map[string]string{"instances.tf": content}, client)

	rule := NewAwsInstanceInvalidAMIRule()
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	expected := helper.Issues{}
	helper.AssertIssues(t, expected, runner.Runner.(*helper.Runner).Issues)
}

func Test_AwsInstanceInvalidAMI_error(t *testing.T) {
	cases := []struct {
		Name    string
		Content string
		Client  *mockClient
		Error   error
	}{
		{
			Name: "AWS API error",
			Content: `
resource "aws_instance" "valid" {
  ami = "ami-9ad76sd1"
}`,
			Client: &mockClient{
				describeImages: func(in *ec2.DescribeImagesInput) (map[string]bool, error) {
					if !(len(in.ImageIds) == 1 && in.ImageIds[0] == "ami-9ad76sd1") {
						t.Fatal("passed ImageIds should be ami-9ad76sd1")
					}

					return map[string]bool{}, &smithy.GenericAPIError{
						Code:    "MissingRegion",
						Message: "could not find region configuration",
					}
				},
			},
			Error: errors.New("An error occurred while describing images; api error MissingRegion: could not find region configuration"),
		},
		{
			Name: "Unexpected error",
			Content: `
resource "aws_instance" "valid" {
  ami = "ami-9ad76sd1"
}`,
			Client: &mockClient{
				describeImages: func(in *ec2.DescribeImagesInput) (map[string]bool, error) {
					if !(len(in.ImageIds) == 1 && in.ImageIds[0] == "ami-9ad76sd1") {
						t.Fatal("passed ImageIds should be ami-9ad76sd1")
					}

					return map[string]bool{}, errors.New("Unexpected")
				},
			},
			Error: errors.New("An error occurred while describing images; Unexpected"),
		},
	}

	rule := NewAwsInstanceInvalidAMIRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := NewTestRunner(t, map[string]string{"instances.tf": tc.Content}, tc.Client)

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

func Test_AwsInstanceInvalidAMI_AMIError(t *testing.T) {
	cases := []struct {
		Name    string
		Content string
		Client  *mockClient
		Issues  helper.Issues
		Error   bool
	}{
		{
			Name: "not found",
			Content: `
resource "aws_instance" "not_found" {
	ami = "ami-9ad76sd1"
}`,
			Client: &mockClient{
				describeImages: func(in *ec2.DescribeImagesInput) (map[string]bool, error) {
					if !(len(in.ImageIds) == 1 && in.ImageIds[0] == "ami-9ad76sd1") {
						t.Fatal("passed ImageIds should be ami-9ad76sd1")
					}

					return map[string]bool{}, &smithy.GenericAPIError{
						Code:    "InvalidAMIID.NotFound",
						Message: "The image id '[ami-9ad76sd1]' does not exist",
					}
				},
			},
			Issues: helper.Issues{
				{
					Rule:    NewAwsInstanceInvalidAMIRule(),
					Message: "\"ami-9ad76sd1\" is invalid AMI ID.",
					Range: hcl.Range{
						Filename: "instances.tf",
						Start:    hcl.Pos{Line: 3, Column: 8},
						End:      hcl.Pos{Line: 3, Column: 22},
					},
				},
			},
			Error: false,
		},
		{
			Name: "malformed",
			Content: `
resource "aws_instance" "malformed" {
	ami = "image-9ad76sd1"
}`,
			Client: &mockClient{
				describeImages: func(in *ec2.DescribeImagesInput) (map[string]bool, error) {
					if !(len(in.ImageIds) == 1 && in.ImageIds[0] == "image-9ad76sd1") {
						t.Fatalf("passed ImageIds should be [image-9ad76sd1], but got %#v", in.ImageIds)
					}

					return map[string]bool{}, &smithy.GenericAPIError{
						Code:    "InvalidAMIID.Malformed",
						Message: `Invalid id: "image-9ad76sd1"`,
					}
				},
			},
			Issues: helper.Issues{
				{
					Rule:    NewAwsInstanceInvalidAMIRule(),
					Message: "\"image-9ad76sd1\" is invalid AMI ID.",
					Range: hcl.Range{
						Filename: "instances.tf",
						Start:    hcl.Pos{Line: 3, Column: 8},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
			},
			Error: false,
		},
		{
			Name: "unavailable",
			Content: `
resource "aws_instance" "unavailable" {
	ami = "ami-1234567"
}`,
			Client: &mockClient{
				describeImages: func(in *ec2.DescribeImagesInput) (map[string]bool, error) {
					if !(len(in.ImageIds) == 1 && in.ImageIds[0] == "ami-1234567") {
						t.Fatal("passed ImageIds should be ami-1234567")
					}

					return map[string]bool{}, &smithy.GenericAPIError{
						Code:    "InvalidAMIID.Unavailable",
						Message: "The image ID: 'ami-1234567' is no longer available",
					}
				},
			},
			Issues: helper.Issues{
				{
					Rule:    NewAwsInstanceInvalidAMIRule(),
					Message: "\"ami-1234567\" is invalid AMI ID.",
					Range: hcl.Range{
						Filename: "instances.tf",
						Start:    hcl.Pos{Line: 3, Column: 8},
						End:      hcl.Pos{Line: 3, Column: 21},
					},
				},
			},
			Error: false,
		},
	}

	rule := NewAwsInstanceInvalidAMIRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := NewTestRunner(t, map[string]string{"instances.tf": tc.Content}, tc.Client)

			err := rule.Check(runner)
			if err != nil && !tc.Error {
				t.Fatalf("Failed `%s` test: unexpected error occurred: %s", tc.Name, err)
			}
			if err == nil && tc.Error {
				t.Fatalf("Failed `%s` test: expected to return an error, but nothing occurred", tc.Name)
			}

			helper.AssertIssues(t, tc.Issues, runner.Runner.(*helper.Runner).Issues)
		})
	}
}
