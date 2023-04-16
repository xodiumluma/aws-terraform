package sesv2_test

import (
	"fmt"
	"regexp"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSESV2DedicatedIPPoolDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_sesv2_dedicated_ip_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccPreCheckDedicatedIPPool(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDedicatedIPPoolDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedIPPoolDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedIPPoolExists(ctx, dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "pool_name", rName),
					acctest.MatchResourceAttrRegionalARN(dataSourceName, "arn", "ses", regexp.MustCompile(`dedicated-ip-pool/.+`)),
				),
			},
		},
	})
}

func testAccDedicatedIPPoolDataSourceConfig_basic(poolName string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_dedicated_ip_pool" "test" {
  pool_name = %[1]q
}

data "aws_sesv2_dedicated_ip_pool" "test" {
  depends_on = [aws_sesv2_dedicated_ip_pool.test]
  pool_name  = %[1]q
}
`, poolName, poolName)
}
