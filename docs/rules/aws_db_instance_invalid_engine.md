# aws_db_instance_invalid_engine

Disallow using invalid engine name.

## Example

```hcl
resource "aws_db_instance" "default" {
  allocated_storage    = 20
  engine               = "mysql57" // invalid engine name!
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "bar"
  db_subnet_group_name = "my_database_subnet_group"
  parameter_group_name = "default.mysql5.6"
}
```

```
$ tflint
1 issue(s) found:

Warning: "mysql57" is invalid engine. (aws_db_invalid_engine)

  on template.tf line 3:
   3:   engine               = "mysql57" // invalid engine name!
 
```

## Why

Apply will fail. (Plan will succeed with the invalid value though)

## How To Fix

Select valid engine name according to the [document](https://docs.aws.amazon.com/cli/latest/reference/rds/describe-db-engine-versions.html#options)
