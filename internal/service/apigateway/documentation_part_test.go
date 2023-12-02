// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apigateway_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfapigateway "github.com/hashicorp/terraform-provider-aws/internal/service/apigateway"
)

func TestAccAPIGatewayDocumentationPart_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var conf apigateway.DocumentationPart

	rString := sdkacctest.RandString(8)
	apiName := fmt.Sprintf("tf-acc-test_api_doc_part_basic_%s", rString)
	properties := `{"description":"Terraform Acceptance Test"}`
	uProperties := `{"description":"Terraform Acceptance Test Updated"}`

	resourceName := "aws_api_gateway_documentation_part.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDocumentationPartDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDocumentationPartConfig_basic(apiName, strconv.Quote(properties)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDocumentationPartExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "location.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "location.0.type", "API"),
					resource.TestCheckResourceAttr(resourceName, "properties", properties),
					resource.TestCheckResourceAttrSet(resourceName, "rest_api_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDocumentationPartConfig_basic(apiName, strconv.Quote(uProperties)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDocumentationPartExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "location.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "location.0.type", "API"),
					resource.TestCheckResourceAttr(resourceName, "properties", uProperties),
					resource.TestCheckResourceAttrSet(resourceName, "rest_api_id"),
				),
			},
		},
	})
}

func TestAccAPIGatewayDocumentationPart_method(t *testing.T) {
	ctx := acctest.Context(t)
	var conf apigateway.DocumentationPart

	rString := sdkacctest.RandString(8)
	apiName := fmt.Sprintf("tf-acc-test_api_doc_part_method_%s", rString)
	properties := `{"description":"Terraform Acceptance Test"}`
	uProperties := `{"description":"Terraform Acceptance Test Updated"}`

	resourceName := "aws_api_gateway_documentation_part.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDocumentationPartDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDocumentationPartConfig_method(apiName, strconv.Quote(properties)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDocumentationPartExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "location.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "location.0.type", "METHOD"),
					resource.TestCheckResourceAttr(resourceName, "location.0.method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "location.0.path", "/terraform-acc-test"),
					resource.TestCheckResourceAttr(resourceName, "properties", properties),
					resource.TestCheckResourceAttrSet(resourceName, "rest_api_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDocumentationPartConfig_method(apiName, strconv.Quote(uProperties)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDocumentationPartExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "location.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "location.0.type", "METHOD"),
					resource.TestCheckResourceAttr(resourceName, "location.0.method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "location.0.path", "/terraform-acc-test"),
					resource.TestCheckResourceAttr(resourceName, "properties", uProperties),
					resource.TestCheckResourceAttrSet(resourceName, "rest_api_id"),
				),
			},
		},
	})
}

func TestAccAPIGatewayDocumentationPart_responseHeader(t *testing.T) {
	ctx := acctest.Context(t)
	var conf apigateway.DocumentationPart

	rString := sdkacctest.RandString(8)
	apiName := fmt.Sprintf("tf-acc-test_api_doc_part_resp_header_%s", rString)
	properties := `{"description":"Terraform Acceptance Test"}`
	uProperties := `{"description":"Terraform Acceptance Test Updated"}`

	resourceName := "aws_api_gateway_documentation_part.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDocumentationPartDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDocumentationPartConfig_responseHeader(apiName, strconv.Quote(properties)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDocumentationPartExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "location.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "location.0.type", "RESPONSE_HEADER"),
					resource.TestCheckResourceAttr(resourceName, "location.0.method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "location.0.name", "tfacc"),
					resource.TestCheckResourceAttr(resourceName, "location.0.path", "/terraform-acc-test"),
					resource.TestCheckResourceAttr(resourceName, "location.0.status_code", "200"),
					resource.TestCheckResourceAttr(resourceName, "properties", properties),
					resource.TestCheckResourceAttrSet(resourceName, "rest_api_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDocumentationPartConfig_responseHeader(apiName, strconv.Quote(uProperties)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDocumentationPartExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "location.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "location.0.type", "RESPONSE_HEADER"),
					resource.TestCheckResourceAttr(resourceName, "location.0.method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "location.0.name", "tfacc"),
					resource.TestCheckResourceAttr(resourceName, "location.0.path", "/terraform-acc-test"),
					resource.TestCheckResourceAttr(resourceName, "location.0.status_code", "200"),
					resource.TestCheckResourceAttr(resourceName, "properties", uProperties),
					resource.TestCheckResourceAttrSet(resourceName, "rest_api_id"),
				),
			},
		},
	})
}

func TestAccAPIGatewayDocumentationPart_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var conf apigateway.DocumentationPart

	rString := sdkacctest.RandString(8)
	apiName := fmt.Sprintf("tf-acc-test_api_doc_part_basic_%s", rString)
	properties := `{"description":"Terraform Acceptance Test"}`

	resourceName := "aws_api_gateway_documentation_part.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDocumentationPartDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDocumentationPartConfig_basic(apiName, strconv.Quote(properties)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDocumentationPartExists(ctx, resourceName, &conf),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfapigateway.ResourceDocumentationPart(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckDocumentationPartExists(ctx context.Context, n string, res *apigateway.DocumentationPart) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API Gateway Documentation Part ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn(ctx)

		apiId, id, err := tfapigateway.DecodeDocumentationPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		req := &apigateway.GetDocumentationPartInput{
			DocumentationPartId: aws.String(id),
			RestApiId:           aws.String(apiId),
		}
		docPart, err := conn.GetDocumentationPartWithContext(ctx, req)
		if err != nil {
			return err
		}

		*res = *docPart

		return nil
	}
}

func testAccCheckDocumentationPartDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_api_gateway_documentation_part" {
				continue
			}

			apiId, id, err := tfapigateway.DecodeDocumentationPartID(rs.Primary.ID)
			if err != nil {
				return err
			}

			req := &apigateway.GetDocumentationPartInput{
				DocumentationPartId: aws.String(id),
				RestApiId:           aws.String(apiId),
			}
			_, err = conn.GetDocumentationPartWithContext(ctx, req)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
					return nil
				}
				return err
			}

			return fmt.Errorf("API Gateway Documentation Part %q still exists.", rs.Primary.ID)
		}
		return nil
	}
}

func testAccDocumentationPartConfig_basic(apiName, properties string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_documentation_part" "test" {
  location {
    type = "API"
  }
  properties  = %s
  rest_api_id = aws_api_gateway_rest_api.test.id
}

resource "aws_api_gateway_rest_api" "test" {
  name = "%s"
}
`, properties, apiName)
}

func testAccDocumentationPartConfig_method(apiName, properties string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_documentation_part" "test" {
  location {
    type   = "METHOD"
    method = "GET"
    path   = "/terraform-acc-test"
  }
  properties  = %s
  rest_api_id = aws_api_gateway_rest_api.test.id
}

resource "aws_api_gateway_rest_api" "test" {
  name = "%s"
}
`, properties, apiName)
}

func testAccDocumentationPartConfig_responseHeader(apiName, properties string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_documentation_part" "test" {
  location {
    type        = "RESPONSE_HEADER"
    method      = "GET"
    name        = "tfacc"
    path        = "/terraform-acc-test"
    status_code = "200"
  }
  properties  = %s
  rest_api_id = aws_api_gateway_rest_api.test.id
}

resource "aws_api_gateway_rest_api" "test" {
  name = "%s"
}
`, properties, apiName)
}
