# Terraform Provider Uptycs

## Auth ##

This provider allows env var auth as well as provider{} auth inline

```
export UPTYCS_CUSTOMER_ID="your-customer-id"
export UPTYCS_API_SECRET="your-api-secret"
export UPTYCS_API_KEY="your-api-key"
export UPTYCS_HOST=https://test.uptycs.io
```

Inline:

```
provider "uptycs" {
  host = "https://test.uptycs.io"
  customer_id = "your-customer-id"
  api_key = "your-api-key"
  api_secret = "your-api-secret"
}


```

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-uptycs
```

## Test sample configuration

First, bump the version so that its unique and won't pull from the registry:

```shell
$ vim Makefile
```

Next, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory.

```shell
$ cd examples
```

To run this locally you'll need to add a `~/.terraformrc` file with:

```
provider_installation {
  filesystem_mirror {
    path    = "/Users/marcus.young/.terraform.d/plugins"
    include = [
      "github.com/uptycslabs/uptycs",
      "registry.terraform.io/uptycslabs/uptycs",
    ]

  }
  direct {
    exclude = ["uptycslabs/uptycs"]
  }
}
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ rm -rf .terraform .terraform.lock.hcl
$ terraform init
$ export UPTYCS_API_KEY="your-api-key"
$ export UPTYCS_API_SECRET="your-api-secret"
$ terraform plan
```
