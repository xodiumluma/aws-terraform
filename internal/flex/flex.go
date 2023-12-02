// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package flex

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfmaps "github.com/hashicorp/terraform-provider-aws/internal/maps"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
)

const (
	// A common separator to be used for creating resource Ids from a combination of attributes
	ResourceIdSeparator = ","
)

// ExpandStringList the result of flatmap.Expand for an array of strings
// and returns a []*string. Empty strings are skipped.
func ExpandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		if v, ok := v.(string); ok && v != "" { // v != "" may not do anything since in []interface{}, empty string will be nil so !ok
			vs = append(vs, aws.String(v))
		}
	}
	return vs
}

// ExpandStringListEmpty the result of flatmap. Expand for an array of strings
// and returns a []*string. Adds an empty element for every nil or uncastable.
func ExpandStringListEmpty(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		if v, ok := v.(string); ok { // empty string in config turns into nil in []interface{} so !ok
			vs = append(vs, aws.String(v))
		} else {
			vs = append(vs, aws.String(""))
		}
	}
	return vs
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*time.Time
func ExpandStringTimeList(configured []interface{}, format string) []*time.Time {
	vs := make([]*time.Time, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			t, _ := time.Parse(format, v.(string))
			vs = append(vs, aws.Time(t))
		}
	}
	return vs
}

// ExpandStringValueList takes the result of flatmap.Expand for an array of strings
// and returns a []string
func ExpandStringValueList(configured []interface{}) []string {
	return ExpandStringyValueList[string](configured)
}

func ExpandStringyValueList[E ~string](configured []any) []E {
	vs := make([]E, 0, len(configured))
	for _, v := range configured {
		if val, ok := v.(string); ok && val != "" {
			vs = append(vs, E(val))
		}
	}
	return vs
}

// Takes list of pointers to strings. Expand to an array
// of raw strings and returns a []interface{}
// to keep compatibility w/ schema.NewSetschema.NewSet
func FlattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

// Takes list of pointers to time.Time. Expand to an array
// of strings and returns a []interface{}
func FlattenTimeStringList(list []*time.Time, format string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v.Format(format))
	}
	return vs
}

// Takes list of strings. Expand to an array
// of raw strings and returns a []interface{}
// to keep compatibility w/ schema.NewSetschema.NewSet
func FlattenStringValueList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

// Expands a map of string to interface to a map of string to int32
func ExpandInt32Map(m map[string]interface{}) map[string]int32 {
	return tfmaps.ApplyToAllValues(m, func(v any) int32 {
		return int32(v.(int))
	})
}

// Expands a map of string to interface to a map of string to *int64
func ExpandInt64Map(m map[string]interface{}) map[string]*int64 {
	return tfmaps.ApplyToAllValues(m, func(v any) *int64 {
		return aws.Int64(int64(v.(int)))
	})
}

// ExpandInt64ValueMap expands a map of string to interface to a map of string to int64
func ExpandInt64ValueMap(m map[string]interface{}) map[string]int64 {
	return tfmaps.ApplyToAllValues(m, func(v any) int64 {
		return int64(v.(int))
	})
}

// Expands a map of string to interface to a map of string to *string
func ExpandStringMap(m map[string]interface{}) map[string]*string {
	return tfmaps.ApplyToAllValues(m, func(v any) *string {
		return aws.String(v.(string))
	})
}

func ExpandStringyValueMap[M ~map[K]V, K ~string, V ~string](m M) map[string]string {
	return tfmaps.ApplyToAllKeys(tfmaps.ApplyToAllValues(m,
		func(v V) string {
			return string(v)
		}),
		func(k K) string {
			return string(k)
		})
}

// ExpandStringValueMap expands a string map of interfaces to a string map of strings
func ExpandStringValueMap(m map[string]interface{}) map[string]string {
	return tfmaps.ApplyToAllValues(m, func(v any) string {
		return v.(string)
	})
}

// Expands a map of string to interface to a map of string to *bool
func ExpandBoolMap(m map[string]interface{}) map[string]*bool {
	return tfmaps.ApplyToAllValues(m, func(v any) *bool {
		return aws.Bool(v.(bool))
	})
}

// Takes the result of schema.Set of strings and returns a []*string
func ExpandStringSet(configured *schema.Set) []*string {
	return ExpandStringList(configured.List()) // nosemgrep:ci.helper-schema-Set-extraneous-ExpandStringList-with-List
}

func ExpandStringValueSet(configured *schema.Set) []string {
	return ExpandStringValueList(configured.List()) // nosemgrep:ci.helper-schema-Set-extraneous-ExpandStringList-with-List
}

func ExpandStringyValueSet[E ~string](configured *schema.Set) []E {
	return ExpandStringyValueList[E](configured.List())
}

func FlattenStringSet(list []*string) *schema.Set {
	return schema.NewSet(schema.HashString, FlattenStringList(list)) // nosemgrep:ci.helper-schema-Set-extraneous-NewSet-with-FlattenStringList
}

func FlattenStringValueSet(list []string) *schema.Set {
	return schema.NewSet(schema.HashString, FlattenStringValueList(list)) // nosemgrep: helper-schema-Set-extraneous-NewSet-with-FlattenStringList
}

func FlattenStringMap(m map[string]*string) map[string]interface{} {
	return tfmaps.ApplyToAllValues(m, func(v *string) any {
		return aws.StringValue(v)
	})
}

// Takes the result of schema.Set of strings and returns a []*int64
func ExpandInt64Set(configured *schema.Set) []*int64 {
	return ExpandInt64List(configured.List())
}

func FlattenInt64Set(list []*int64) *schema.Set {
	return schema.NewSet(schema.HashInt, FlattenInt64List(list))
}

// Takes the result of flatmap.Expand for an array of int64
// and returns a []*int64
func ExpandInt64List(configured []interface{}) []*int64 {
	return tfslices.ApplyToAll(configured, func(v any) *int64 {
		return aws.Int64(int64(v.(int)))
	})
}

// Takes the result of flatmap.Expand for an array of float64
// and returns a []*float64
func ExpandFloat64List(configured []interface{}) []*float64 {
	return tfslices.ApplyToAll(configured, func(v any) *float64 {
		return aws.Float64(v.(float64))
	})
}

// Takes list of pointers to int64s. Expand to an array
// of raw ints and returns a []interface{}
// to keep compatibility w/ schema.NewSet
func FlattenInt64List(list []*int64) []interface{} {
	return tfslices.ApplyToAll(list, func(v *int64) any {
		return int(aws.Int64Value(v))
	})
}

// Takes list of pointers to float64s. Expand to an array
// of raw floats and returns a []interface{}
// to keep compatibility w/ schema.NewSet
func FlattenFloat64List(list []*float64) []interface{} {
	return tfslices.ApplyToAll(list, func(v *float64) any {
		return int(aws.Float64Value(v))
	})
}

// Takes a string of resource attributes separated by the ResourceIdSeparator constant, an expected number of Id Parts, and a boolean specifying if empty parts are to be allowed
// Returns a list of the resource attributes strings used to construct the unique Id or an error message if the resource id does not parse properly
func ExpandResourceId(id string, partCount int, allowEmptyPart bool) ([]string, error) {
	idParts := strings.Split(id, ResourceIdSeparator)

	if len(idParts) <= 1 {
		return nil, fmt.Errorf("unexpected format for ID (%v), expected more than one part", idParts)
	}

	if len(idParts) != partCount {
		return nil, fmt.Errorf("unexpected format for ID (%s), expected (%d) parts separated by (%s)", id, partCount, ResourceIdSeparator)
	}

	if !allowEmptyPart {
		var emptyPart bool
		emptyParts := make([]int, 0, partCount)
		for index, part := range idParts {
			if part == "" {
				emptyPart = true
				emptyParts = append(emptyParts, index)
			}
		}

		if emptyPart {
			return nil, fmt.Errorf("unexpected format for ID (%[1]s), the following id parts indexes are blank (%v)", id, emptyParts)
		}
	}
	return idParts, nil
}

// Takes a list of the resource attributes as strings used to construct the unique Id, an expected number of Id Parts, and a boolean specifying if empty parts are to be allowed
// Returns a string of resource attributes separated by the ResourceIdSeparator constant or an error message if the id parts do not parse properly
func FlattenResourceId(idParts []string, partCount int, allowEmptyPart bool) (string, error) {
	if len(idParts) <= 1 {
		return "", fmt.Errorf("unexpected format for ID parts (%v), expected more than one part", idParts)
	}

	if len(idParts) != partCount {
		return "", fmt.Errorf("unexpected format for ID parts (%v), expected (%d) parts", idParts, partCount)
	}

	if !allowEmptyPart {
		var emptyPart bool
		emptyParts := make([]int, 0, len(idParts))
		for index, part := range idParts {
			if part == "" {
				emptyPart = true
				emptyParts = append(emptyParts, index)
			}
		}

		if emptyPart {
			return "", fmt.Errorf("unexpected format for ID parts (%v), the following id parts indexes are blank (%v)", idParts, emptyParts)
		}
	}

	return strings.Join(idParts, ResourceIdSeparator), nil
}

// StringToBoolValue converts a string pointer to a Go bool value.
// Only the string "true" is converted to true, all other values return false.
func StringToBoolValue(v *string) bool {
	return aws.StringValue(v) == strconv.FormatBool(true)
}

// Takes a string of resource attributes separated by the ResourceIdSeparator constant
// returns the number of parts
func ResourceIdPartCount(id string) int {
	idParts := strings.Split(id, ResourceIdSeparator)
	return len(idParts)
}

type Set[T comparable] []T

// Difference find the elements in two sets that are not similar.
func (s Set[T]) Difference(ns Set[T]) Set[T] {
	m := make(map[T]struct{})
	for _, v := range ns {
		m[v] = struct{}{}
	}

	var result []T
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// DiffStringMaps returns the set of keys and values that must be created, the set of keys
// and values that must be destroyed, and the set of keys and values that are unchanged.
func DiffStringMaps(oldMap, newMap map[string]interface{}) (map[string]*string, map[string]*string, map[string]*string) {
	// First, we're creating everything we have.
	add := ExpandStringMap(newMap)

	// Build the maps of what to remove and what is unchanged.
	remove := make(map[string]*string)
	unchanged := make(map[string]*string)
	for k, v := range oldMap {
		v := v.(string)
		if old, ok := add[k]; !ok || aws.StringValue(old) != v {
			// Delete it!
			remove[k] = aws.String(v)
		} else if ok {
			unchanged[k] = aws.String(v)
			// Already present, so remove from new.
			delete(add, k)
		}
	}

	return add, remove, unchanged
}
