// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package licensemanager_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/licensemanager"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/envvar"
	tflicensemanager "github.com/hashicorp/terraform-provider-aws/internal/service/licensemanager"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	homeRegionKey = "TF_AWS_LICENSE_MANAGER_GRANT_HOME_REGION"
	licenseARNKey = "TF_AWS_LICENSE_MANAGER_GRANT_LICENSE_ARN"
	principalKey  = "TF_AWS_LICENSE_MANAGER_GRANT_PRINCIPAL"
)

const (
	envVarHomeRegionError    = "The region where the license has been imported into the current account."
	envVarLicenseARNKeyError = "ARN for a license imported into the current account."
	envVarPrincipalKeyError  = "ARN of a principal to share the license with. Either a root user, Organization, or Organizational Unit."
)

func TestAccLicenseManagerGrant_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"grant": {
			"basic":      testAccGrant_basic,
			"disappears": testAccGrant_disappears,
			"name":       testAccGrant_name,
		},
		"grant_accepter": {
			"basic":      testAccGrantAccepter_basic,
			"disappears": testAccGrantAccepter_disappears,
		},
		"grant_data_source": {
			"basic": testAccGrantsDataSource_basic,
			"empty": testAccGrantsDataSource_noMatch,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}

func testAccGrant_basic(t *testing.T) {
	ctx := acctest.Context(t)
	licenseARN := envvar.SkipIfEmpty(t, licenseARNKey, envVarLicenseARNKeyError)
	principal := envvar.SkipIfEmpty(t, principalKey, envVarPrincipalKeyError)
	homeRegion := envvar.SkipIfEmpty(t, homeRegionKey, envVarHomeRegionError)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_licensemanager_grant.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, licensemanager.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_basic(licenseARN, rName, principal, homeRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					acctest.MatchResourceAttrGlobalARN(resourceName, "arn", "license-manager", regexache.MustCompile(`grant:g-.+`)),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_operations.*", "ListPurchasedLicenses"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_operations.*", "CheckoutLicense"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_operations.*", "CheckInLicense"),
					resource.TestCheckTypeSetElemAttr(resourceName, "allowed_operations.*", "ExtendConsumptionLicense"),
					resource.TestCheckResourceAttr(resourceName, "home_region", homeRegion),
					resource.TestCheckResourceAttr(resourceName, "license_arn", licenseARN),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "parent_arn"),
					resource.TestCheckResourceAttr(resourceName, "principal", principal),
					resource.TestCheckResourceAttr(resourceName, "status", "PENDING_ACCEPT"),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
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

func testAccGrant_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	licenseARN := envvar.SkipIfEmpty(t, licenseARNKey, envVarLicenseARNKeyError)
	principal := envvar.SkipIfEmpty(t, principalKey, envVarPrincipalKeyError)
	homeRegion := envvar.SkipIfEmpty(t, homeRegionKey, envVarHomeRegionError)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_licensemanager_grant.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, licensemanager.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_basic(licenseARN, rName, principal, homeRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tflicensemanager.ResourceGrant(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccGrant_name(t *testing.T) {
	ctx := acctest.Context(t)
	licenseARN := envvar.SkipIfEmpty(t, licenseARNKey, envVarLicenseARNKeyError)
	principal := envvar.SkipIfEmpty(t, principalKey, envVarPrincipalKeyError)
	homeRegion := envvar.SkipIfEmpty(t, homeRegionKey, envVarHomeRegionError)
	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_licensemanager_grant.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, licensemanager.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_basic(licenseARN, rName1, principal, homeRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGrantConfig_basic(licenseARN, rName2, principal, homeRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func testAccCheckGrantExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No License Manager License Configuration ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LicenseManagerConn(ctx)

		out, err := tflicensemanager.FindGrantByARN(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		if out == nil {
			return fmt.Errorf("Grant %q does not exist", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckGrantDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).LicenseManagerConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_licensemanager_grant" {
				continue
			}

			_, err := tflicensemanager.FindGrantByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("License Manager Grant %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccGrantConfig_basic(licenseARN, rName, principal, homeRegion string) string {
	return acctest.ConfigCompose(
		acctest.ConfigRegionalProvider(homeRegion),
		fmt.Sprintf(`
data "aws_licensemanager_received_license" "test" {
  license_arn = %[1]q
}

locals {
  allowed_operations = [for i in data.aws_licensemanager_received_license.test.received_metadata[0].allowed_operations : i if i != "CreateGrant"]
}

resource "aws_licensemanager_grant" "test" {
  name               = %[2]q
  allowed_operations = local.allowed_operations
  license_arn        = data.aws_licensemanager_received_license.test.license_arn
  principal          = %[3]q
}
`, licenseARN, rName, principal),
	)
}
