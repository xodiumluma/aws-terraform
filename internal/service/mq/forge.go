// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mq

import (
	"github.com/YakDriver/regexache"
	"github.com/beevik/etree"
)

// CanonicalXML reads XML in a string and re-writes it canonically, used for
// comparing XML for logical equivalency
func CanonicalXML(s string) (string, error) {
	doc := etree.NewDocument()
	doc.WriteSettings.CanonicalEndTags = true
	if err := doc.ReadFromString(s); err != nil {
		return "", err
	}

	rawString, err := doc.WriteToString()
	if err != nil {
		return "", err
	}

	re := regexache.MustCompile(`\s`)
	results := re.ReplaceAllString(rawString, "")
	return results, nil
}
