package entities

type ConfigSt struct {
	Rotation ConfigRotationSt `bson:"rotation" json:"rotation"`
}

type ConfigRotationSt struct {
	DefaultDur int64                       `bson:"default_dur" json:"default_dur"`
	Exceptions []ConfigRotationExceptionSt `bson:"exceptions" json:"exceptions"`
}

type ConfigRotationExceptionSt struct {
	Tag string `bson:"tag" json:"tag"`
	Dur int64  `bson:"dur" json:"dur"`
}
