// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccDocDBEngineVersionDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_docdb_engine_version.test"
	engine := "docdb"
	version := "3.6.0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_basic(engine, version),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "engine", engine),
					resource.TestCheckResourceAttr(dataSourceName, "version", version),
					resource.TestCheckResourceAttrSet(dataSourceName, "engine_description"),
					resource.TestMatchResourceAttr(dataSourceName, "exportable_log_types.#", regexache.MustCompile(`^[1-9][0-9]*`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameter_group_family"),
					resource.TestCheckResourceAttrSet(dataSourceName, "supports_log_exports_to_cloudwatch"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_description"),
				),
			},
		},
	})
}

func TestAccDocDBEngineVersionDataSource_preferred(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_docdb_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_preferred(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "version", "3.6.0"),
				),
			},
		},
	})
}

func TestAccDocDBEngineVersionDataSource_defaultOnly(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_docdb_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_defaultOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "engine", "docdb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
				),
			},
		},
	})
}

func testAccEngineVersionPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).DocDBConn(ctx)

	input := &docdb.DescribeDBEngineVersionsInput{
		Engine:      aws.String("docdb"),
		DefaultOnly: aws.Bool(true),
	}

	_, err := conn.DescribeDBEngineVersionsWithContext(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccEngineVersionDataSourceConfig_basic(engine, version string) string {
	return fmt.Sprintf(`
data "aws_docdb_engine_version" "test" {
  engine  = %[1]q
  version = %[2]q
}
`, engine, version)
}

func testAccEngineVersionDataSourceConfig_preferred() string {
	return `
data "aws_docdb_engine_version" "test" {
  preferred_versions = ["34.6.1", "3.6.0", "2.6.0"]
}
`
}

func testAccEngineVersionDataSourceConfig_defaultOnly() string {
	return `
data "aws_docdb_engine_version" "test" {}
`
}
