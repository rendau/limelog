package cns

import "time"

const (
	AppName = "LimeLog"
	AppUrl  = "https://limeLog.com"

	MaxPageSize = 1000
)

var (
	AppTimeLocation = time.FixedZone("AST", 21600) // +0600
)

// Static file directories
const (
	SFDUsrAva = "usr_avatar"
)

// Usr types
const (
	UsrTypeUndefined = 0
	UsrTypeAdmin     = 1
)

func UsrTypeIsValid(v int) bool {
	return v == UsrTypeUndefined ||
		v == UsrTypeAdmin
}

// Notification types
const (
	NfTypeRefreshProfile = "refresh-profile"
	NfTypeRefreshNumbers = "refresh-numbers"
)
