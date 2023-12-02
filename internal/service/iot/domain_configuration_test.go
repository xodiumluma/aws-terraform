// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/iot"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfiot "github.com/hashicorp/terraform-provider-aws/internal/service/iot"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccIoTDomainConfiguration_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rootDomain := acctest.ACMCertificateDomainFromEnv(t)
	domain := acctest.ACMCertificateRandomSubDomain(rootDomain)
	resourceName := "aws_iot_domain_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDomainConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigurationConfig_basic(rName, rootDomain, domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDomainConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", domain),
					resource.TestCheckResourceAttr(resourceName, "domain_type", "CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "server_certificate_arns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "DATA"),
					resource.TestCheckResourceAttr(resourceName, "status", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tls_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "tls_config.0.security_policy"),
					resource.TestCheckResourceAttr(resourceName, "validation_certificate_arn", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIoTDomainConfiguration_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rootDomain := acctest.ACMCertificateDomainFromEnv(t)
	domain := acctest.ACMCertificateRandomSubDomain(rootDomain)
	resourceName := "aws_iot_domain_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDomainConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigurationConfig_basic(rName, rootDomain, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfiot.ResourceDomainConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIoTDomainConfiguration_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rootDomain := acctest.ACMCertificateDomainFromEnv(t)
	domain := acctest.ACMCertificateRandomSubDomain(rootDomain)
	resourceName := "aws_iot_domain_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDomainConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigurationConfig_tags1(rName, rootDomain, domain, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDomainConfigurationConfig_tags2(rName, rootDomain, domain, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccDomainConfigurationConfig_tags1(rName, rootDomain, domain, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccIoTDomainConfiguration_update(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rootDomain := acctest.ACMCertificateDomainFromEnv(t)
	domain := acctest.ACMCertificateRandomSubDomain(rootDomain)
	resourceName := "aws_iot_domain_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDomainConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigurationConfig_securityPolicy(rName, rootDomain, domain, "IoTSecurityPolicy_TLS13_1_3_2022_10", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "authorizer_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_config.0.allow_authorizer_override", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tls_config.0.security_policy", "IoTSecurityPolicy_TLS13_1_3_2022_10"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDomainConfigurationConfig_securityPolicy(rName, rootDomain, domain, "IoTSecurityPolicy_TLS13_1_2_2022_10", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "authorizer_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_config.0.allow_authorizer_override", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tls_config.0.security_policy", "IoTSecurityPolicy_TLS13_1_2_2022_10"),
				),
			},
		},
	})
}

func TestAccIoTDomainConfiguration_awsManaged(t *testing.T) { // nosemgrep:ci.aws-in-func-name
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_iot_domain_configuration.test"

	acctest.SkipIfEnvVarNotSet(t, "IOT_DOMAIN_CONFIGURATION_TEST_AWS_MANAGED")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfigurationConfig_awsManaged(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDomainConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_config.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttr(resourceName, "domain_type", "AWS_MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "server_certificate_arns.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "DATA"),
					resource.TestCheckResourceAttr(resourceName, "status", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tls_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "tls_config.0.security_policy"),
					resource.TestCheckResourceAttr(resourceName, "validation_certificate_arn", ""),
				),
			},
		},
	})
}

func testAccCheckDomainConfigurationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).IoTConn(ctx)

		_, err := tfiot.FindDomainConfigurationByName(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckDomainConfigurationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).IoTConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_iot_domain_configuration" {
				continue
			}

			_, err := tfiot.FindDomainConfigurationByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("IoT Domain Configuration %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccDomainConfigurationConfig_base(rootDomain, domain string) string {
	return fmt.Sprintf(`
resource "aws_acm_certificate" "test" {
  domain_name       = %[2]q
  validation_method = "DNS"
}

data "aws_route53_zone" "test" {
  name         = %[1]q
  private_zone = false
}

resource "aws_route53_record" "test" {
  allow_overwrite = true
  name            = tolist(aws_acm_certificate.test.domain_validation_options)[0].resource_record_name
  records         = [tolist(aws_acm_certificate.test.domain_validation_options)[0].resource_record_value]
  ttl             = 60
  type            = tolist(aws_acm_certificate.test.domain_validation_options)[0].resource_record_type
  zone_id         = data.aws_route53_zone.test.zone_id
}

resource "aws_acm_certificate_validation" "test" {
  depends_on = [aws_route53_record.test]

  certificate_arn = aws_acm_certificate.test.arn
}
`, rootDomain, domain)
}

func testAccDomainConfigurationConfig_basic(rName, rootDomain, domain string) string {
	return acctest.ConfigCompose(testAccDomainConfigurationConfig_base(rootDomain, domain), fmt.Sprintf(`
resource "aws_iot_domain_configuration" "test" {
  depends_on = [aws_acm_certificate_validation.test]

  name                    = %[1]q
  domain_name             = %[2]q
  server_certificate_arns = [aws_acm_certificate.test.arn]
}
`, rName, domain))
}

func testAccDomainConfigurationConfig_tags1(rName, rootDomain, domain, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccDomainConfigurationConfig_base(rootDomain, domain), fmt.Sprintf(`
resource "aws_iot_domain_configuration" "test" {
  depends_on = [aws_acm_certificate_validation.test]

  name                    = %[1]q
  domain_name             = %[2]q
  server_certificate_arns = [aws_acm_certificate.test.arn]

  tags = {
    %[3]q = %[4]q
  }
}
`, rName, domain, tagKey1, tagValue1))
}

func testAccDomainConfigurationConfig_tags2(rName, rootDomain, domain, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccDomainConfigurationConfig_base(rootDomain, domain), fmt.Sprintf(`
resource "aws_iot_domain_configuration" "test" {
  depends_on = [aws_acm_certificate_validation.test]

  name                    = %[1]q
  domain_name             = %[2]q
  server_certificate_arns = [aws_acm_certificate.test.arn]

  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, rName, domain, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccDomainConfigurationConfig_securityPolicy(rName, rootDomain, domain, securityPolicy string, allowAuthorizerOverride bool) string {
	return acctest.ConfigCompose(testAccAuthorizerConfig_basic(rName), testAccDomainConfigurationConfig_base(rootDomain, domain), fmt.Sprintf(`
resource "aws_iot_domain_configuration" "test" {
  depends_on = [aws_acm_certificate_validation.test]

  authorizer_config {
    allow_authorizer_override = %[4]t
    default_authorizer_name   = aws_iot_authorizer.test.name
  }

  name                    = %[1]q
  domain_name             = %[2]q
  server_certificate_arns = [aws_acm_certificate.test.arn]

  tls_config {
    security_policy = %[3]q
  }
}
`, rName, domain, securityPolicy, allowAuthorizerOverride))
}

func testAccDomainConfigurationConfig_awsManaged(rName string) string { // nosemgrep:ci.aws-in-func-name
	return fmt.Sprintf(`
resource "aws_iot_domain_configuration" "test" {
  name = %[1]q
}
`, rName)
}
