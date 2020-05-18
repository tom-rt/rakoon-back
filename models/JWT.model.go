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
	Name string `json:"name"`
	Iat  int    `json:"iat"`
	Exp  int    `json:"exp"`
}

type JwtInput struct {
	Name    string `json:"name"`
	IsAdmin *bool  `json:"isAdmin"`
}
