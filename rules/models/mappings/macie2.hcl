import = "aws-sdk-go/models/apis/macie2/2020-01-01/api-2.json"

mapping "aws_macie2_account" {
  finding_publishing_frequency = FindingPublishingFrequency
  status = MacieStatus
}

mapping "aws_macie2_classification_job" {
  schedule_frequency = JobScheduleFrequency
  job_type = JobType
  s3_job_definition = S3JobDefinition
  tags = TagMap
  job_status = JobStatus
}

mapping "aws_macie2_custom_data_identifier" {
  tags = TagMap
}

mapping "aws_macie2_findings_filter" {
  finding_criteria = FindingCriteria
  action = FindingsFilterAction
  tags = TagMap
}

mapping "aws_macie2_member" {
  tags = TagMap
}
