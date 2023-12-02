// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chime

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	voiceConnectorResourcePropagationTimeout = 1 * time.Minute
)

func FindVoiceConnectorResourceWithRetry[T any](ctx context.Context, isNewResource bool, f func() (T, error)) (T, error) {
	var resp T
	err := tfresource.Retry(ctx, voiceConnectorResourcePropagationTimeout, func() *retry.RetryError {
		var err error
		resp, err = f()
		if isNewResource && tfresource.NotFound(err) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}, tfresource.WithDelay(5*time.Second))

	return resp, err
}
