package models

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

// This test is a manual test on length
// it is implemented for only one rule to test that the minimal logic works correctly.
func Test_Length(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It is too short",
			Content: `
resource "aws_launch_template" "foo" {
	name = "go"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLaunchTemplateInvalidNameRule(),
					Message: `name must be 3 characters or higher`,
				},
			},
		},
		{
			Name: "It is too long",
			Content: `
resource "aws_launch_template" "foo" {
	name = "Lorem_ipsum_dolor_sit_amet_consectetur_adipisicing_elit_sed_do_eiusmod_tempor_incididunt_ut_labore_et_dolore_magna_aliqua.Ut_enim_ad_minim_veniam_quis_nostrud_exercitation_ullamco_laboris_nisi_ut_aliquip_ex_ea_commodo_consequat.Duis_aute_irure_dolor_in_reprehenderit_in_voluptate_velit_esse_cillum_dolore_eu_fugiat_nulla_pariatur.Excepteur_sint_occaecat_cupidatat_non_proident_sunt_in_culpa_qui_officia_deserunt_mollit_anim_id_est_laborum."
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLaunchTemplateInvalidNameRule(),
					Message: `name must be 128 characters or less`,
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_launch_template" "foo" {
	name = "foo"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLaunchTemplateInvalidNameRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
