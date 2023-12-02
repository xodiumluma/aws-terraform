// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rds_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	tfrds "github.com/hashicorp/terraform-provider-aws/internal/service/rds"
)

func TestInstanceStateUpgradeV0(t *testing.T) {
	ctx := acctest.Context(t)
	t.Parallel()

	testCases := []struct {
		Description   string
		InputState    map[string]interface{}
		ExpectedState map[string]interface{}
	}{
		{
			Description:   "missing state",
			InputState:    nil,
			ExpectedState: nil,
		},
		{
			Description: "adds delete_automated_backups",
			InputState: map[string]interface{}{
				"allocated_storage": 10,
				"engine":            "mariadb",
				"identifier":        "my-test-instance",
				"instance_class":    "db.t2.micro",
				"password":          "avoid-plaintext-passwords",
				"username":          "tfacctest",
				"tags":              map[string]interface{}{"key1": "value1"},
			},
			ExpectedState: map[string]interface{}{
				"allocated_storage":        10,
				"delete_automated_backups": true,
				"engine":                   "mariadb",
				"identifier":               "my-test-instance",
				"instance_class":           "db.t2.micro",
				"password":                 "avoid-plaintext-passwords",
				"username":                 "tfacctest",
				"tags":                     map[string]interface{}{"key1": "value1"},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			got, err := tfrds.InstanceStateUpgradeV0(ctx, testCase.InputState, nil)
			if err != nil {
				t.Fatalf("error migrating state: %s", err)
			}

			if !reflect.DeepEqual(testCase.ExpectedState, got) {
				t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", testCase.ExpectedState, got)
			}
		})
	}
}

func TestInstanceStateUpgradeV1(t *testing.T) {
	ctx := acctest.Context(t)
	t.Parallel()

	testCases := []struct {
		Description   string
		InputState    map[string]interface{}
		ExpectedState map[string]interface{}
	}{
		{
			Description:   "missing state",
			InputState:    nil,
			ExpectedState: nil,
		},
		{
			Description: "change id to resource id",
			InputState: map[string]interface{}{
				"allocated_storage": 10,
				"engine":            "mariadb",
				"id":                "my-test-instance",
				"identifier":        "my-test-instance",
				"instance_class":    "db.t2.micro",
				"password":          "avoid-plaintext-passwords",
				"resource_id":       "db-cnuap2ilnbmok4eunzklfvwjca",
				"tags":              map[string]interface{}{"key1": "value1"},
				"username":          "tfacctest",
			},
			ExpectedState: map[string]interface{}{
				"allocated_storage": 10,
				"engine":            "mariadb",
				"id":                "db-cnuap2ilnbmok4eunzklfvwjca",
				"identifier":        "my-test-instance",
				"instance_class":    "db.t2.micro",
				"password":          "avoid-plaintext-passwords",
				"resource_id":       "db-cnuap2ilnbmok4eunzklfvwjca",
				"tags":              map[string]interface{}{"key1": "value1"},
				"username":          "tfacctest",
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()

			got, err := tfrds.InstanceStateUpgradeV1(ctx, testCase.InputState, nil)
			if err != nil {
				t.Fatalf("error migrating state: %s", err)
			}

			if !reflect.DeepEqual(testCase.ExpectedState, got) {
				t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", testCase.ExpectedState, got)
			}
		})
	}
}
