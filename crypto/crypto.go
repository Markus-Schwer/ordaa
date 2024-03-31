package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

// HashSalt struct used to store
// generated hash and salt used to
// generate the hash.
type HashSalt struct {
	Hash, Salt []byte
}

type Argon2idHash struct {
	// time represents the number of
	// passed over the specified memory.
	time uint32
	// cpu memory to be used.
	memory uint32
	// threads for parallelism aspect
	// of the algorithm.
	threads uint8
	// keyLen of the generate hash key.
	keyLen uint32
	// saltLen the length of the salt used.
	saltLen uint32
}

// NewArgon2idHash constructor function for
// Argon2idHash.
func NewArgon2idHash(time, saltLen uint32, memory uint32, threads uint8, keyLen uint32) *Argon2idHash {
	return &Argon2idHash{
		time:    time,
		saltLen: saltLen,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
}

func NewDefaultArgon2idHash() *Argon2idHash {
	return &Argon2idHash{
		time:    1,
		saltLen: 32,
		memory:  64 * 1024,
		threads: 32,
		keyLen:  256,
	}
}

func randomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// GenerateHash using the password and provided salt.
// If not salt value provided fallback to random value
// generated of a given length.
func (a *Argon2idHash) generateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	// If salt is not provided generate a salt of
	// the configured salt length.
	if len(salt) == 0 {
		salt, err = randomSecret(a.saltLen)
	}
	if err != nil {
		return nil, err
	}
	// Generate hash
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)
	// Return the generated hash and salt used for storage.
	return &HashSalt{Hash: hash, Salt: salt}, nil
}

func GeneratePasswordHash(password string) (string, error) {
	argon2IDHash := NewDefaultArgon2idHash()
	salt, err := randomSecret(argon2IDHash.saltLen)
	if err != nil {
		return "", err
	}
	hashSalt, err := argon2IDHash.generateHash([]byte(password), salt)
	if err != nil {
		return "", err
	}

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(hashSalt.Salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hashSalt.Hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argon2IDHash.memory, argon2IDHash.time, argon2IDHash.threads, b64Salt, b64Hash)

	return encodedHash, nil
}

func ComparePasswordAndHash(password string, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	argon2IDHash, hashSalt, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash, err := argon2IDHash.generateHash([]byte(password), hashSalt.Salt)
	if err != nil {
		return false, err
	}

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hashSalt.Hash, otherHash.Hash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (argon2IDHash *Argon2idHash, hashSalt *HashSalt, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}

	argon2IDHash = &Argon2idHash{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &argon2IDHash.memory, &argon2IDHash.time, &argon2IDHash.threads)
	if err != nil {
		return nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}
	argon2IDHash.saltLen = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
	}
	argon2IDHash.keyLen = uint32(len(hash))

	return argon2IDHash, &HashSalt{Hash: hash, Salt: salt}, nil
}
