// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package comprehend

import (
	"github.com/YakDriver/regexache"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	modelIdentifierMaxLen       = 63 // Documentation says 256, Console says 63
	modelIdentifierPrefixMaxLen = modelIdentifierMaxLen - id.UniqueIDSuffixLength
)

var validModelName = validIdentifier
var validModelVersionName = validation.Any( // nosemgrep:ci.avoid-string-is-empty-validation
	validation.StringIsEmpty,
	validIdentifier,
)
var validModelVersionNamePrefix = validIdentifierPrefix

var validIdentifier = validation.All(
	validation.StringLenBetween(1, modelIdentifierMaxLen),
	validIdentifierPattern,
)

var validIdentifierPrefix = validation.All(
	validation.StringLenBetween(1, modelIdentifierPrefixMaxLen),
	validIdentifierPattern,
)

var validIdentifierPattern = validation.StringMatch(regexache.MustCompile(`^[[:alnum:]-]+$`), "must contain A-Z, a-z, 0-9, and hypen (-)")
