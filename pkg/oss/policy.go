package oss

import _ "embed"

//go:embed policies/upload.json
var uploadBucketPolicy string

//go:embed policies/private.json
var privateBucketPolicy string

//go:embed policies/readonly.json
var readonlyBucketPolicy string
