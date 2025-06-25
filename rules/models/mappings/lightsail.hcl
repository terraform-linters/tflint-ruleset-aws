import = "api-models-aws/models/lightsail/service/2016-11-28/lightsail-2016-11-28.json"

mapping "aws_lightsail_domain" {
  domain_name = DomainName
}

mapping "aws_lightsail_instance" {
  name              = string
  availability_zone = string
  blueprint_id      = NonEmptyString
  bundle_id         = NonEmptyString
  key_pair_name     = ResourceName
  user_data         = string
}

mapping "aws_lightsail_instance_public_ports" {
  instance_name = ResourceName
  port_info = PortInfo
}

mapping "aws_lightsail_key_pair" {
  name       = ResourceName
  pgp_key    = any
  public_key = any
}

mapping "aws_lightsail_static_ip" {
  name = ResourceName
}

mapping "aws_lightsail_static_ip_attachment" {
  static_ip_name = ResourceName
  instance_name  = ResourceName
}
