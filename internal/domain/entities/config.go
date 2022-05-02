package entities

import (
	"time"
)

type ConfigSt struct {
	Rotation ConfigRotationSt `bson:"rotation" json:"rotation"`
}

type ConfigRotationSt struct {
	DefaultDur time.Duration               `bson:"default_dur" json:"default_dur"`
	Exceptions []ConfigRotationExceptionSt `bson:"exceptions" json:"exceptions"`
}

type ConfigRotationExceptionSt struct {
	Tag string        `bson:"tag" json:"tag"`
	Dur time.Duration `bson:"dur" json:"dur"`
}
