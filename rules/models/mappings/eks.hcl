import = "aws-sdk-ruby/apis/eks/2017-11-01/api-2.json"

mapping "aws_eks_addon" {
  cluster_name = ClusterName
  tags = TagMap
  service_account_role_arn = RoleArn
}

mapping "aws_eks_cluster" {
  name                      = ClusterName
  role_arn                  = String
  vpc_config                = VpcConfigRequest
  enabled_cluster_log_types = Logging
  version                   = String
}

mapping "aws_eks_fargate_profile" {
  selector = FargateProfileSelectors
  tags = TagMap
}

mapping "aws_eks_identity_provider_config" {
  oidc = OidcIdentityProviderConfigRequest
  tags = TagMap
}

mapping "aws_eks_node_group" {
  scaling_config = NodegroupScalingConfig
  ami_type = AMITypes
  capacity_type = CapacityTypes
  labels = labelsMap
  launch_template = LaunchTemplateSpecification
  remote_access = RemoteAccessConfig
  tags = TagMap
  taint = taintsList
}

test "aws_eks_cluster" "name" {
  ok = "example"
  ng = "@example"
}
