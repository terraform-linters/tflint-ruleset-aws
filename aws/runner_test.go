package aws

import (
	"sync/atomic"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_NewRunner_sharesClientsAcrossRunners(t *testing.T) {
	// Three providers, two distinct credentials (default + alias "west" share
	// us-east-1, alias "other" uses us-west-2).
	content := `
provider "aws" {
  region     = "us-east-1"
  access_key = "AKID"
  secret_key = "SECRET"
}

provider "aws" {
  alias      = "west"
  region     = "us-east-1"
  access_key = "AKID"
  secret_key = "SECRET"
}

provider "aws" {
  alias      = "other"
  region     = "us-west-2"
  access_key = "AKID"
  secret_key = "SECRET"
}
`

	cases := []struct {
		name      string
		config    *Config
		wantCalls int64
	}{
		{
			name:      "deep check shares clients by credential across runners",
			config:    &Config{DeepCheck: true},
			wantCalls: 2,
		},
		{
			name:      "deep check disabled never builds clients",
			config:    &Config{DeepCheck: false},
			wantCalls: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var count int64
			cache := newClientCache(func(Credentials) (Client, error) {
				atomic.AddInt64(&count, 1)
				return &AwsClient{}, nil
			})

			for range 2 {
				runner := helper.TestRunner(t, map[string]string{"providers.tf": content})
				if _, err := NewRunner(runner, tc.config, cache); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
			}

			if got := atomic.LoadInt64(&count); got != tc.wantCalls {
				t.Fatalf("expected factory to be called %d times, got %d", tc.wantCalls, got)
			}
		})
	}
}
