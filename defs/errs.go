package defs

import "errors"

var (
	ERR_NO_ENOUGH_TOKEN_BUCKET = errors.New("No enough token bucket")
	ERR_GET_TOKEN_BUCKET_FAIL  = errors.New("Get token bucket fail")
)
