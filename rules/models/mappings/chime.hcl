import = "aws-sdk-ruby/apis/chime/2018-05-01/api-2.json"

mapping "aws_chime_voice_connector" {
  name = VoiceConnectorName
  aws_region = VoiceConnectorAwsRegion
}

mapping "aws_chime_voice_connector_group" {
  name = VoiceConnectorGroupName
}

mapping "aws_chime_voice_connector_logging" {
  voice_connector_id = NonEmptyString
}

mapping "aws_chime_voice_connector_origination" {
  voice_connector_id = NonEmptyString
}

mapping "aws_chime_voice_connector_streaming" {
  voice_connector_id = NonEmptyString
  data_retention = DataRetentionInHours
  streaming_notification_targets = StreamingNotificationTargetList
}

mapping "aws_chime_voice_connector_termination" {
  voice_connector_id = NonEmptyString
  cidr_allow_list = StringList
  calling_regions = CallingRegionList
  default_phone_number = E164PhoneNumber
  cps_limit = CpsLimit
}

mapping "aws_chime_voice_connector_termination_credentials" {
  voice_connector_id = NonEmptyString
  credentials = CredentialList
}
