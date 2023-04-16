package elasticbeanstalk_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccElasticBeanstalkApplicationDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := fmt.Sprintf("tf-acc-test-%s", sdkacctest.RandString(5))
	dataSourceResourceName := "data.aws_elastic_beanstalk_application.test"
	resourceName := "aws_elastic_beanstalk_application.tftest"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceResourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceResourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceResourceName, "description"),
					resource.TestCheckResourceAttr(dataSourceResourceName, "appversion_lifecycle.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "appversion_lifecycle.0.service_role", dataSourceResourceName, "appversion_lifecycle.0.service_role"),
					resource.TestCheckResourceAttrPair(resourceName, "appversion_lifecycle.0.max_age_in_days", dataSourceResourceName, "appversion_lifecycle.0.max_age_in_days"),
					resource.TestCheckResourceAttrPair(resourceName, "appversion_lifecycle.0.delete_source_from_s3", dataSourceResourceName, "appversion_lifecycle.0.delete_source_from_s3"),
				),
			},
		},
	})
}

func testAccApplicationDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "aws_elastic_beanstalk_application" "test" {
  name = aws_elastic_beanstalk_application.tftest.name
}
`, testAccApplicationConfig_maxAge(rName))
}
