# aws_db_instance_invalid_type

Disallow using invalid type.

## Example

```hcl
resource "aws_db_instance" "default" {
  allocated_storage    = 20
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t1.micro" // invalid type!
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

Warning: "db.t1.micro" is previous generation instance type. (aws_db_instance_previous_type)

  on template.tf line 5:
   5:   instance_class       = "db.t1.micro" // invalid type!

```

## Why

Apply will fail. (Plan will succeed with the invalid value though)

## How To Fix

Select valid type according to the [document](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.html#Concepts.DBInstanceClass.Types)

