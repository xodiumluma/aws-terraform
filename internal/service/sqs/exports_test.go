// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sqs

// Exports for use in tests only.
var (
	ResourceQueue                   = resourceQueue
	ResourceQueuePolicy             = resourceQueuePolicy
	ResourceQueueRedriveAllowPolicy = resourceQueueRedriveAllowPolicy
	ResourceQueueRedrivePolicy      = resourceQueueRedrivePolicy

	FindQueueAttributesByURL = findQueueAttributesByURL

	DefaultQueueDelaySeconds                  = defaultQueueDelaySeconds
	DefaultQueueKMSDataKeyReusePeriodSeconds  = defaultQueueKMSDataKeyReusePeriodSeconds
	DefaultQueueMaximumMessageSize            = defaultQueueMaximumMessageSize
	DefaultQueueMessageRetentionPeriod        = defaultQueueMessageRetentionPeriod
	DefaultQueueReceiveMessageWaitTimeSeconds = defaultQueueReceiveMessageWaitTimeSeconds
	DefaultQueueVisibilityTimeout             = defaultQueueVisibilityTimeout
	FIFOQueueNameSuffix                       = fifoQueueNameSuffix
	QueueDeletedTimeout                       = queueDeletedTimeout
	QueueNameFromURL                          = queueNameFromURL
)
