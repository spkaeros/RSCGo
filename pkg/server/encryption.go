package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"io/ioutil"
	"runtime"

	"os"

	"bitbucket.org/zlacki/rscgo/pkg/isaac"
	rscrand "bitbucket.org/zlacki/rscgo/pkg/rand"
	"golang.org/x/crypto/argon2"
)

//RsaKey The RSA key for use in decoding the login packet
var RsaKey *rsa.PrivateKey

//IsaacStream Container struct for 2 instances of the ISAAC+ CSPRNG, one for incoming data, the other outgoing data.
type IsaacStream struct {
	encoder, decoder *isaac.ISAAC
}

//loadRsaKey Read the RSA key into memory.
func loadRsaKey() {
	buf, err := ioutil.ReadFile(TomlConfig.DataDir + TomlConfig.Crypto.RsaKeyFile)
	if err != nil {
		LogError.Printf("Could not read RSA key from file:%v", err)
		os.Exit(103)
	}
	key, err := x509.ParsePKCS8PrivateKey(buf)
	if err != nil {
		LogWarning.Printf("Could not parse RSA key:%v", err)
		os.Exit(104)
	}
	RsaKey = key.(*rsa.PrivateKey)
}

//SeedOpcodeCipher Initialize the ISAAC+ PRNG for use as a stream cipher for this client.
func (c *Client) SeedOpcodeCipher(clientSeed uint64, serverSeed uint64) *IsaacStream {
	if serverSeed != c.player.ServerSeed() {
		LogWarning.Printf("Session encryption key for command cipher received from client doesn't match the one we supplied it.\n")
		return nil
	}
	decodingStream := isaac.New([]uint64{clientSeed, serverSeed})
	encodingStream := isaac.New([]uint64{clientSeed + 50, serverSeed + 50})

	return &IsaacStream{encodingStream, decodingStream}
}

//GenerateSessionID Generates a new 64-bit long using the systems CSPRNG.
// For use as a seed with the ISAAC cipher (or similar secure stream cipher) used to encrypt packet data.
func GenerateSessionID() uint64 {
	return rscrand.Uint64S()
}

//HashPassword Takes a plaintext password as input, returns a hexidecimal string representation of the SHAKE256 hash as output.
func HashPassword(password string) string {
	return hex.EncodeToString(argon2.IDKey([]byte(password), []byte(TomlConfig.Crypto.HashSalt), uint32(TomlConfig.Crypto.HashComplexity), uint32(TomlConfig.Crypto.HashMemory*1024), uint8(runtime.NumCPU()), uint32(TomlConfig.Crypto.HashLength)))
}

//DecryptRSABlock Attempts to decrypt the payload buffer.  Returns the decrypted buffer upon success, otherwise returns nil.
func DecryptRSABlock(payload []byte) []byte {
	buf, err := rsa.DecryptPKCS1v15(rand.Reader, RsaKey, payload)
	if err != nil {
		LogWarning.Println("Could not decrypt RSA block:", err)
		return nil
	}
	return buf
}
