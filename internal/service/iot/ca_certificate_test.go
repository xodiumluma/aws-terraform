// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfiot "github.com/hashicorp/terraform-provider-aws/internal/service/iot"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccIoTCACertificate_basic(t *testing.T) {
	ctx := acctest.Context(t)
	caKey := acctest.TLSRSAPrivateKeyPEM(t, 2048)
	caCertificate := acctest.TLSRSAX509SelfSignedCACertificatePEM(t, caKey)
	resourceName := "aws_iot_ca_certificate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCACertificateDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCACertificateConfig_basic(caCertificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCACertificateExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_auto_registration", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "ca_certificate_pem"),
					resource.TestCheckResourceAttr(resourceName, "certificate_mode", "SNI_ONLY"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_version"),
					resource.TestCheckResourceAttrSet(resourceName, "generation_id"),
					resource.TestCheckResourceAttr(resourceName, "registration_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "validity.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "validity.0.not_after"),
					resource.TestCheckResourceAttrSet(resourceName, "validity.0.not_before"),
				),
			},
		},
	})
}

func TestAccIoTCACertificate_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	caKey := acctest.TLSRSAPrivateKeyPEM(t, 2048)
	caCertificate := acctest.TLSRSAX509SelfSignedCACertificatePEM(t, caKey)
	resourceName := "aws_iot_ca_certificate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCACertificateDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCACertificateConfig_basic(caCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCACertificateExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfiot.ResourceCACertificate(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIoTCACertificate_tags(t *testing.T) {
	ctx := acctest.Context(t)
	caKey := acctest.TLSRSAPrivateKeyPEM(t, 2048)
	caCertificate := acctest.TLSRSAX509SelfSignedCACertificatePEM(t, caKey)
	resourceName := "aws_iot_ca_certificate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCACertificateDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCACertificateConfig_tags1(caCertificate, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCACertificateExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccCACertificateConfig_tags2(caCertificate, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCACertificateExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccCACertificateConfig_tags1(caCertificate, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCACertificateExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccIoTCACertificate_defaultMode(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_iot_ca_certificate.test"
	testExternalProviders := map[string]resource.ExternalProvider{
		"tls": {
			Source:            "hashicorp/tls",
			VersionConstraint: "4.0.4",
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		ExternalProviders:        testExternalProviders,
		CheckDestroy:             testAccCheckCACertificateDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCACertificateConfig_defaultMode(false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCACertificateExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
					resource.TestCheckResourceAttr(resourceName, "allow_auto_registration", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "ca_certificate_pem"),
					resource.TestCheckResourceAttr(resourceName, "certificate_mode", "DEFAULT"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_version"),
					resource.TestCheckResourceAttrSet(resourceName, "generation_id"),
					resource.TestCheckResourceAttr(resourceName, "registration_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "validity.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "validity.0.not_after"),
					resource.TestCheckResourceAttrSet(resourceName, "validity.0.not_before"),
				),
			},
			{
				Config: testAccCACertificateConfig_defaultMode(true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCACertificateExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_auto_registration", "true"),
				),
			},
		},
	})
}

func testAccCheckCACertificateExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).IoTConn(ctx)

		_, err := tfiot.FindCACertificateByID(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckCACertificateDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).IoTConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_iot_ca_certificate" {
				continue
			}

			_, err := tfiot.FindCACertificateByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("IoT CA Certificate %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCACertificateConfig_basic(caCertificate string) string {
	return fmt.Sprintf(`
resource "aws_iot_ca_certificate" "test" {
  active                  = true
  allow_auto_registration = true
  ca_certificate_pem      = "%[1]s"
  certificate_mode        = "SNI_ONLY"
}
`, acctest.TLSPEMEscapeNewlines(caCertificate))
}

func testAccCACertificateConfig_tags1(caCertificate, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_iot_ca_certificate" "test" {
  active                  = true
  allow_auto_registration = true
  ca_certificate_pem      = "%[1]s"
  certificate_mode        = "SNI_ONLY"

  tags = {
    %[2]q = %[3]q
  }
}
`, acctest.TLSPEMEscapeNewlines(caCertificate), tagKey1, tagValue1)
}

func testAccCACertificateConfig_tags2(caCertificate, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_iot_ca_certificate" "test" {
  active                  = true
  allow_auto_registration = true
  ca_certificate_pem      = "%[1]s"
  certificate_mode        = "SNI_ONLY"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, acctest.TLSPEMEscapeNewlines(caCertificate), tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccCACertificateConfig_defaultMode(active, allowAutoRegistration bool) string {
	return fmt.Sprintf(`
resource "tls_self_signed_cert" "ca" {
  private_key_pem = tls_private_key.ca.private_key_pem
  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }
  validity_period_hours = 12
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
  is_ca_certificate = true
}

resource "tls_private_key" "ca" {
  algorithm = "RSA"
}

resource "tls_cert_request" "verification" {
  private_key_pem = tls_private_key.verification.private_key_pem
  subject {
    common_name = data.aws_iot_registration_code.test.registration_code
  }
}

resource "tls_private_key" "verification" {
  algorithm = "RSA"
}

resource "tls_locally_signed_cert" "verification" {
  cert_request_pem      = tls_cert_request.verification.cert_request_pem
  ca_private_key_pem    = tls_private_key.ca.private_key_pem
  ca_cert_pem           = tls_self_signed_cert.ca.cert_pem
  validity_period_hours = 12
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "aws_iot_ca_certificate" "test" {
  active                       = %[1]t
  allow_auto_registration      = %[2]t
  ca_certificate_pem           = tls_self_signed_cert.ca.cert_pem
  certificate_mode             = "DEFAULT"
  verification_certificate_pem = tls_locally_signed_cert.verification.cert_pem
}

data "aws_iot_registration_code" "test" {}
`, active, allowAutoRegistration)
}
