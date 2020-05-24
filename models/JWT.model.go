package models

type Jwt struct {
	Header    string
	Payload   string
	Signature string
}

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JwtPayload struct {
	ID      int  `json:"id"`
	IsAdmin bool `json:"isAdmin"`
	Iat     int  `json:"iat"`
	Exp     int  `json:"exp"`
}
