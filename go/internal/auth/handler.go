package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"splitsavvy/internal/database"
	"strings"

	"golang.org/x/crypto/argon2"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Params struct {
    Memory      uint32
    Iterations  uint32
    Parallelism uint8
    SaltLength  uint32
    KeyLength   uint32
}

var (
    ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
    ErrIncompatibleVariant = errors.New("incompatible variant of argon2")
    ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func DecodeHash(hash string) (params *Params, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, ErrIncompatibleVariant
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params = &Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}

func ComparePasswordAndHash(password, hash string) (match bool, err error) {
	match, _, err = CheckHash(password, hash)
	return match, err
}

func CheckHash(password, hash string) (match bool, params *Params, err error) {
	params, salt, key, err := DecodeHash(hash)
	if err != nil {
		return false, nil, err
	}

	otherKey := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	keyLen := int32(len(key))
	otherKeyLen := int32(len(otherKey))

	if subtle.ConstantTimeEq(keyLen, otherKeyLen) == 0 {
		return false, params, nil
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 {
		return true, params, nil
	}
	return false, params, nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "Invalid Request Format", http.StatusBadRequest)
		return
	}
	if len(req.Identifier) == 0 {
        http.Error(w, "Identifier is required", http.StatusBadRequest)
        return
    }
	var(
		passwordHash	string
		username		string
		firstName 		string
		phoneVerified	bool
	)
	isPhoneLogin := false
	var query string
	if req.Identifier[0] == '+'{
		isPhoneLogin = true
		query = `
			SELECT hashed_password, username, first_name, phone_number_verified
			FROM users
			WHERE phone_number = $1`
	} else {
		query = `
			SELECT hashed_password, username, first_name, phone_number_verified
			FROM users
			WHERE username = $1`
	}
	err := database.DB.QueryRow(r.Context(), query, req.Identifier).Scan(&passwordHash, &username, &firstName, &phoneVerified)
	if err != nil{
		http.Error(w, "User not found", http.StatusUnauthorized)
        return
	}
	if isPhoneLogin && !phoneVerified {
        http.Error(w, "Phone number not verified. Please login with Username.", http.StatusForbidden)
        return
    }

	match, err := ComparePasswordAndHash(req.Password, passwordHash)
    if err != nil {
        http.Error(w, "Error verifying credentials", http.StatusInternalServerError)
        return
    }

	if !match {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

	response := map[string]string{
        "username":   username,
        "first_name": firstName,
		"status":     "success",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

}