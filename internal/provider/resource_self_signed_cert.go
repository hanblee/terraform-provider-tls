package provider

import (
	"crypto/x509"
	"fmt"
	"net"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSelfSignedCert() *schema.Resource {
	s := map[string]*schema.Schema{}

	setCertificateCommonSchema(s)
	setCertificateSubjectSchema(s)

	return &schema.Resource{
		Create:        CreateSelfSignedCert,
		Delete:        deleteCertificate,
		Read:          readCertificate,
		Update:        updateCertificate,
		CustomizeDiff: customizeCertificateDiff,
		Schema:        s,
		Description: "Creates a **self-signed** TLS certificate in " +
			"[PEM (RFC 1421)](https://datatracker.ietf.org/doc/html/rfc1421) format.",
	}
}

func CreateSelfSignedCert(d *schema.ResourceData, meta interface{}) error {
	key, err := parsePrivateKey(d, "private_key_pem", "key_algorithm")
	if err != nil {
		return err
	}

	subjectConfs := d.Get("subject").([]interface{})
	if len(subjectConfs) != 1 {
		return fmt.Errorf("must have exactly one 'subject' block")
	}
	subjectConf, ok := subjectConfs[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("subject block cannot be empty")
	}
	subject := distinguishedNamesFromSubjectAttributes(subjectConf)

	cert := x509.Certificate{
		Subject:               *subject,
		BasicConstraintsValid: true,
	}

	dnsNamesI := d.Get("dns_names").([]interface{})
	for _, nameI := range dnsNamesI {
		cert.DNSNames = append(cert.DNSNames, nameI.(string))
	}
	ipAddressesI := d.Get("ip_addresses").([]interface{})
	for _, ipStrI := range ipAddressesI {
		ip := net.ParseIP(ipStrI.(string))
		if ip == nil {
			return fmt.Errorf("invalid IP address %#v", ipStrI.(string))
		}
		cert.IPAddresses = append(cert.IPAddresses, ip)
	}
	urisI := d.Get("uris").([]interface{})
	for _, uriStrI := range urisI {
		uri, err := url.Parse(uriStrI.(string))
		if err != nil {
			return fmt.Errorf("invalid URI %#v", uriStrI.(string))
		}
		cert.URIs = append(cert.URIs, uri)
	}

	return createCertificate(d, &cert, &cert, toPublicKey(key), key)
}
