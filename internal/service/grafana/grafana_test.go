// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package grafana_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccGrafana_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"Workspace": {
			"saml":                     testAccWorkspace_saml,
			"sso":                      testAccWorkspace_sso,
			"disappears":               testAccWorkspace_disappears,
			"organization":             testAccWorkspace_organization,
			"dataSources":              testAccWorkspace_dataSources,
			"permissionType":           testAccWorkspace_permissionType,
			"notificationDestinations": testAccWorkspace_notificationDestinations,
			"tags":                     testAccWorkspace_tags,
			"vpc":                      testAccWorkspace_vpc,
			"configuration":            testAccWorkspace_configuration,
			"networkAccess":            testAccWorkspace_networkAccess,
			"version":                  testAccWorkspace_version,
		},
		"ApiKey": {
			"basic": testAccWorkspaceAPIKey_basic,
		},
		"DataSource": {
			"basic": testAccWorkspaceDataSource_basic,
		},
		"LicenseAssociation": {
			"enterpriseFreeTrial": testAccLicenseAssociation_freeTrial,
		},
		"SamlConfiguration": {
			"basic":         testAccWorkspaceSAMLConfiguration_basic,
			"loginValidity": testAccWorkspaceSAMLConfiguration_loginValidity,
			"assertions":    testAccWorkspaceSAMLConfiguration_assertions,
		},
		"RoleAssociation": {
			"usersAdmin":           testAccRoleAssociation_usersAdmin,
			"usersEditor":          testAccRoleAssociation_usersEditor,
			"groupsAdmin":          testAccRoleAssociation_groupsAdmin,
			"groupsEditor":         testAccRoleAssociation_groupsEditor,
			"usersAndGroupsAdmin":  testAccRoleAssociation_usersAndGroupsAdmin,
			"usersAndGroupsEditor": testAccRoleAssociation_usersAndGroupsEditor,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
