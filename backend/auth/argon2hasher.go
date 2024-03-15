package auth

// Referenced: https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/errx"
	"golang.org/x/crypto/argon2"
)

type Argon2Config struct {
	Memory     uint32 `json:"memory"`
	Iterations uint32 `json:"iterations"`
	Threads    uint8  `json:"numThreads"`
	SaltLen    uint32 `json:"saltLen"`
	KeyLen     uint32 `json:"keyLen"`
}

type argon2Hasher struct {
	config *Argon2Config
}

func NewArgon2Hasher() core.Hasher {
	return NewArgon2HasherWithConfig(&Argon2Config{
		Memory:     64 * 1024,
		Iterations: 3,
		Threads:    1,
		SaltLen:    16,
		KeyLen:     32,
	})
}

func NewArgon2HasherWithConfig(config *Argon2Config) core.Hasher {
	return &argon2Hasher{
		config: config,
	}
}

// Hash implements core.Hasher.
func (ah *argon2Hasher) Hash(pw string) (string, error) {
	saltBytes := make([]byte, ah.config.SaltLen)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", errx.Errf(err,
			"failed to generate random salt for hasing password")
	}

	hashBytes := argon2.IDKey(
		[]byte(pw),
		saltBytes,
		ah.config.Iterations,
		ah.config.Memory,
		ah.config.Threads,
		ah.config.KeyLen)

	hash := base64.StdEncoding.EncodeToString(hashBytes)
	salt := base64.StdEncoding.EncodeToString(saltBytes)

	hashForStorage := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		ah.config.Memory,
		ah.config.Iterations,
		ah.config.Threads,
		salt,
		hash)

	return hashForStorage, nil
}

// Verify implements core.Hasher.
func (*argon2Hasher) Verify(pw string, hashStr string) (bool, error) {
	comps := strings.Split(hashStr, "$")
	if len(comps) != 6 {
		return false, errors.New("invalid hash format given")
	}

	version, config := 0, Argon2Config{}
	if _, err := fmt.Sscanf(comps[2], "v=%d", &version); err != nil {
		return false, errx.Errf(err, "invalid format for argon version id")
	}
	if argon2.Version != version {
		return false, errx.Fmt("argon2 version mismatch, expected %d found %d",
			argon2.Version, version)
	}

	if _, err := fmt.Sscanf(comps[3], "m=%d,t=%d,p=%d", &config); err != nil {
		return false, errx.Errf(err, "argon2 config mismatch")
	}

	salt, err := base64.StdEncoding.Strict().DecodeString(comps[4])
	if err != nil {
		return false, errx.Errf(err, "invalid encoding detected for salt")
	}
	config.SaltLen = uint32(len(salt))

	hash, err := base64.StdEncoding.Strict().DecodeString(comps[5])
	if err != nil {
		return false, errx.Errf(err, "invalid encoding detected for hash")
	}
	config.KeyLen = uint32(len(hash))

	newHash := argon2.IDKey(
		[]byte(pw),
		salt,
		config.Iterations,
		config.Memory,
		config.Threads,
		config.KeyLen)

	if subtle.ConstantTimeCompare(hash, newHash) == 1 {
		return true, nil
	}
	return false, nil

}
