# aws_elasticache_replication_group_previous_type

Disallow using previous node types.

## Example

```hcl
resource "aws_elasticache_replication_group" "default" {
  node_type                     = "cache.t1.micro" // previous node type!
  at_rest_encryption_enabled    = true
  automatic_failover_enabled    = true
  engine_version                = "6.x"
  maintenance_window            = "thu:02:30-thu:03:30"
  apply_immediately             = false
  parameter_group_name          = "custom.redis6.x.cluster.on"
  port                          = 6379
  replication_group_description = " "
  replication_group_id          = "replication_group_id"
  snapshot_retention_limit      = 1
  subnet_group_name             = aws_elasticache_subnet_group.private.name
  security_group_ids            = [aws_security_group.redis_service.id]

  cluster_mode {
    replicas_per_node_group = 1
    num_node_groups         = 2
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: "cache.t1.micro" is previous generation node type. (aws_elasticache_replication_group_previous_type)

  on template.tf line 6:
   2:   node_type            = "cache.t1.micro" // previous node type!

```

## Why

Previous node types are inferior to current generation in terms of performance and fee. Unless there is a special reason, you should avoid to use these ones.

## How To Fix

Select a current generation node type according to the [upgrade paths](https://aws.amazon.com/elasticache/previous-generation/).
