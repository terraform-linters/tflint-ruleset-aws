resource "aws_instance" "invalid" {
  instance_type = "x2.2xlarge"
}

resource "aws_instance" "previous" {
  instance_type = "t1.micro"
}
