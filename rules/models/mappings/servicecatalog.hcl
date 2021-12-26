import = "aws-sdk-go/models/apis/servicecatalog/2015-12-10/api-2.json"

mapping "aws_servicecatalog_budget_resource_association" {
  budget_name = BudgetName
  resource_id = Id
}

mapping "aws_servicecatalog_constraint" {
  parameters = ConstraintParameters
  portfolio_id = Id
  product_id = Id
  type = ConstraintType
  accept_language = AcceptLanguage
  description = ConstraintDescription
}

mapping "aws_servicecatalog_portfolio" {
  name          = PortfolioDisplayName
  description   = PortfolioDescription
  provider_name = ProviderName
  tags          = AddTags
}

mapping "aws_servicecatalog_portfolio_share" {
  portfolio_id = Id
  principal_id = Id
  type = PortfolioShareType
  accept_language = AcceptLanguage
}

mapping "aws_servicecatalog_principal_portfolio_association" {
  portfolio_id = Id
  principal_arn = PrincipalARN
  accept_language = AcceptLanguage
  principal_type = PrincipalType
}

mapping "aws_servicecatalog_product" {
  name = ProductViewName
  owner = ProductViewOwner
  provisioning_artifact_parameters = ProvisioningArtifactProperties
  type = ProductType
  accept_language = AcceptLanguage
  description = ProductViewShortDescription
  distributor = ProductViewOwner
  support_description = SupportDescription
  support_email = SupportEmail
  support_url = SupportUrl
  tags = AddTags
}

mapping "aws_servicecatalog_product_portfolio_association" {
  portfolio_id = Id
  product_id = Id
  accept_language = AcceptLanguage
  source_portfolio_id = Id
}

mapping "aws_servicecatalog_provisioned_product" {
  name = ProvisionedProductName
  accept_language = AcceptLanguage
  notification_arns = NotificationArns
  path_id = Id
  path_name = PortfolioDisplayName
  product_id = Id
  product_name = ProductViewName
  provisioning_artifact_id = Id
  provisioning_artifact_name = ProvisioningArtifactName
  provisioning_parameters = ProvisioningParameters
}

mapping "aws_servicecatalog_provisioning_artifact" {
  product_id = Id
  accept_language = AcceptLanguage
  description = ProvisioningArtifactDescription
  disable_template_validation = DisableTemplateValidation
  guidance = ProvisioningArtifactGuidance
  name = ProvisioningArtifactName
  type = ProvisioningArtifactType
}

mapping "aws_servicecatalog_service_action" {
  definition = ServiceActionDefinitionMap
  name = ServiceActionName
  accept_language = AcceptLanguage
  description = ServiceActionDescription
}

mapping "aws_servicecatalog_tag_option" {
  key = TagOptionKey
  value = TagOptionValue
}

mapping "aws_servicecatalog_tag_option_resource_association" {
  resource_id = ResourceId
  tag_option_id = TagOptionId
}
