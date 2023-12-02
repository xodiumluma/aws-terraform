// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package emr

import (
	"fmt"
	"strings"

	"github.com/YakDriver/regexache"
)

const IdentityIdPattern = `([0-9a-f]{10}-|)[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}`

var IdentityIdPatternRegexp = regexache.MustCompile(IdentityIdPattern)

func isIdentityId(identityIdOrName string) bool {
	return IdentityIdPatternRegexp.MatchString(identityIdOrName)
}

func readStudioSessionMapping(id string) (studioId, identityType, identityIdOrName string, err error) {
	idOrNameParts := strings.Split(id, ":")
	if len(idOrNameParts) == 3 {
		return idOrNameParts[0], idOrNameParts[1], idOrNameParts[2], nil
	}

	if isIdentityId(identityIdOrName) {
		err = fmt.Errorf("expected ID in format studio-id:identity-type:identity-id, received: %s", identityIdOrName)
	} else {
		err = fmt.Errorf("expected ID in format studio-id:identity-type:identity-name, received: %s", identityIdOrName)
	}

	return "", "", "", err
}
