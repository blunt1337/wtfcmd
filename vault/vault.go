package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

type recipient struct {
	PublicKey *publicKey `json:"key"`
	SymKey    []byte     `json:"sym"`
	Name      string     `json:"name"`
}
type vault struct {
	PrivateKey       *privateKey       `json:"-"`
	Recipients       []*recipient      `json:"users"`
	Passwords        map[string]string `json:"-"`
	EncodedPasswords []byte            `json:"content"`
}

// InitVault generates private key if missing, and .wtfvault if missing.
func newVault(init bool, load bool) (res *vault, err error) {
	var pk *privateKey

	// Check if private key exists.
	var pkPath string
	pkPath, err = PrivateKeyPath()
	if err != nil {
		return
	}

	if fileExists(pkPath) {
		// Load private key
		var value []byte
		value, err = ioutil.ReadFile(pkPath)
		if err != nil {
			return
		}
		pk = &privateKey{value: value}
	} else if init {
		// Generate private key
		pk, err = newPrivateKey(4096)
		if err != nil {
			return
		}

		// Save
		err = ioutil.WriteFile(pkPath, pk.Value(), 0600)
		if err != nil {
			return
		}
	}

	// Check if vault exists
	if fileExists("./.wtfvault") {
		// Load vault content if possible
		var value []byte
		value, err = ioutil.ReadFile("./.wtfvault")
		if err != nil {
			return
		}

		// Parse vault
		res = &vault{}
		err = json.Unmarshal(value, &res)
		if err != nil {
			return
		}
		res.PrivateKey = pk

		if load {
			// Find my recipient
			if pk != nil {
				var pubKey *publicKey
				pubKey, err = pk.PublicKey()
				if err != nil {
					return
				}
				myFingerprint := pubKey.Fingerprint()
				for _, r := range res.Recipients {
					if r.PublicKey.Fingerprint() == myFingerprint {
						// Decode encryption key
						var key []byte
						key, err = pk.Decrypt(r.SymKey)
						if err != nil {
							return
						}

						// Decode passwords
						var data []byte
						data, err = decryptSymetrical(res.EncodedPasswords, key)
						if err != nil {
							return
						}

						// Parse passwords
						err = json.Unmarshal(data, &res.Passwords)
						return
					}
				}
			}

			err = errors.New("you're not a member of the vault")
		}
	} else if init {
		// Generate vault file
		var username string
		username, err = PrivateKeyDefaultUsername()
		if err != nil {
			return
		}

		var pubKey *publicKey
		pubKey, err = pk.PublicKey()
		if err != nil {
			return
		}

		res = &vault{
			PrivateKey: pk,
			Recipients: []*recipient{
				{
					Name:      username,
					PublicKey: pubKey,
				},
			},
			Passwords: make(map[string]string),
		}
	} else {
		// No vault
		err = errors.New("no vault in directory")
	}
	return
}

// Save encodes and save as file.
func (v *vault) Save() (err error) {
	// Passwords to json
	data, err := json.Marshal(v.Passwords)
	if err != nil {
		return
	}

	// Generate new key
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		return
	}

	// Encode passwords
	v.EncodedPasswords, err = encryptSymetrical(data, key)
	if err != nil {
		return
	}

	// Set new key for everyone
	for _, r := range v.Recipients {
		r.SymKey, err = r.PublicKey.Encrypt(key)
		if err != nil {
			return
		}
	}

	// JSON vault
	data, err = json.Marshal(v)
	if err != nil {
		return
	}

	// Write
	err = ioutil.WriteFile("./.wtfvault", data, 0600)
	return
}

func encryptSymetrical(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptSymetrical(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
