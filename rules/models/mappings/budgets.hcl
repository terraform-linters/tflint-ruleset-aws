import = "api-models-aws/models/budgets/service/2016-10-20/budgets-2016-10-20.json"

mapping "aws_budgets_budget" {
  account_id  = AccountId
  name        = any //BudgetName
  budget_type = BudgetType
  time_unit   = TimeUnit
}

test "aws_budgets_budget" "account_id" {
  ok = "123456789012"
  ng = "abcdefghijkl"
}

test "aws_budgets_budget" "budget_type" {
  ok = "USAGE"
  ng = "MONEY"
}

test "aws_budgets_budget" "time_unit" {
  ok = "MONTHLY"
  ng = "HOURLY"
}
