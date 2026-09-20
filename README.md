# Terraform Provider for voidauth

A Terraform provider for managing voidauth.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25 (for development)

## Installation

### Using the Provider

```terraform
terraform {
  required_providers {
    voidauth = {
      source = "SimonPrinz/voidauth"
    }
  }
}

provider "voidauth" {
  url      = "http://localhost:3000"
  # username = "auth_admin"
  # password = "password"
}
```

## Documentation

Full documentation is available in the [docs/](./docs/) directory:

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
- API Client sourced from [voidauth](https://github.com/voidauth/voidauth)
