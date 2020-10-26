package vault

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os/user"
	"path"
)

// privateKey holds the multiple key format.
type privateKey struct {
	value  []byte
	rsaKey *rsa.PrivateKey
}

// newPrivateKey creates a RSA Private Key of specified byte size.
func newPrivateKey(bitSize int) (*privateKey, error) {
	// Private Key generation
	rsaKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	//TODO: passphrase

	// Validate Private Key
	err = rsaKey.Validate()
	if err != nil {
		return nil, err
	}

	// Encode the private key
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "WTF VAULT COMMAND KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(rsaKey),
		},
	)

	return &privateKey{
		value:  pemdata,
		rsaKey: rsaKey,
	}, nil
}

// PrivateKeyPath returns the path to the private key.
func PrivateKeyPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(u.HomeDir, ".wtfvault.key"), nil
}

// PrivateKeyDefaultUsername returns the default username when creating the vault.
func PrivateKeyDefaultUsername() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

// Value returns the value.
func (k *privateKey) Value() []byte {
	return k.value
}

// RsaKey parse the string value.
func (k *privateKey) RsaKey() (*rsa.PrivateKey, error) {
	if k.rsaKey != nil {
		return k.rsaKey, nil
	}

	block, _ := pem.Decode(k.value)
	if block == nil {
		return nil, errors.New("Failed to parse file containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	k.rsaKey = priv
	return priv, nil
}

// Decrypt decrypts encoded message.
func (k *privateKey) Decrypt(ciphertext []byte) ([]byte, error) {
	key, err := k.RsaKey()
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, key, ciphertext)
}

// PublicKey returns the public key.
func (k *privateKey) PublicKey() (*publicKey, error) {
	pk, err := k.RsaKey()
	if err != nil {
		return nil, err
	}

	pub := publicKey{}
	pub.SetRsaKey(&pk.PublicKey)
	return &pub, nil
}
