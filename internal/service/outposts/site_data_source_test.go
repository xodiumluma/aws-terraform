// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package outposts_test

import (
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/outposts"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccOutpostsSiteDataSource_id(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_outposts_site.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckSites(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, outposts.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteDataSourceConfig_id(),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrAccountID(dataSourceName, "account_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestMatchResourceAttr(dataSourceName, "id", regexache.MustCompile(`^os-.+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "name", regexache.MustCompile(`^.+$`)),
				),
			},
		},
	})
}

func TestAccOutpostsSiteDataSource_name(t *testing.T) {
	ctx := acctest.Context(t)
	sourceDataSourceName := "data.aws_outposts_site.source"
	dataSourceName := "data.aws_outposts_site.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckSites(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, outposts.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteDataSourceConfig_name(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "account_id", sourceDataSourceName, "account_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", sourceDataSourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "id", sourceDataSourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", sourceDataSourceName, "name"),
				),
			},
		},
	})
}

func testAccSiteDataSourceConfig_id() string {
	return `
data "aws_outposts_sites" "test" {}

data "aws_outposts_site" "test" {
  id = tolist(data.aws_outposts_sites.test.ids)[0]
}
`
}

func testAccSiteDataSourceConfig_name() string {
	return `
data "aws_outposts_sites" "test" {}

data "aws_outposts_site" "source" {
  id = tolist(data.aws_outposts_sites.test.ids)[0]
}

data "aws_outposts_site" "test" {
  name = data.aws_outposts_site.source.name
}
`
}
