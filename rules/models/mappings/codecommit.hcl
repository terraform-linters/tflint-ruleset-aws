import = "api-models-aws/models/codecommit/service/2015-04-13/codecommit-2015-04-13.json"

mapping "aws_codecommit_approval_rule_template" {
  content = ApprovalRuleTemplateContent
  name = ApprovalRuleTemplateName
  description = ApprovalRuleTemplateDescription
}

mapping "aws_codecommit_approval_rule_template_association" {
  approval_rule_template_name = ApprovalRuleTemplateName
  repository_name = RepositoryName
}

mapping "aws_codecommit_repository" {
  repository_name = RepositoryName
  description     = RepositoryDescription
  default_branch  = BranchName
}

mapping "aws_codecommit_trigger" {
  repository_name = RepositoryName
}

test "aws_codecommit_repository" "repository_name" {
  ok = "MyTestRepository"
  ng = "mytest@repository"
}
