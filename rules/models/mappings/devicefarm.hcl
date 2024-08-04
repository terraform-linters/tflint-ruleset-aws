import = "aws-sdk-ruby/apis/devicefarm/2015-06-23/api-2.json"

mapping "aws_devicefarm_device_pool" {
  name = Name
  project_arn = AmazonResourceName
  rule = Rules
  description = Message
  tags = TagList
}
    
mapping "aws_devicefarm_network_profile" {
  description = Message
  downlink_loss_percent = PercentInteger
  name = Name
  uplink_loss_percent = PercentInteger
  project_arn = AmazonResourceName
  type = NetworkProfileType
  tags = TagList
}

mapping "aws_devicefarm_project" {
  name = Name
}
    
mapping "aws_devicefarm_upload" {
  content_type = ContentType
  name = Name
  project_arn = AmazonResourceName
  type = UploadType
}
