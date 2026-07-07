terraform {

  required_providers {

    thales = {

      source  = "provider/thales"

      version = "1.0.0"

    }

  }

}

provider "thales" {

  endpoint = "http://localhost:8080"

}

resource "thales_keystore" "payment" {

  name = "payment-keystore"

}