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

## Using provider locally

While developing you'll likely want to use this provider locally.

If you need to use a temporary branch for the `uptycs-client-go`:

 * push the commit for the uptycs-client-go to the uptycslabs upstream
 * note the SHA (hint: `git log -1`)
 * pull the temporary sha via `go get github.com/uptycslabs/uptycs-client-go@{sha from git log -1}`


Next, Bump the `VERSION` in the `Makefile`, then build and install the provider with `$ make install`, use this for your init configuration:

```
terraform {
  required_providers {
    uptycs = {
      source  = "github.com/uptycslabs/uptycs"
      version = "0.0.5" # the version you bumped above
    }
  }
}
```

Do not use the s3 state to prevent corruption. Use a local state file.
Make sure to remove cached versions, state, etc: `rm .terraform.lock.hcl; rm -rf .terraform; rm terraform.tfstate*`

The github.com namespace is what lets you run this locally. When running the officially published one from the terraform registry note the lack of that part.

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory.

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```
