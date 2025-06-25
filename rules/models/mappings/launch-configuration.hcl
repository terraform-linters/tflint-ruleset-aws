import = "api-models-aws/models/ec2/service/2016-11-15/ec2-2016-11-15.json"

// NOTE: `aws_launch_configuration` mapping is already defined in autoscaling.hcl
//       The following mapping is to import ec2 types.
mapping "aws_launch_configuration" {
  instance_type = InstanceType
}
