package lambda_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tflambda "github.com/hashicorp/terraform-provider-aws/internal/service/lambda"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccLambdaProvisionedConcurrencyConfig_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	lambdaFunctionResourceName := "aws_lambda_function.test"
	resourceName := "aws_lambda_provisioned_concurrency_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.LambdaEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedConcurrencyConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedConcurrencyConfigConfig_qualifierFunctionVersion(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisionedConcurrencyExistsConfig(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "function_name", lambdaFunctionResourceName, "function_name"),
					resource.TestCheckResourceAttr(resourceName, "provisioned_concurrent_executions", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "qualifier", lambdaFunctionResourceName, "version"),
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

func TestAccLambdaProvisionedConcurrencyConfig_Disappears_lambdaFunction(t *testing.T) {
	ctx := acctest.Context(t)
	var function lambda.GetFunctionOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	lambdaFunctionResourceName := "aws_lambda_function.test"
	resourceName := "aws_lambda_provisioned_concurrency_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.LambdaEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedConcurrencyConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedConcurrencyConfigConfig_concurrentExecutions(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(ctx, lambdaFunctionResourceName, &function),
					testAccCheckProvisionedConcurrencyExistsConfig(ctx, resourceName),
					testAccCheckFunctionDisappears(ctx, &function),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLambdaProvisionedConcurrencyConfig_Disappears_lambdaProvisionedConcurrency(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lambda_provisioned_concurrency_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.LambdaEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedConcurrencyConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedConcurrencyConfigConfig_concurrentExecutions(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisionedConcurrencyExistsConfig(ctx, resourceName),
					testAccCheckProvisionedConcurrencyDisappearsConfig(ctx, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLambdaProvisionedConcurrencyConfig_provisionedConcurrentExecutions(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lambda_provisioned_concurrency_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.LambdaEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedConcurrencyConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedConcurrencyConfigConfig_concurrentExecutions(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisionedConcurrencyExistsConfig(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "function_name", rName),
					resource.TestCheckResourceAttr(resourceName, "provisioned_concurrent_executions", "1"),
					resource.TestCheckResourceAttr(resourceName, "qualifier", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProvisionedConcurrencyConfigConfig_concurrentExecutions(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisionedConcurrencyExistsConfig(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "function_name", rName),
					resource.TestCheckResourceAttr(resourceName, "provisioned_concurrent_executions", "2"),
					resource.TestCheckResourceAttr(resourceName, "qualifier", "1"),
				),
			},
		},
	})
}

func TestAccLambdaProvisionedConcurrencyConfig_Qualifier_aliasName(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	lambdaAliasResourceName := "aws_lambda_alias.test"
	resourceName := "aws_lambda_provisioned_concurrency_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.LambdaEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedConcurrencyConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedConcurrencyConfigConfig_qualifierAliasName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisionedConcurrencyExistsConfig(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "qualifier", lambdaAliasResourceName, "name"),
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

func testAccCheckProvisionedConcurrencyConfigDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).LambdaClient()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_lambda_provisioned_concurrency_config" {
				continue
			}

			functionName, qualifier, err := tflambda.ProvisionedConcurrencyConfigParseID(rs.Primary.ID)

			if err != nil {
				return err
			}

			input := &lambda.GetProvisionedConcurrencyConfigInput{
				FunctionName: aws.String(functionName),
				Qualifier:    aws.String(qualifier),
			}

			output, err := conn.GetProvisionedConcurrencyConfig(ctx, input)
			if err != nil {
				var pccnfe *types.ProvisionedConcurrencyConfigNotFoundException
				if errors.As(err, &pccnfe) {
					continue
				}
				var nfe *types.ResourceNotFoundException
				if errors.As(err, &nfe) {
					continue
				}
				return err
			}

			if output != nil {
				return fmt.Errorf("Lambda Provisioned Concurrency Config (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckProvisionedConcurrencyDisappearsConfig(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LambdaClient()

		functionName, qualifier, err := tflambda.ProvisionedConcurrencyConfigParseID(rs.Primary.ID)

		if err != nil {
			return err
		}

		input := &lambda.DeleteProvisionedConcurrencyConfigInput{
			FunctionName: aws.String(functionName),
			Qualifier:    aws.String(qualifier),
		}

		_, err = conn.DeleteProvisionedConcurrencyConfig(ctx, input)

		return err
	}
}

func testAccCheckProvisionedConcurrencyExistsConfig(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LambdaClient()

		functionName, qualifier, err := tflambda.ProvisionedConcurrencyConfigParseID(rs.Primary.ID)

		if err != nil {
			return err
		}

		input := &lambda.GetProvisionedConcurrencyConfigInput{
			FunctionName: aws.String(functionName),
			Qualifier:    aws.String(qualifier),
		}

		output, err := conn.GetProvisionedConcurrencyConfig(ctx, input)

		if err != nil {
			return err
		}

		if got, want := output.Status, types.ProvisionedConcurrencyStatusEnumReady; got != want {
			return fmt.Errorf("Lambda Provisioned Concurrency Config (%s) expected status (%s), got: %s", rs.Primary.ID, want, got)
		}

		return nil
	}
}

func testAccProvisionedConcurrencyConfig_base(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "test" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.test.id
}

resource "aws_lambda_function" "test" {
  filename      = "test-fixtures/lambdapinpoint.zip"
  function_name = %[1]q
  role          = aws_iam_role.test.arn
  handler       = "lambdapinpoint.handler"
  publish       = true
  runtime       = "nodejs16.x"

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName)
}

func testAccProvisionedConcurrencyConfigConfig_concurrentExecutions(rName string, provisionedConcurrentExecutions int) string {
	return testAccProvisionedConcurrencyConfig_base(rName) + fmt.Sprintf(`
resource "aws_lambda_provisioned_concurrency_config" "test" {
  function_name                     = aws_lambda_function.test.function_name
  provisioned_concurrent_executions = %[1]d
  qualifier                         = aws_lambda_function.test.version
}
`, provisionedConcurrentExecutions)
}

func testAccProvisionedConcurrencyConfigConfig_qualifierAliasName(rName string) string {
	return testAccProvisionedConcurrencyConfig_base(rName) + `
resource "aws_lambda_alias" "test" {
  function_name    = aws_lambda_function.test.function_name
  function_version = aws_lambda_function.test.version
  name             = "test"
}

resource "aws_lambda_provisioned_concurrency_config" "test" {
  function_name                     = aws_lambda_alias.test.function_name
  provisioned_concurrent_executions = 1
  qualifier                         = aws_lambda_alias.test.name
}
`
}

func testAccProvisionedConcurrencyConfigConfig_qualifierFunctionVersion(rName string) string {
	return testAccProvisionedConcurrencyConfig_base(rName) + `
resource "aws_lambda_provisioned_concurrency_config" "test" {
  function_name                     = aws_lambda_function.test.function_name
  provisioned_concurrent_executions = 1
  qualifier                         = aws_lambda_function.test.version
}
`
}
