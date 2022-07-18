# Terraform Provider Uptycs

## Auth ##

This provider allows env var auth as well as provider{} auth inline

```
export UPTYCS_CUSTOMER_ID=your-customer-id
export UPTYCS_API_SECRET=your-api-secret
export UPTYCS_API_KEY=your-api-key
export UPTYCS_HOST=https://test.uptycs.io
```

Inline:

```
provider "uptycs" {
  host = "https://test.uptycs.io"
  customer_id = "11111111-1111-1111-1111-11111111111"
  api_key = "2222222222222222222222"
  api_secret = "234444444444433333333333222222221111111"
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

To run this locally you'll need to use this for `init.tf`:

```
terraform {
  required_providers {
    uptycs = {
      source  = "github.com/uptycslabs/uptycs" # this source has HOSTNAME to make it unique and valid for local testing
      # source  = "uptycslabs/uptycs" # this is what you'll use when not local, pulling from the registry (with signature checking etc)
      version = "0.0.5"
    }
  }
}

provider "uptycs" {
  host = "https://thor.uptycs.io"
  customer_id = "fda3f46b-c262-439c-bc93-5de6ee6993b6"
}
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ rm -rf .terraform .terraform.lock.hcl
$ terraform init 
$ export UPTYCS_API_KEY="changeme"
$ export UPTYCS_API_SECRET="changeme"
$ terraform plan
```
