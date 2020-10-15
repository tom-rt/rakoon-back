package models

// Jwt struct
type Jwt struct {
	Header    string
	Payload   string
	Signature string
}

// JwtHeader struct
type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// JwtPayload struct
type JwtPayload struct {
	ID      int  `json:"id"`
	IsAdmin bool `json:"isAdmin"`
	Iat     int  `json:"iat"`
	Exp     int  `json:"exp"`
}
