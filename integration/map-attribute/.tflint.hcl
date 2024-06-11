plugin "terraform" {
  enabled = false
}

plugin "aws" {
  enabled = true
}

rule "aws_resource_tags" {
  enabled = true
  tags = ["Environment", "Name", "Type"]
}
