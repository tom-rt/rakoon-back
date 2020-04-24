package types

type Jwt struct {
    Header string
    Payload string
    Signature string
}

type JwtHeader struct {
    Alg string `json:"alg"`
    Typ string `json:"typ"`
}

type JwtPayload struct {
    Name string `json:"name"`
}