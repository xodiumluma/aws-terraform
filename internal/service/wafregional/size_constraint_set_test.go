// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wafregional_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfwafregional "github.com/hashicorp/terraform-provider-aws/internal/service/wafregional"
)

func TestAccWAFRegionalSizeConstraintSet_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var constraints waf.SizeConstraintSet
	sizeConstraintSet := fmt.Sprintf("sizeConstraintSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_size_constraint_set.size_constraint_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSizeConstraintSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSizeConstraintSetConfig_basic(sizeConstraintSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSizeConstraintSetExists(ctx, resourceName, &constraints),
					resource.TestCheckResourceAttr(
						resourceName, "name", sizeConstraintSet),
					resource.TestCheckResourceAttr(
						resourceName, "size_constraints.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "size_constraints.*", map[string]string{
						"comparison_operator": "EQ",
						"field_to_match.#":    "1",
						"size":                "4096",
						"text_transformation": "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "size_constraints.*.field_to_match.*", map[string]string{
						"data": "",
						"type": "BODY",
					}),
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

func TestAccWAFRegionalSizeConstraintSet_changeNameForceNew(t *testing.T) {
	ctx := acctest.Context(t)
	var before, after waf.SizeConstraintSet
	sizeConstraintSet := fmt.Sprintf("sizeConstraintSet-%s", sdkacctest.RandString(5))
	sizeConstraintSetNewName := fmt.Sprintf("sizeConstraintSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_size_constraint_set.size_constraint_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSizeConstraintSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSizeConstraintSetConfig_basic(sizeConstraintSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSizeConstraintSetExists(ctx, resourceName, &before),
					resource.TestCheckResourceAttr(
						resourceName, "name", sizeConstraintSet),
					resource.TestCheckResourceAttr(
						resourceName, "size_constraints.#", "1"),
				),
			},
			{
				Config: testAccSizeConstraintSetConfig_changeName(sizeConstraintSetNewName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSizeConstraintSetExists(ctx, resourceName, &after),
					resource.TestCheckResourceAttr(
						resourceName, "name", sizeConstraintSetNewName),
					resource.TestCheckResourceAttr(
						resourceName, "size_constraints.#", "1"),
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

func TestAccWAFRegionalSizeConstraintSet_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var constraints waf.SizeConstraintSet
	sizeConstraintSet := fmt.Sprintf("sizeConstraintSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_size_constraint_set.size_constraint_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSizeConstraintSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSizeConstraintSetConfig_basic(sizeConstraintSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSizeConstraintSetExists(ctx, resourceName, &constraints),
					testAccCheckSizeConstraintSetDisappears(ctx, &constraints),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWAFRegionalSizeConstraintSet_changeConstraints(t *testing.T) {
	ctx := acctest.Context(t)
	var before, after waf.SizeConstraintSet
	setName := fmt.Sprintf("sizeConstraintSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_size_constraint_set.size_constraint_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSizeConstraintSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSizeConstraintSetConfig_basic(setName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSizeConstraintSetExists(ctx, resourceName, &before),
					resource.TestCheckResourceAttr(
						resourceName, "name", setName),
					resource.TestCheckResourceAttr(
						resourceName, "size_constraints.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "size_constraints.*", map[string]string{
						"comparison_operator": "EQ",
						"field_to_match.#":    "1",
						"size":                "4096",
						"text_transformation": "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "size_constraints.*.field_to_match.*", map[string]string{
						"data": "",
						"type": "BODY",
					}),
				),
			},
			{
				Config: testAccSizeConstraintSetConfig_changes(setName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSizeConstraintSetExists(ctx, resourceName, &after),
					resource.TestCheckResourceAttr(
						resourceName, "name", setName),
					resource.TestCheckResourceAttr(
						resourceName, "size_constraints.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "size_constraints.*", map[string]string{
						"comparison_operator": "GE",
						"field_to_match.#":    "1",
						"size":                "1024",
						"text_transformation": "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "size_constraints.*.field_to_match.*", map[string]string{
						"data": "",
						"type": "BODY",
					}),
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

func TestAccWAFRegionalSizeConstraintSet_noConstraints(t *testing.T) {
	ctx := acctest.Context(t)
	var constraints waf.SizeConstraintSet
	setName := fmt.Sprintf("sizeConstraintSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_size_constraint_set.size_constraint_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSizeConstraintSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSizeConstraintSetConfig_nos(setName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSizeConstraintSetExists(ctx, resourceName, &constraints),
					resource.TestCheckResourceAttr(
						resourceName, "name", setName),
					resource.TestCheckResourceAttr(
						resourceName, "size_constraints.#", "0"),
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

func testAccCheckSizeConstraintSetDisappears(ctx context.Context, constraints *waf.SizeConstraintSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn(ctx)
		region := acctest.Provider.Meta().(*conns.AWSClient).Region

		wr := tfwafregional.NewRetryer(conn, region)
		_, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
			req := &waf.UpdateSizeConstraintSetInput{
				ChangeToken:         token,
				SizeConstraintSetId: constraints.SizeConstraintSetId,
			}

			for _, sizeConstraint := range constraints.SizeConstraints {
				sizeConstraintUpdate := &waf.SizeConstraintSetUpdate{
					Action: aws.String("DELETE"),
					SizeConstraint: &waf.SizeConstraint{
						FieldToMatch:       sizeConstraint.FieldToMatch,
						ComparisonOperator: sizeConstraint.ComparisonOperator,
						Size:               sizeConstraint.Size,
						TextTransformation: sizeConstraint.TextTransformation,
					},
				}
				req.Updates = append(req.Updates, sizeConstraintUpdate)
			}
			return conn.UpdateSizeConstraintSetWithContext(ctx, req)
		})
		if err != nil {
			return fmt.Errorf("Error updating SizeConstraintSet: %s", err)
		}

		_, err = wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
			opts := &waf.DeleteSizeConstraintSetInput{
				ChangeToken:         token,
				SizeConstraintSetId: constraints.SizeConstraintSetId,
			}
			return conn.DeleteSizeConstraintSetWithContext(ctx, opts)
		})

		return err
	}
}

func testAccCheckSizeConstraintSetExists(ctx context.Context, n string, constraints *waf.SizeConstraintSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAF SizeConstraintSet ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn(ctx)
		resp, err := conn.GetSizeConstraintSetWithContext(ctx, &waf.GetSizeConstraintSetInput{
			SizeConstraintSetId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if *resp.SizeConstraintSet.SizeConstraintSetId == rs.Primary.ID {
			*constraints = *resp.SizeConstraintSet
			return nil
		}

		return fmt.Errorf("WAF SizeConstraintSet (%s) not found", rs.Primary.ID)
	}
}

func testAccCheckSizeConstraintSetDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_wafregional_size_contraint_set" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn(ctx)
			resp, err := conn.GetSizeConstraintSetWithContext(ctx, &waf.GetSizeConstraintSetInput{
				SizeConstraintSetId: aws.String(rs.Primary.ID),
			})

			if err == nil {
				if *resp.SizeConstraintSet.SizeConstraintSetId == rs.Primary.ID {
					return fmt.Errorf("WAF SizeConstraintSet %s still exists", rs.Primary.ID)
				}
			}

			// Return nil if the SizeConstraintSet is already destroyed
			if tfawserr.ErrCodeEquals(err, wafregional.ErrCodeWAFNonexistentItemException) {
				return nil
			}

			return err
		}

		return nil
	}
}

func testAccSizeConstraintSetConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_size_constraint_set" "size_constraint_set" {
  name = "%s"

  size_constraints {
    text_transformation = "NONE"
    comparison_operator = "EQ"
    size                = "4096"

    field_to_match {
      type = "BODY"
    }
  }
}
`, name)
}

func testAccSizeConstraintSetConfig_changeName(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_size_constraint_set" "size_constraint_set" {
  name = "%s"

  size_constraints {
    text_transformation = "NONE"
    comparison_operator = "EQ"
    size                = "4096"

    field_to_match {
      type = "BODY"
    }
  }
}
`, name)
}

func testAccSizeConstraintSetConfig_changes(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_size_constraint_set" "size_constraint_set" {
  name = "%s"

  size_constraints {
    text_transformation = "NONE"
    comparison_operator = "GE"
    size                = "1024"

    field_to_match {
      type = "BODY"
    }
  }
}
`, name)
}

func testAccSizeConstraintSetConfig_nos(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_size_constraint_set" "size_constraint_set" {
  name = "%s"
}
`, name)
}
