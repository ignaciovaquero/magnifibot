terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.73"
    }
  }
}

module "name" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "magnifibot-dev-handle-telegram-command"
  description   = "Handle Telegram command"
  runtime       = "go1.x"

}
