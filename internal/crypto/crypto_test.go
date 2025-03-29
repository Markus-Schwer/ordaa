package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePasswordHash(t *testing.T) {
	password := "supersecretpassword1234!"
	hash, err := GeneratePasswordHash(password)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := ComparePasswordAndHash(password, hash)
	if err != nil || !ok {
		t.Fatal("password hash comparison failed")
	}
}

func TestGeneratePasswordHash2(t *testing.T) {
	password := "test"
	hash, err := GeneratePasswordHash(password)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := ComparePasswordAndHash(password, hash)
	if err != nil || !ok {
		t.Fatal("password hash comparison failed")
	}
}

func TestGeneratePasswordHash3(t *testing.T) {
	password := "test"
	hash, err := GeneratePasswordHash(password)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := ComparePasswordAndHash("test2", hash)
	assert.True(t, err == nil || ok, "password hash comparison was successful, but should not have been")
}
