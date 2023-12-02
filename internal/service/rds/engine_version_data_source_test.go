// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rds_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccRDSEngineVersionDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_rds_engine_version.test"
	engine := "oracle-ee"
	version := "19.0.0.0.ru-2020-07.rur-2020-07.r1"
	paramGroup := "oracle-ee-19"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_basic(engine, version, paramGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "engine", engine),
					resource.TestCheckResourceAttr(dataSourceName, "version", version),
					resource.TestCheckResourceAttr(dataSourceName, "parameter_group_family", paramGroup),

					resource.TestCheckResourceAttrSet(dataSourceName, "default_character_set"),
					resource.TestCheckResourceAttrSet(dataSourceName, "engine_description"),
					resource.TestMatchResourceAttr(dataSourceName, "exportable_log_types.#", regexache.MustCompile(`^[1-9][0-9]*`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestMatchResourceAttr(dataSourceName, "supported_character_sets.#", regexache.MustCompile(`^[1-9][0-9]*`)),
					resource.TestMatchResourceAttr(dataSourceName, "supported_feature_names.#", regexache.MustCompile(`^[1-9][0-9]*`)),
					resource.TestMatchResourceAttr(dataSourceName, "supported_modes.#", regexache.MustCompile(`^[0-9]*`)),
					resource.TestMatchResourceAttr(dataSourceName, "supported_timezones.#", regexache.MustCompile(`^[0-9]*`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "supports_global_databases"),
					resource.TestCheckResourceAttrSet(dataSourceName, "supports_log_exports_to_cloudwatch"),
					resource.TestCheckResourceAttrSet(dataSourceName, "supports_parallel_query"),
					resource.TestCheckResourceAttrSet(dataSourceName, "supports_read_replica"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version_description"),
				),
			},
		},
	})
}

func TestAccRDSEngineVersionDataSource_upgradeTargets(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_rds_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_upgradeTargets(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "valid_upgrade_targets.#", regexache.MustCompile(`^[1-9][0-9]*`)),
				),
			},
		},
	})
}

func TestAccRDSEngineVersionDataSource_preferred(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_rds_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_preferred(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "version", "8.0.32"),
				),
			},
		},
	})
}

func TestAccRDSEngineVersionDataSource_defaultOnlyImplicit(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_rds_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_defaultOnlyImplicit(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
				),
			},
		},
	})
}

func TestAccRDSEngineVersionDataSource_defaultOnlyExplicit(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_rds_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_defaultOnlyExplicit(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "version", regexache.MustCompile(`^8\.0\.`)),
				),
			},
		},
	})
}

func TestAccRDSEngineVersionDataSource_includeAll(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_rds_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_includeAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "version", "8.0.20"),
				),
			},
		},
	})
}

func TestAccRDSEngineVersionDataSource_filter(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_rds_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccEngineVersionPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, rds.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEngineVersionDataSourceConfig_filter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "version", "13.9"),
					resource.TestCheckResourceAttr(dataSourceName, "supported_modes.0", "serverless"),
				),
			},
		},
	})
}

func testAccEngineVersionPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).RDSConn(ctx)

	input := &rds.DescribeDBEngineVersionsInput{
		Engine:      aws.String("mysql"),
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

func testAccEngineVersionDataSourceConfig_basic(engine, version, paramGroup string) string {
	return fmt.Sprintf(`
data "aws_rds_engine_version" "test" {
  engine                 = %[1]q
  version                = %[2]q
  parameter_group_family = %[3]q
}
`, engine, version, paramGroup)
}

func testAccEngineVersionDataSourceConfig_upgradeTargets() string {
	return `
data "aws_rds_engine_version" "test" {
  engine  = "mysql"
  version = "8.0.32"
}
`
}

func testAccEngineVersionDataSourceConfig_preferred() string {
	return `
data "aws_rds_engine_version" "test" {
  engine             = "mysql"
  preferred_versions = ["85.9.12", "8.0.32", "8.0.31"]
}
`
}

func testAccEngineVersionDataSourceConfig_defaultOnlyImplicit() string {
	return `
data "aws_rds_engine_version" "test" {
  engine = "mysql"
}
`
}

func testAccEngineVersionDataSourceConfig_defaultOnlyExplicit() string {
	return `
data "aws_rds_engine_version" "test" {
  engine       = "mysql"
  version      = "8.0"
  default_only = true
}
`
}

func testAccEngineVersionDataSourceConfig_includeAll() string {
	return `
data "aws_rds_engine_version" "test" {
  engine      = "mysql"
  version     = "8.0.20"
  include_all = true
}
`
}

func testAccEngineVersionDataSourceConfig_filter() string {
	return `
data "aws_rds_engine_version" "test" {
  engine      = "aurora-postgresql"
  version     = "13.9"
  include_all = true

  filter {
    name   = "engine-mode"
    values = ["serverless"]
  }
}
`
}
