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
