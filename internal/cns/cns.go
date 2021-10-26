package cns

import "time"

const (
	AppName = "LimeLog"

	SystemFieldPrefix  = "sf_"
	SfTsFieldName      = SystemFieldPrefix + "ts"
	SfMessageFieldName = SystemFieldPrefix + "message"
	MessageFieldName   = "message"

	MaxPageSize = 1000
)

var (
	AppTimeLocation = time.FixedZone("AST", 21600) // +0600
)
