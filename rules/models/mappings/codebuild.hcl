import = "api-models-aws/models/codebuild/service/2016-10-06/codebuild-2016-10-06.json"

mapping "aws_codebuild_project" {
  build_timeout = TimeOut
  description   = ProjectDescription
}

mapping "aws_codebuild_report_group" {
  name = ReportGroupName
  type = ReportType
  export_config = ReportExportConfig
  tags = TagList
}

mapping "aws_codebuild_source_credential" {
  auth_type   = AuthType
  server_type = ServerType
  token       = SensitiveNonEmptyString
  user_name   = NonEmptyString
}
