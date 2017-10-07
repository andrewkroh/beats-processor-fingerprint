package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/sha3"
)

type hashAlgorithm struct {
	name     string
	producer hashProducer
}

type hashProducer func() hash.Hash

var hashProducerTypes = map[string]hashProducer{
	"md5":        md5.New,
	"sha1":       sha1.New,
	"sha224":     sha256.New224,
	"sha256":     sha256.New,
	"sha384":     sha512.New384,
	"sha512":     sha512.New,
	"sha512_224": sha512.New512_224,
	"sha512_256": sha512.New512_256,
	"sha3_224":   sha3.New224,
	"sha3_256":   sha3.New256,
	"sha3_384":   sha3.New384,
	"sha3_512":   sha3.New512,
}

func (a *hashAlgorithm) Unpack(s string) error {
	s = strings.ToLower(s)
	p, found := hashProducerTypes[strings.ToLower(s)]
	if !found {
		return fmt.Errorf("invalid hash type '%v'", s)
	}

	*a = hashAlgorithm{name: s, producer: p}
	return nil
}

type encodingType encoder

type encoder func([]byte) string

var encodingTypes = map[string]encodingType{
	"hex":    hex.EncodeToString,
	"base32": base32.StdEncoding.EncodeToString,
	"base64": base64.StdEncoding.EncodeToString,
}

func (t *encodingType) Unpack(s string) error {
	e, found := encodingTypes[strings.ToLower(s)]
	if !found {
		return fmt.Errorf("invalid encoding type '%v'", s)
	}

	*t = e
	return nil
}

type FingerprintConfig struct {
	// Hash function used for calculating the fingerprint.
	// This value is case-insensitive. The default is sha256.
	Hash hashAlgorithm `config:"hash"`

	// Encoding type for the output. Either hex or base64. Default is hex.
	Encoder encodingType `config:"encoding"`

	// Fields is a list of fields whose values are to be used as the input to
	// the hash function. The field values are concatenated before the hashing
	// is performed. All fields must be present in the event otherwise an error
	// will be returned by the filter. The default is message.
	Fields []string `config:"fields"`

	// Target field for the hash value. The default is fingerprint.
	Target string `config:"target"`
}

var defaultFingerprintConfig = FingerprintConfig{
	Fields:  []string{"message"},
	Target:  "fingerprint",
	Hash:    hashAlgorithm{"sha256", sha256.New},
	Encoder: hex.EncodeToString,
}
