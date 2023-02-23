import = "aws-sdk-go/models/apis/cur/2017-01-06/api-2.json"

mapping "aws_cur_report_definition" {
  report_name = ReportName
  time_unit   = TimeUnit
  format      = ReportFormat
  compression = CompressionFormat
  s3_bucket   = S3Bucket
  s3_prefix   = S3Prefix
  s3_region   = AWSRegion
}

test "aws_cur_report_definition" "report_name" {
  valid   = ["example-cur-report-definition"]
  invalid = ["example/cur-report-definition"]
}

test "aws_cur_report_definition" "time_unit" {
  valid   = ["HOURLY"]
  invalid = ["FORNIGHTLY"]
}

test "aws_cur_report_definition" "format" {
  valid   = ["textORcsv"]
  invalid = ["textORjson"]
}

test "aws_cur_report_definition" "compression" {
  valid   = ["ZIP"]
  invalid = ["TAR"]
}

test "aws_cur_report_definition" "s3_region" {
  valid   = ["us-east-1"]
  invalid = ["us-gov-east-1"]
}
