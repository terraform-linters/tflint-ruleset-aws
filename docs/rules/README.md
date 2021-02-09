# Rules

This documentation describes a list of rules available by enabling this ruleset.

### Possible Errors

These rules warn of possible errors that can occur at `terraform apply`. Rules marked with `Deep` are only used when enabling [Deep Checking](../deep_checking.md):

|Rule|Deep|
| --- | --- |
|aws_alb_invalid_security_group|✔|
|aws_alb_invalid_subnet|✔|
|aws_db_instance_invalid_db_subnet_group|✔|
|aws_db_instance_invalid_engine||
|aws_db_instance_invalid_option_group|✔|
|aws_db_instance_invalid_parameter_group|✔|
|aws_db_instance_invalid_type||
|aws_db_instance_invalid_vpc_security_group|✔|
|aws_elasticache_cluster_invalid_parameter_group|✔|
|aws_elasticache_cluster_invalid_security_group|✔|
|aws_elasticache_cluster_invalid_subnet_group|✔|
|aws_elasticache_cluster_invalid_type||
|aws_elb_invalid_instance|✔|
|aws_elb_invalid_security_group|✔|
|aws_elb_invalid_subnet|✔|
|[aws_iam_policy_document_gov_friendly](aws_iam_policy_document_gov_friendly.md)||
|[aws_iam_policy_gov_friendly](aws_iam_policy_gov_friendly.md)||
|[aws_iam_role_policy_gov_friendly](aws_iam_role_policy_gov_friendly.md)||
|aws_instance_invalid_ami|✔|
|aws_instance_invalid_iam_profile|✔|
|aws_instance_invalid_key_name|✔|
|aws_instance_invalid_subnet|✔|
|aws_instance_invalid_vpc_security_group|✔|
|aws_launch_configuration_invalid_iam_profile|✔|
|aws_launch_configuration_invalid_image_id|✔|
|aws_route_invalid_egress_only_gateway|✔|
|aws_route_invalid_gateway|✔|
|aws_route_invalid_instance|✔|
|aws_route_invalid_nat_gateway|✔|
|aws_route_invalid_network_interface|✔|
|aws_route_invalid_route_table|✔|
|aws_route_invalid_vpc_peering_connection|✔|
|[aws_s3_bucket_name](aws_s3_bucket_name.md)||
|[aws_route_not_specified_target](aws_route_not_specified_target.md)||
|[aws_route_specified_multiple_targets](aws_route_specified_multiple_targets.md)||

#### SDK-based Validations

700+ rules based on the aws-sdk validations are also available. See [full list](../../rules/awsrules/models/).

### Best Practices

These rules suggest to better ways.

- [aws_instance_previous_type](aws_instance_previous_type.md)
- [aws_db_instance_invalid_engine](aws_db_instance_invalid_engine.md)
- [aws_db_instance_previous_type](aws_db_instance_previous_type.md)
- [aws_db_instance_default_parameter_group](aws_db_instance_default_parameter_group.md)
- [aws_elasticache_cluster_previous_type](aws_elasticache_cluster_previous_type.md)
- [aws_elasticache_cluster_default_parameter_group](aws_elasticache_cluster_default_parameter_group.md)
