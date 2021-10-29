package cns

import "time"

const (
	AppName = "LimeLog"

	SystemFieldPrefix  = "sf_"
	SfTsFieldName      = SystemFieldPrefix + "ts"
	SfTagFieldName     = SystemFieldPrefix + "tag"
	SfMessageFieldName = SystemFieldPrefix + "message"
	MessageFieldName   = "message"

	MaxPageSize = 1000
)

var (
	AppTimeLocation = time.FixedZone("AST", 21600) // +0600
)
