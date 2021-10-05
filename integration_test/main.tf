terraform {
  required_version = ">= 1.0.0"
  required_providers {
    statuscake = {
      source  = "marceloboeira/statuscake"
      version = "1.0.0"
    }
  }
}

variable "STATUSCAKE_APIKEY" {
  description = "The StatusCake API key"
  default = "io5N8pQxIeHcVh1fXW7K"
}

provider "statuscake" {
  apikey   = var.STATUSCAKE_APIKEY
}

resource "statuscake_contact_group" "test1" {
  name  = "SeatGeek"
  ping_url     = "http://marceloboeira.com"
}
