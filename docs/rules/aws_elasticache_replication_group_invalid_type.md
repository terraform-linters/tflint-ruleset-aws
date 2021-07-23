# aws_elasticache_replication_group_invalid_type

Disallow using invalid type.

## Example

```hcl
resource "aws_elasticache_replication_group" "default" {
  node_type                     = "cache.t3.mini" // invalid type!
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

Warning: "cache.t3.mini" is an invalid node type. (aws_elasticache_replication_group_invalid_type)

  on template.tf line 5:
   2:   node_type       = "cache.t3.mini" // invalid type!

```

## Why

Apply will fail. (Plan will succeed with the invalid value though)

## How To Fix

Select valid type according to the [document](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/CacheNodes.SupportedTypes.html#CacheNodes.SupportedTypesByRegion)
