// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apigateway_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/apigateway"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfapigateway "github.com/hashicorp/terraform-provider-aws/internal/service/apigateway"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccAPIGatewayMethodResponse_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var conf apigateway.MethodResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_method_response.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMethodResponseDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMethodResponseConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMethodResponseExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "response_models.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "response_models.application/json", "Error"),
					resource.TestCheckResourceAttr(resourceName, "response_parameters.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "response_parameters.method.response.header.Content-Type", "true"),
					resource.TestCheckResourceAttr(resourceName, "status_code", "400"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccMethodResponseImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccMethodResponseConfig_update(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMethodResponseExists(ctx, resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "response_models.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "response_models.application/json", "Empty"),
					resource.TestCheckResourceAttr(resourceName, "response_parameters.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "response_parameters.method.response.header.Host", "false"),
					resource.TestCheckResourceAttr(resourceName, "status_code", "400"),
				),
			},
		},
	})
}

func TestAccAPIGatewayMethodResponse_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var conf apigateway.MethodResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_api_gateway_method_response.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckAPIGatewayTypeEDGE(t) },
		ErrorCheck:               acctest.ErrorCheck(t, apigateway.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMethodResponseDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMethodResponseConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMethodResponseExists(ctx, resourceName, &conf),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfapigateway.ResourceMethodResponse(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckMethodResponseExists(ctx context.Context, n string, v *apigateway.MethodResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API Gateway Method Response ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn(ctx)

		output, err := tfapigateway.FindMethodResponseByFourPartKey(ctx, conn, rs.Primary.Attributes["http_method"], rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["rest_api_id"], rs.Primary.Attributes["status_code"])

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckMethodResponseDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).APIGatewayConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_api_gateway_method_response" {
				continue
			}

			_, err := tfapigateway.FindMethodResponseByFourPartKey(ctx, conn, rs.Primary.Attributes["http_method"], rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["rest_api_id"], rs.Primary.Attributes["status_code"])

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("API Gateway Method Response %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccMethodResponseImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s/%s/%s", rs.Primary.Attributes["rest_api_id"], rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["http_method"], rs.Primary.Attributes["status_code"]), nil
	}
}

func testAccMethodResponseConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  parent_id   = aws_api_gateway_rest_api.test.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = aws_api_gateway_rest_api.test.id
  resource_id   = aws_api_gateway_resource.test.id
  http_method   = "GET"
  authorization = "NONE"

  request_models = {
    "application/json" = "Error"
  }
}

resource "aws_api_gateway_method_response" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  resource_id = aws_api_gateway_resource.test.id
  http_method = aws_api_gateway_method.test.http_method
  status_code = "400"

  response_models = {
    "application/json" = "Error"
  }

  response_parameters = {
    "method.response.header.Content-Type" = true
  }
}
`, rName)
}

func testAccMethodResponseConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  parent_id   = aws_api_gateway_rest_api.test.root_resource_id
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = aws_api_gateway_rest_api.test.id
  resource_id   = aws_api_gateway_resource.test.id
  http_method   = "GET"
  authorization = "NONE"

  request_models = {
    "application/json" = "Error"
  }
}

resource "aws_api_gateway_method_response" "test" {
  rest_api_id = aws_api_gateway_rest_api.test.id
  resource_id = aws_api_gateway_resource.test.id
  http_method = aws_api_gateway_method.test.http_method
  status_code = "400"

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Host" = false
  }
}
`, rName)
}
