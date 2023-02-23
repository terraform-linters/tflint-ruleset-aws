import = "aws-sdk-go/models/apis/ds/2015-04-16/api-2.json"

mapping "aws_directory_service_directory" {
  name             = DirectoryName
  password         = ConnectPassword
  size             = DirectorySize
  vpc_settings     = DirectoryVpcSettings
  connect_settings = DirectoryConnectSettings
  // alias         = AliasName
  description      = Description
  short_name       = DirectoryShortName
  type             = DirectoryType
  edition          = DirectoryEdition
  tags             = Tags
}

mapping "aws_directory_service_conditional_forwarder" {
  directory_id       = DirectoryId
  dns_ips            = DnsIpAddrs
  remote_domain_name = RemoteDomainName
}

mapping "aws_directory_service_log_subscription" {
  directory_id   = DirectoryId
  log_group_name = LogGroupName
}

test "aws_directory_service_directory" "name" {
  valid   = ["corp.notexample.com"]
  invalid = ["@example.com"]
}

test "aws_directory_service_directory" "size" {
  valid   = ["Small"]
  invalid = ["Micro"]
}

test "aws_directory_service_directory" "short_name" {
  valid   = ["CORP"]
  invalid = ["CORP:EXAMPLE"]
}

test "aws_directory_service_directory" "description" {
  valid   = ["example"]
  invalid = ["@example"]
}

test "aws_directory_service_directory" "type" {
  valid   = ["SimpleAD"]
  invalid = ["ActiveDirectory"]
}

test "aws_directory_service_directory" "edition" {
  valid   = ["Enterprise"]
  invalid = ["Free"]
}

test "aws_directory_service_conditional_forwarder" "directory_id" {
  valid   = ["d-1234567890"]
  invalid = ["1234567890"]
}

test "aws_directory_service_conditional_forwarder" "remote_domain_name" {
  valid   = ["example.com"]
  invalid = ["example^com"]
}
