package aws

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	awsbase "github.com/hashicorp/aws-sdk-go-base/v2"
	"github.com/mitchellh/go-homedir"
)

func Test_getBaseConfig(t *testing.T) {
	home, err := homedir.Expand("~/")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Name     string
		Creds    Credentials
		Expected *awsbase.Config
	}{
		{
			Name: "static credentials",
			Creds: Credentials{
				AccessKey: "AWS_ACCESS_KEY",
				SecretKey: "AWS_SECRET_KEY",
				Region:    "us-east-1",
			},
			Expected: &awsbase.Config{
				AccessKey:              "AWS_ACCESS_KEY",
				SecretKey:              "AWS_SECRET_KEY",
				Region:                 "us-east-1",
				CallerDocumentationURL: "https://github.com/terraform-linters/tflint-ruleset-aws/blob/master/docs/deep_checking.md",
				CallerName:             "tflint-ruleset-aws",
			},
		},
		{
			Name: "shared credentials",
			Creds: Credentials{
				Profile:   "default",
				CredsFile: "~/.aws/creds",
				Region:    "us-east-1",
			},
			Expected: &awsbase.Config{
				Profile:                "default",
				SharedCredentialsFiles: []string{filepath.Join(home, ".aws", "creds")},
				Region:                 "us-east-1",
				CallerDocumentationURL: "https://github.com/terraform-linters/tflint-ruleset-aws/blob/master/docs/deep_checking.md",
				CallerName:             "tflint-ruleset-aws",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			base, err := getBaseConfig(tc.Creds)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tc.Expected, base) {
				t.Fatalf("diff=%s", cmp.Diff(tc.Expected, base))
			}
		})
	}
}

func Test_Merge(t *testing.T) {
	cases := []struct {
		Name     string
		Self     Credentials
		Other    Credentials
		Expected Credentials
	}{
		{
			Name: "self is empty",
			Self: Credentials{},
			Other: Credentials{
				AccessKey:             "AWS_ACCESS_KEY",
				SecretKey:             "AWS_SECRET_KEY",
				Profile:               "default",
				CredsFile:             "~/.aws/creds",
				AssumeRoleARN:         "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME",
				AssumeRoleSessionName: "SESSION_NAME",
				AssumeRoleExternalID:  "EXTERNAL_ID",
				AssumeRolePolicy:      "POLICY_NAME",
				Region:                "us-east-1",
			},
			Expected: Credentials{
				AccessKey:             "AWS_ACCESS_KEY",
				SecretKey:             "AWS_SECRET_KEY",
				Profile:               "default",
				CredsFile:             "~/.aws/creds",
				AssumeRoleARN:         "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME",
				AssumeRoleSessionName: "SESSION_NAME",
				AssumeRoleExternalID:  "EXTERNAL_ID",
				AssumeRolePolicy:      "POLICY_NAME",
				Region:                "us-east-1",
			},
		},
		{
			Name: "other is empty",
			Self: Credentials{
				AccessKey:             "AWS_ACCESS_KEY_2",
				SecretKey:             "AWS_SECRET_KEY_2",
				Profile:               "staging",
				CredsFile:             "~/.aws/creds_stg",
				AssumeRoleARN:         "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME",
				AssumeRoleSessionName: "SESSION_NAME",
				AssumeRoleExternalID:  "EXTERNAL_ID",
				AssumeRolePolicy:      "POLICY_NAME",
				Region:                "ap-northeast-1",
			},
			Other: Credentials{},
			Expected: Credentials{
				AccessKey:             "AWS_ACCESS_KEY_2",
				SecretKey:             "AWS_SECRET_KEY_2",
				Profile:               "staging",
				CredsFile:             "~/.aws/creds_stg",
				AssumeRoleARN:         "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME",
				AssumeRoleSessionName: "SESSION_NAME",
				AssumeRoleExternalID:  "EXTERNAL_ID",
				AssumeRolePolicy:      "POLICY_NAME",
				Region:                "ap-northeast-1",
			},
		},
		{
			Name: "merged",
			Self: Credentials{
				AccessKey:             "AWS_ACCESS_KEY_2",
				SecretKey:             "AWS_SECRET_KEY_2",
				Profile:               "staging",
				CredsFile:             "~/.aws/creds_stg",
				AssumeRoleARN:         "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME_2",
				AssumeRoleSessionName: "SESSION_NAME_2",
				AssumeRoleExternalID:  "EXTERNAL_ID_2",
				AssumeRolePolicy:      "POLICY_NAME_2",
				Region:                "ap-northeast-1",
			},
			Other: Credentials{
				AccessKey:             "AWS_ACCESS_KEY",
				SecretKey:             "AWS_SECRET_KEY",
				Profile:               "default",
				CredsFile:             "~/.aws/creds",
				AssumeRoleARN:         "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME",
				AssumeRoleSessionName: "SESSION_NAME",
				AssumeRoleExternalID:  "EXTERNAL_ID",
				AssumeRolePolicy:      "POLICY_NAME",
				Region:                "us-east-1",
			},
			Expected: Credentials{
				AccessKey:             "AWS_ACCESS_KEY",
				SecretKey:             "AWS_SECRET_KEY",
				Profile:               "default",
				CredsFile:             "~/.aws/creds",
				AssumeRoleARN:         "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME",
				AssumeRoleSessionName: "SESSION_NAME",
				AssumeRoleExternalID:  "EXTERNAL_ID",
				AssumeRolePolicy:      "POLICY_NAME",
				Region:                "us-east-1",
			},
		},
	}

	for _, tc := range cases {
		ret := tc.Self.Merge(tc.Other)
		if !cmp.Equal(tc.Expected, ret) {
			t.Fatalf("Failed `%s` test: Diff=%s", tc.Name, cmp.Diff(tc.Expected, ret))
		}
	}
}
