// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../generate/tags/main.go -ServiceTagsSlice
//go:generate go run ../../generate/tags/main.go -AWSSDKVersion=2 -SkipAWSServiceImp -ServiceTagsSlice -TagsFunc=tagsV2 -KeyValueTagsFunc=keyValueTagsV2 -GetTagsInFunc=getTagsInV2 -SetTagsOutFunc=setTagsOutV2 -- tagsv2_gen.go
//go:generate go run ../../generate/servicepackage/main.go
// ONLY generate directives and package declaration! Do not add anything else to this file.

package s3
