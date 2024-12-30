import = "aws-sdk-ruby/apis/servicediscovery/2017-03-14/api-2.json"

mapping "aws_service_discovery_http_namespace" {
  name        = any //NamespaceName
  description = ResourceDescription
}

mapping "aws_service_discovery_instance" {
  instance_id = InstanceId
  service_id = ResourceId
  attributes = Attributes
}

mapping "aws_service_discovery_private_dns_namespace" {
  name        = any //NamespaceName
  vpc         = ResourceId
  description = ResourceDescription
}

mapping "aws_service_discovery_public_dns_namespace" {
  name        = any //NamespaceName
  description = ResourceDescription
}

mapping "aws_service_discovery_service" {
  name                       = any // ServiceName
  description                = ResourceDescription
  dns_config                 = DnsConfig
  health_check_config        = HealthCheckConfig
  health_check_custom_config = HealthCheckCustomConfig
}
