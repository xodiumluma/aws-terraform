// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestValidSecretName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "testing123!",
			ErrCount: 1,
		},
		{
			Value:    "testing 123",
			ErrCount: 1,
		},
		{
			Value:    sdkacctest.RandStringFromCharSet(513, sdkacctest.CharSetAlpha),
			ErrCount: 1,
		},
	}
	for _, tc := range cases {
		_, errors := validSecretName(tc.Value, "aws_secretsmanager_secret")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the AWS Secretsmanager Secret Name to not trigger a validation error for %q", tc.Value)
		}
	}
}

func TestValidSecretNamePrefix(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "testing123!",
			ErrCount: 1,
		},
		{
			Value:    "testing 123",
			ErrCount: 1,
		},
		{
			Value:    sdkacctest.RandStringFromCharSet(512, sdkacctest.CharSetAlpha),
			ErrCount: 1,
		},
	}
	for _, tc := range cases {
		_, errors := validSecretNamePrefix(tc.Value, "aws_secretsmanager_secret")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the AWS Secretsmanager Secret Name to not trigger a validation error for %q", tc.Value)
		}
	}
}
