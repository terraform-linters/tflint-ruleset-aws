import = "aws-sdk-go/models/apis/budgets/2016-10-20/api-2.json"

mapping "aws_budgets_budget" {
  account_id  = AccountId
  name        = BudgetName
  budget_type = BudgetType
  time_unit   = TimeUnit
}

test "aws_budgets_budget" "account_id" {
  valid   = ["123456789012"]
  invalid = ["abcdefghijkl"]
}

test "aws_budgets_budget" "name" {
  valid   = ["budget-ec2-monthly"]
  invalid = ["budget:ec2:monthly"]
}

test "aws_budgets_budget" "budget_type" {
  valid   = ["USAGE"]
  invalid = ["MONEY"]
}

test "aws_budgets_budget" "time_unit" {
  valid   = ["MONTHLY"]
  invalid = ["HOURLY"]
}
