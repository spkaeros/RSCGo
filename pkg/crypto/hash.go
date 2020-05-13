package crypto

import (
	"encoding/hex"
	"runtime"

	"github.com/spkaeros/rscgo/pkg/config"
	"golang.org/x/crypto/argon2"
)

//Hash Takes a plaintext password as input, returns a hexadecimal string representation of the Argon2id hash as output.
func Hash(password string) string {
	return hex.EncodeToString(argon2.IDKey([]byte(password), []byte(config.HashSalt()), uint32(config.HashComplexity()), uint32(config.HashMemory()*1024), uint8(runtime.NumCPU()), uint32(config.HashLength())))
}

////NewHash Takes a plaintext password as input, returns a hexadecimal string representation of the Argon2id hash as output, with updated variables.
//func NewHash(password string) string {
//	return hex.EncodeToString(argon2.IDKey([]byte(password), []byte(config.HashSalt()), 15, uint32(8*1024), uint8(runtime.NumCPU()), uint32(config.HashLength())))
//}
