package vault

import "fmt"

// UserAddCmd adds an user.
func UserAddCmd(args map[string]interface{}, askSecure func(question ...interface{}) string) error {
	v, err := newVault(false, true)
	if err != nil {
		return err
	}

	// Create public key from params
	name := args["name"].(string)
	publicKey := publicKey{value: args["public_key"].(string)}

	// Verify public key
	_, err = publicKey.RsaKey()
	if err != nil {
		return err
	}

	// Check existing
	for _, r := range v.Recipients {
		if r.Name == name {
			r.PublicKey = &publicKey
			return v.Save()
		}
	}

	// Add
	v.Recipients = append(v.Recipients, &recipient{
		PublicKey: &publicKey,
		Name:      name,
	})
	return v.Save()
}

// UserRemoveCmd removes an user.
func UserRemoveCmd(args map[string]interface{}, askSecure func(question ...interface{}) string) error {
	v, err := newVault(false, false)
	if err != nil {
		return err
	}

	// Param
	param := args["public_key"].(string)
	isFingerprint := len(param) == 32

	res := []*recipient{}
	for _, r := range v.Recipients {
		if !(r.Name == param || r.PublicKey.Value() == param || (isFingerprint && r.PublicKey.Fingerprint() == param)) {
			res = append(res, r)
		}
	}
	v.Recipients = res
	return v.Save()
}

// UserListCmd lists users.
func UserListCmd(args map[string]interface{}, askSecure func(question ...interface{}) string) error {
	v, err := newVault(false, false)
	if err != nil {
		return err
	}

	fmt.Println("Allowed users in vault:")
	for _, r := range v.Recipients {
		fmt.Printf("     - %s: %s \n", r.Name, r.PublicKey.Fingerprint())
	}
	return nil
}

// PublicKeyCmd prints public key.
func PublicKeyCmd(args map[string]interface{}, askSecure func(question ...interface{}) string) error {
	v, err := newVault(true, false)
	if err != nil {
		return err
	}
	fmt.Println("A")
	pubKey, err := v.PrivateKey.PublicKey()
	if err != nil {
		return err
	}

	if args["raw"].(bool) {
		fmt.Println(pubKey.Value())
		return nil
	}

	fmt.Println("Your public key to share is:")
	fmt.Println(pubKey.Value())
	fmt.Println("")
	fmt.Println("It's fingerrpint is:", pubKey.Fingerprint())
	fmt.Println("Import it to the vault with 'wtf --vault user-add <public-key>'")
	return nil
}

// SetCmd sets a value in the store.
func SetCmd(args map[string]interface{}, askSecure func(question ...interface{}) string) error {
	v, err := newVault(true, true)
	if err != nil {
		return err
	}

	// Params
	key := args["key"].(string)
	value := args["value"].(string)
	ask := args["ask"].(string)

	if ask != "" {
		value = askSecure(ask)
	}

	if value == "" {
		delete(v.Passwords, key)
	} else {
		v.Passwords[key] = value
	}
	return v.Save()
}

// GetCmd gets a value in the store.
func GetCmd(args map[string]interface{}, askSecure func(question ...interface{}) string) error {
	v, err := newVault(false, true)
	if err != nil {
		return err
	}

	// Params
	key := args["key"].(string)
	required := args["required"].(bool)

	value, ok := v.Passwords[key]
	if ok {
		fmt.Print(value)
	} else if required {
		return fmt.Errorf("Missing vault value '%s'", key)
	}
	return nil
}
