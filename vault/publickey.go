package vault

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
)

// publicKey holds the multiple key format.
type publicKey struct {
	// Cached values
	value       string
	rsaKey      *rsa.PublicKey
	fingerprint string
}

// Value returns public key's base64 string.
func (k *publicKey) Value() string {
	return k.value
}

// Fingerprint returns public key's fingerprint.
func (k *publicKey) Fingerprint() string {
	if k.fingerprint != "" {
		return k.fingerprint
	}

	h := md5.New()
	io.WriteString(h, k.Value())
	k.fingerprint = fmt.Sprintf("%x", h.Sum(nil))

	return k.fingerprint
}

// RsaKey decodes the base64 Value version.
func (k *publicKey) RsaKey() (*rsa.PublicKey, error) {
	if k.rsaKey != nil {
		return k.rsaKey, nil
	}

	// Decode base64
	keyBytes, err := base64.StdEncoding.DecodeString(k.Value())
	if err != nil {
		return nil, err
	}

	// Check length
	keyLen := len(keyBytes)
	if keyLen < 6 {
		return nil, errors.New("Public key too small")
	}

	// e length
	eLen := int(keyBytes[0])<<24 | int(keyBytes[1])<<16 | int(keyBytes[2])<<8 | int(keyBytes[3])
	if eLen > 24 {
		return nil, errors.New("Invalid E length")
	}

	// Get big ints
	eBigInt := new(big.Int).SetBytes(keyBytes[4 : 4+eLen])
	nBigInt := new(big.Int).SetBytes(keyBytes[4+eLen:])

	// E
	e := eBigInt.Int64()
	if e < 3 || e&1 == 0 {
		return nil, errors.New("Invalid E")
	}

	// To object
	res := rsa.PublicKey{}
	res.E = int(e)
	res.N = nBigInt
	k.rsaKey = &res
	return &res, nil
}

// SetRsaKey changes the RSA key.
func (k *publicKey) SetRsaKey(key *rsa.PublicKey) {
	// Convert big ints to byte array
	e := new(big.Int).SetInt64(int64(key.E))

	eBytes := e.Bytes()
	nBytes := key.N.Bytes()
	eLen := len(eBytes)

	out := make([]byte, 0, 4+eLen+len(nBytes))
	out = append(out, byte(eLen>>24&0xFF), byte(eLen>>16&0xFF), byte(eLen>>8&0xFF), byte(eLen&0xFF))
	out = append(out, eBytes...)
	out = append(out, nBytes...)

	// Encode base64
	k.value = base64.StdEncoding.EncodeToString(out)

	// Reset
	k.fingerprint = ""
	k.rsaKey = key
}

// Encrypt encrypts with the public key.
func (k *publicKey) Encrypt(plaintext []byte) ([]byte, error) {
	key, err := k.RsaKey()
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, key, plaintext)
}

// MarshalJSON encodes as json.
func (k *publicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Value())
}

// UnmarshalJSON decodes a json.
func (k *publicKey) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &k.value)
}
