// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inspector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/inspector"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfinspector "github.com/hashicorp/terraform-provider-aws/internal/service/inspector"
)

func TestAccInspectorAssessmentTarget_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var assessmentTarget1 inspector.AssessmentTarget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_inspector_assessment_target.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, inspector.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetAssessmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAssessmentTargetConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetExists(ctx, resourceName, &assessmentTarget1),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "inspector", regexache.MustCompile(`target/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "resource_group_arn", ""),
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

func TestAccInspectorAssessmentTarget_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var assessmentTarget1 inspector.AssessmentTarget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_inspector_assessment_target.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, inspector.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetAssessmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAssessmentTargetConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetExists(ctx, resourceName, &assessmentTarget1),
					testAccCheckTargetDisappears(ctx, &assessmentTarget1),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccInspectorAssessmentTarget_name(t *testing.T) {
	ctx := acctest.Context(t)
	var assessmentTarget1, assessmentTarget2 inspector.AssessmentTarget
	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_inspector_assessment_target.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, inspector.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetAssessmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAssessmentTargetConfig_basic(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetExists(ctx, resourceName, &assessmentTarget1),
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAssessmentTargetConfig_basic(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetExists(ctx, resourceName, &assessmentTarget2),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccInspectorAssessmentTarget_resourceGroupARN(t *testing.T) {
	ctx := acctest.Context(t)
	var assessmentTarget1, assessmentTarget2, assessmentTarget3, assessmentTarget4 inspector.AssessmentTarget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	inspectorResourceGroupResourceName1 := "aws_inspector_resource_group.test1"
	inspectorResourceGroupResourceName2 := "aws_inspector_resource_group.test2"
	resourceName := "aws_inspector_assessment_target.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, inspector.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetAssessmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAssessmentTargetConfig_resourceGroupARN(rName, inspectorResourceGroupResourceName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetExists(ctx, resourceName, &assessmentTarget1),
					resource.TestCheckResourceAttrPair(resourceName, "resource_group_arn", inspectorResourceGroupResourceName1, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAssessmentTargetConfig_resourceGroupARN(rName, inspectorResourceGroupResourceName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetExists(ctx, resourceName, &assessmentTarget2),
					resource.TestCheckResourceAttrPair(resourceName, "resource_group_arn", inspectorResourceGroupResourceName2, "arn"),
				),
			},
			{
				Config: testAccAssessmentTargetConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetExists(ctx, resourceName, &assessmentTarget3),
					resource.TestCheckResourceAttr(resourceName, "resource_group_arn", ""),
				),
			},
			{
				Config: testAccAssessmentTargetConfig_resourceGroupARN(rName, inspectorResourceGroupResourceName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetExists(ctx, resourceName, &assessmentTarget4),
					resource.TestCheckResourceAttrPair(resourceName, "resource_group_arn", inspectorResourceGroupResourceName1, "arn"),
				),
			},
		},
	})
}

func testAccCheckTargetAssessmentDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).InspectorConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_inspector_assessment_target" {
				continue
			}

			assessmentTarget, err := tfinspector.DescribeAssessmentTarget(ctx, conn, rs.Primary.ID)

			if err != nil {
				return fmt.Errorf("finding Inspector Classic Assessment Target: %s", err)
			}

			if assessmentTarget != nil {
				return fmt.Errorf("Inspector Classic Assessment Target (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckTargetExists(ctx context.Context, name string, target *inspector.AssessmentTarget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).InspectorConn(ctx)

		assessmentTarget, err := tfinspector.DescribeAssessmentTarget(ctx, conn, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("finding Inspector Classic Assessment Target: %s", err)
		}

		if assessmentTarget == nil {
			return fmt.Errorf("Inspector Classic Assessment Target (%s) not found", rs.Primary.ID)
		}

		*target = *assessmentTarget

		return nil
	}
}

func testAccCheckTargetDisappears(ctx context.Context, assessmentTarget *inspector.AssessmentTarget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).InspectorConn(ctx)

		input := &inspector.DeleteAssessmentTargetInput{
			AssessmentTargetArn: assessmentTarget.Arn,
		}

		_, err := conn.DeleteAssessmentTargetWithContext(ctx, input)

		return err
	}
}

func testAccAssessmentTargetConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_inspector_assessment_target" "test" {
  name = %q
}
`, rName)
}

func testAccAssessmentTargetConfig_resourceGroupARN(rName, inspectorResourceGroupResourceName string) string {
	return fmt.Sprintf(`
resource "aws_inspector_resource_group" "test1" {
  tags = {
    Name = "%s1"
  }
}

resource "aws_inspector_resource_group" "test2" {
  tags = {
    Name = "%s2"
  }
}

resource "aws_inspector_assessment_target" "test" {
  name               = %q
  resource_group_arn = %s.arn
}
`, rName, rName, rName, inspectorResourceGroupResourceName)
}
