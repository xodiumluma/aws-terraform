// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/schemas"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfschemas "github.com/hashicorp/terraform-provider-aws/internal/service/schemas"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	testAccSchemaContent = `
{
  "openapi": "3.0.0",
  "info": {
    "version": "1.0.0",
    "title": "Event"
  },
  "paths": {},
  "components": {
    "schemas": {
      "Event": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string"
          }
        }
      }
    }
  }
}
`

	testAccSchemaContentUpdated = `
{
  "openapi": "3.0.0",
  "info": {
    "version": "2.0.0",
    "title": "Event"
  },
  "paths": {},
  "components": {
    "schemas": {
      "Event": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string"
          },
          "created_at": {
            "type": "string",
            "format": "date-time"
          }
        }
      }
    }
  }
}
`

	testAccJSONSchemaContent = `
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/product.schema.json",
  "title": "Event",
  "description": "An generic example",
  "type": "object",
  "properties": {
    "name": {
      "description": "The unique identifier for a product",
      "type": "string"
    },
    "created_at": {
      "description": "Date-time format",
      "type": "string",
      "format": "date-time"
    }
  },
  "required": [ "name" ]
}
`
)

func TestAccSchemasSchema_openAPI3(t *testing.T) {
	ctx := acctest.Context(t)
	var v schemas.DescribeSchemaOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_schemas_schema.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, schemas.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, schemas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSchemaDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "schemas", fmt.Sprintf("schema/%s/%s", rName, rName)),
					resource.TestCheckResourceAttrSet(resourceName, "content"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "registry_name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "OpenApi3"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "version_created_date"),
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

func TestAccSchemasSchema_jsonSchemaDraftv4(t *testing.T) {
	ctx := acctest.Context(t)
	var v schemas.DescribeSchemaOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_schemas_schema.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, schemas.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, schemas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSchemaDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaConfig_jsonSchema(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "schemas", fmt.Sprintf("schema/%s/%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "content", testAccJSONSchemaContent),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "registry_name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "JSONSchemaDraft4"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "version_created_date"),
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

func TestAccSchemasSchema_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v schemas.DescribeSchemaOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_schemas_schema.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, schemas.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, schemas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSchemaDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfschemas.ResourceSchema(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSchemasSchema_contentDescription(t *testing.T) {
	ctx := acctest.Context(t)
	var v schemas.DescribeSchemaOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_schemas_schema.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, schemas.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, schemas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSchemaDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaConfig_contentDescription(rName, testAccSchemaContent, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "content", testAccSchemaContent),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSchemaConfig_contentDescription(rName, testAccSchemaContentUpdated, "description2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "content", testAccSchemaContentUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
					resource.TestCheckResourceAttr(resourceName, "version", "2"),
				),
			},
			{
				Config: testAccSchemaConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "version", "3"),
				),
			},
		},
	})
}

func TestAccSchemasSchema_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v schemas.DescribeSchemaOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_schemas_schema.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, schemas.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, schemas.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSchemaDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
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
				Config: testAccSchemaConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccSchemaConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckSchemaDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).SchemasConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_schemas_schema" {
				continue
			}

			name, registryName, err := tfschemas.SchemaParseResourceID(rs.Primary.ID)

			if err != nil {
				return err
			}

			_, err = tfschemas.FindSchemaByNameAndRegistryName(ctx, conn, name, registryName)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("EventBridge Schemas Schema %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckSchemaExists(ctx context.Context, n string, v *schemas.DescribeSchemaOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EventBridge Schemas Schema ID is set")
		}

		name, registryName, err := tfschemas.SchemaParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SchemasConn(ctx)

		output, err := tfschemas.FindSchemaByNameAndRegistryName(ctx, conn, name, registryName)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccSchemaConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_schemas_registry" "test" {
  name = %[1]q
}

resource "aws_schemas_schema" "test" {
  name          = %[1]q
  registry_name = aws_schemas_registry.test.name
  type          = "OpenApi3"
  content       = %[2]q
}
`, rName, testAccSchemaContent)
}

func testAccSchemaConfig_jsonSchema(rName string) string {
	return fmt.Sprintf(`
resource "aws_schemas_registry" "test" {
  name = %[1]q
}

resource "aws_schemas_schema" "test" {
  name          = %[1]q
  registry_name = aws_schemas_registry.test.name
  type          = "JSONSchemaDraft4"
  content       = %[2]q
}
`, rName, testAccJSONSchemaContent)
}

func testAccSchemaConfig_contentDescription(rName, content, description string) string {
	return fmt.Sprintf(`
resource "aws_schemas_registry" "test" {
  name = %[1]q
}

resource "aws_schemas_schema" "test" {
  name          = %[1]q
  registry_name = aws_schemas_registry.test.name
  type          = "OpenApi3"
  content       = %[2]q
  description   = %[3]q
}
`, rName, content, description)
}

func testAccSchemaConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_schemas_registry" "test" {
  name = %[1]q
}

resource "aws_schemas_schema" "test" {
  name          = %[1]q
  registry_name = aws_schemas_registry.test.name
  type          = "OpenApi3"
  content       = %[2]q

  tags = {
    %[3]q = %[4]q
  }
}
`, rName, testAccSchemaContent, tagKey1, tagValue1)
}

func testAccSchemaConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`	
resource "aws_schemas_registry" "test" {
  name = %[1]q
}

resource "aws_schemas_schema" "test" {
  name          = %[1]q
  registry_name = aws_schemas_registry.test.name
  type          = "OpenApi3"
  content       = %[2]q

  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, rName, testAccSchemaContent, tagKey1, tagValue1, tagKey2, tagValue2)
}
