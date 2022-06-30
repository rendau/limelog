package entities

type ProfileAuthReqSt struct {
	Password string `json:"password"`
}

type ProfileAuthRepSt struct {
	Token string `json:"token"`
}

type ProfileSt struct {
}
