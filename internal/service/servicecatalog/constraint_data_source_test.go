// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicecatalog_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/servicecatalog"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccServiceCatalogConstraintDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_servicecatalog_constraint.test"
	dataSourceName := "data.aws_servicecatalog_constraint.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, servicecatalog.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConstraintDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConstraintDataSourceConfig_basic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConstraintExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "accept_language", dataSourceName, "accept_language"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "owner", dataSourceName, "owner"),
					resource.TestCheckResourceAttrPair(resourceName, "parameters", dataSourceName, "parameters"),
					resource.TestCheckResourceAttrPair(resourceName, "portfolio_id", dataSourceName, "portfolio_id"),
					resource.TestCheckResourceAttrPair(resourceName, "product_id", dataSourceName, "product_id"),
					resource.TestCheckResourceAttrPair(resourceName, "status", dataSourceName, "status"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
				),
			},
		},
	})
}

func testAccConstraintDataSourceConfig_basic(rName, description string) string {
	return acctest.ConfigCompose(testAccConstraintConfig_basic(rName, description), `
data "aws_servicecatalog_constraint" "test" {
  id = aws_servicecatalog_constraint.test.id
}
`)
}
