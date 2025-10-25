provider "aws" {
  region = "us-east-1"

  default_tags {
    tags = {
      region            = var.aws_region
      app-id            = var.app_id
      environment       = var.environment
      notification_mail = var.notification_mail
    }
  }
}
