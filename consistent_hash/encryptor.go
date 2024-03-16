package consistent_hash

import (
	"github.com/spaolacci/murmur3"
	"math"
)

type Encryptor interface {
	Encrypt(s string) int32
}

type MurmurHasher struct {
}

func NewMurmurHasher() *MurmurHasher {
	return &MurmurHasher{}
}
func (m *MurmurHasher) Encrypt(s string) int32 {
	hasher := murmur3.New32()
	_, _ = hasher.Write([]byte(s))
	return int32(hasher.Sum32() % math.MaxInt32)
}
