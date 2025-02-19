---
layout: "tls"
page_title: "TLS: tls_cert_request"
description: |-
  Creates a PEM-encoded certificate request.
---

# tls\_cert\_request

Generates a *Certificate Signing Request* (CSR) in PEM format, which is the
typical format used to request a certificate from a certificate authority.

This resource is intended to be used in conjunction with a Terraform provider
for a particular certificate authority in order to provision a new certificate.
This is a *logical resource*, so it contributes only to the current Terraform
state and does not create any external managed resources.

~> **Compatibility Note** From Terraform 0.7.0 to 0.7.4 this resource was
converted to a data source, and the resource form of it was deprecated. This
turned out to be a design error since a cert request includes a random number
in the form of the signature nonce, and so the data source form of this
resource caused non-convergent configuration. The data source form is no longer
supported as of Terraform 0.7.5 and any users should return to using the
resource form.

## Example Usage

```hcl
resource "tls_cert_request" "example" {
  key_algorithm   = "ECDSA"
  private_key_pem = "${file("private_key.pem")}"

  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `key_algorithm` - (Required) The name of the algorithm for the key provided
in `private_key_pem`.

* `private_key_pem` - (Required) PEM-encoded private key data. This can be
read from a separate file using the ``file`` interpolation function. Only
an irreversable secure hash of the private key will be stored in the Terraform
state.

* `subject` - (Required) The subject for which a certificate is being requested. This is
a nested configuration block whose structure is described below.

* `dns_names` - (Optional) List of DNS names for which a certificate is being requested.

* `ip_addresses` - (Optional) List of IP addresses for which a certificate is being requested.

* `uris` - (Optional) List of URIs for which a certificate is being requested.

The nested `subject` block accepts the following arguments, all optional, with their meaning
corresponding to the similarly-named attributes defined in
[RFC5280](https://tools.ietf.org/html/rfc5280#section-4.1.2.4):

* `common_name` (string)

* `organization` (string)

* `organizational_unit` (string)

* `street_address` (list of strings)

* `locality` (string)

* `province` (string)

* `country` (string)

* `postal_code` (string)

* `serial_number` (string)

-> **Note:** Versions of this provider prior to 1.2.0 may generate a
certificate request which cannot be validated by your Certificate Authority if
the `*` character is used in any of the `subject` fields, for example as part
of the `common_name` when generating a request for a wildcard certificate.
Strings containing a `*` and passed to the `dns_names` argument are encoded
correctly.

## Attributes Reference

The following attributes are exported:

* `cert_request_pem` - The certificate request data in PEM format.
