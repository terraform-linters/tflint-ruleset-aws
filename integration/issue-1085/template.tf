variable "aws_lb_controller_enabled" {
  type    = bool
  default = false
}

locals {
  namespaces = {
    system = merge({},
      var.aws_lb_controller_enabled ? {
        lb_controller = "aws-lb-controller"
      } : {},
    )
  }

  ingresses = [
    local.namespaces.system.lb_controller
  ]
}

resource "kubernetes_network_policy" "default" {
  for_each = toset(local.ingresses)

  metadata {
    namespace = each.key
    name      = "test-${each.key}"
  }

  spec {
    pod_selector {}

    ingress {
      from {
        namespace_selector {
          match_labels = {
            role = "test"
          }
        }
      }
    }

    policy_types = ["Ingress"]
  }
}

resource "aws_iam_policy" "test" {
  name_prefix = "test"
  path        = "/"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ValidSid",
      "Effect": "Allow",
      "Action": ["es:*"],
      "Resource": ["*"]
    }
  ]
}
EOF
}
