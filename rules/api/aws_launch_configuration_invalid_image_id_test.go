package api

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLaunchConfigurationInvalidImageID_invalid(t *testing.T) {
	content := `
resource "aws_launch_configuration" "invalid" {
	image_id = "ami-1234abcd"
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

	rule := NewAwsLaunchConfigurationInvalidImageIDRule()
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	expected := helper.Issues{
		{
			Rule:    NewAwsLaunchConfigurationInvalidImageIDRule(),
			Message: "\"ami-1234abcd\" is invalid image ID.",
			Range: hcl.Range{
				Filename: "instances.tf",
				Start:    hcl.Pos{Line: 3, Column: 13},
				End:      hcl.Pos{Line: 3, Column: 27},
			},
		},
	}

	helper.AssertIssues(t, expected, runner.Runner.(*helper.Runner).Issues)
}

func Test_AwsLaunchConfigurationInvalidImageID_valid(t *testing.T) {
	content := `
resource "aws_launch_configuration" "valid" {
	image_id = "ami-9ad76sd1"
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

	rule := NewAwsLaunchConfigurationInvalidImageIDRule()
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	expected := helper.Issues{}
	helper.AssertIssues(t, expected, runner.Runner.(*helper.Runner).Issues)
}

func Test_AwsLaunchConfigurationInvalidImageID_error(t *testing.T) {
	cases := []struct {
		Name    string
		Content string
		Client  *mockClient
		Error   error
	}{
		{
			Name: "AWS API error",
			Content: `
resource "aws_launch_configuration" "valid" {
  image_id = "ami-9ad76sd1"
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
resource "aws_launch_configuration" "valid" {
	image_id = "ami-9ad76sd1"
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

	rule := NewAwsLaunchConfigurationInvalidImageIDRule()

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

func Test_AwsLaunchConfigurationInvalidImageID_AMIError(t *testing.T) {
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
resource "aws_launch_configuration" "not_found" {
	image_id = "ami-9ad76sd1"
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
					Rule:    NewAwsLaunchConfigurationInvalidImageIDRule(),
					Message: "\"ami-9ad76sd1\" is invalid image ID.",
					Range: hcl.Range{
						Filename: "instances.tf",
						Start:    hcl.Pos{Line: 3, Column: 13},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
			Error: false,
		},
		{
			Name: "malformed",
			Content: `
resource "aws_launch_configuration" "malformed" {
	image_id = "image-9ad76sd1"
}`,
			Client: &mockClient{
				describeImages: func(in *ec2.DescribeImagesInput) (map[string]bool, error) {
					if !(len(in.ImageIds) == 1 && in.ImageIds[0] == "image-9ad76sd1") {
						t.Fatalf("passed ImageIds should be [image-9ad76sd1], but got %#v", in.ImageIds)
					}

					return map[string]bool{}, &smithy.GenericAPIError{
						Code:    "InvalidAMIID.Malformed",
						Message: `Invalid id: "image-9ad76sd1" (expecting "ami-...")`,
					}
				},
			},
			Issues: helper.Issues{
				{
					Rule:    NewAwsLaunchConfigurationInvalidImageIDRule(),
					Message: "\"image-9ad76sd1\" is invalid image ID.",
					Range: hcl.Range{
						Filename: "instances.tf",
						Start:    hcl.Pos{Line: 3, Column: 13},
						End:      hcl.Pos{Line: 3, Column: 29},
					},
				},
			},
			Error: false,
		},
		{
			Name: "unavailable",
			Content: `
resource "aws_launch_configuration" "unavailable" {
	image_id = "ami-1234567"
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
					Rule:    NewAwsLaunchConfigurationInvalidImageIDRule(),
					Message: "\"ami-1234567\" is invalid image ID.",
					Range: hcl.Range{
						Filename: "instances.tf",
						Start:    hcl.Pos{Line: 3, Column: 13},
						End:      hcl.Pos{Line: 3, Column: 26},
					},
				},
			},
			Error: false,
		},
	}

	rule := NewAwsLaunchConfigurationInvalidImageIDRule()

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
