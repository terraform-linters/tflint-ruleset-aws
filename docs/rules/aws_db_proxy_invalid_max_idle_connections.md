# aws_db_proxy_invalid_max_idle_connections

Disallow invalid `max_idle_connections_percent`.

## Example

```hcl
resource "aws_db_proxy_default_target_group" "example" {
  connection_pool_config {
    max_connections_percent      = 10
    max_idle_connections_percent = 50
  }
}
```

```
$ tflint
1 issue(s) found:

Error: max_idle_connections_percent 50 must be less than max_connections_percent 10 (aws_db_proxy_invalid_max_idle_connections)

  on resource.tf line 4:
   4:     max_idle_connections_percent = 50
 
```

## Why

Apply will fail. (Plan will succeed with the invalid value though)

## How To Fix

Set `max_idle_connections_percent` to the value less than `max_connections_percent` according to the [document](https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_ConnectionPoolConfiguration.html).
