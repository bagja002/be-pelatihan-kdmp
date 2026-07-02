// Package hash implements Argon2id password hashing, the algorithm
// recommended by OWASP for password storage.
package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32 // KiB
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// defaultParams follow current OWASP guidance for Argon2id.
var defaultParams = params{
	memory:      64 * 1024, // 64 MiB
	iterations:  2,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

var (
	// ErrInvalidHash is returned when an encoded hash is malformed.
	ErrInvalidHash = errors.New("hash: invalid encoded hash")
	// ErrIncompatibleVersion is returned for a mismatched argon2 version.
	ErrIncompatibleVersion = errors.New("hash: incompatible argon2 version")
)

// Password hashes plain using Argon2id and returns a PHC-formatted string
// that embeds the algorithm, parameters, salt and derived key.
func Password(plain string) (string, error) {
	p := defaultParams
	salt := make([]byte, p.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	digest := argon2.IDKey([]byte(plain), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Digest := base64.RawStdEncoding.EncodeToString(digest)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Digest,
	), nil
}

// Verify reports whether plain matches the given encoded Argon2id hash.
// The comparison is constant-time to resist timing attacks.
func Verify(plain, encoded string) (bool, error) {
	p, salt, digest, err := decode(encoded)
	if err != nil {
		return false, err
	}
	other := argon2.IDKey([]byte(plain), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	if subtle.ConstantTimeEq(int32(len(digest)), int32(len(other))) == 0 {
		return false, nil
	}
	return subtle.ConstantTimeCompare(digest, other) == 1, nil
}

func decode(encoded string) (params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return params{}, nil, nil, ErrInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return params{}, nil, nil, ErrInvalidHash
	}
	if version != argon2.Version {
		return params{}, nil, nil, ErrIncompatibleVersion
	}

	var p params
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism); err != nil {
		return params{}, nil, nil, ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return params{}, nil, nil, ErrInvalidHash
	}
	digest, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return params{}, nil, nil, ErrInvalidHash
	}

	p.saltLength = uint32(len(salt))
	p.keyLength = uint32(len(digest))
	return p, salt, digest, nil
}
