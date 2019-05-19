package lib

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"

	jwt "github.com/gbrlsnchs/jwt/v2"
)

type ISigner interface {
	Sign([]byte) ([]byte, error)
}

type ES256Signer struct {
	Key string
}

func (signer ES256Signer) Sign(payload []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(signer.Key))
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	es256 := jwt.NewES256(privateKey, &privateKey.PublicKey)

	if err != nil {
		return []byte{}, err
	}

	header, _ := json.Marshal(map[string]string{
		"alg": "ES256",
		"typ": "JWT",
	})

	headerEnc := base64.StdEncoding.EncodeToString(header)
	payloadEnc := base64.StdEncoding.EncodeToString(payload)
	signed, err := es256.Sign([]byte(headerEnc + "." + payloadEnc))

	if err != nil {
		return []byte{}, err
	}

	return signed, err
}
