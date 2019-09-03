package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"io/ioutil"

	"os"

	"bitbucket.org/zlacki/rscgo/pkg/isaac"
	rscrand "bitbucket.org/zlacki/rscgo/pkg/rand"
	"golang.org/x/crypto/sha3"
)

//RsaKey The RSA key for use in decoding the login packet
var RsaKey *rsa.PrivateKey

//ShakeHash The SHA3 hashing function state and reference point.
var ShakeHash sha3.ShakeHash

//IsaacStream Container struct for 2 instances of the ISAAC+ CSPRNG, one for incoming data, the other outgoing data.
type IsaacStream struct {
	encoder, decoder *isaac.ISAAC
}

//InitializeCrypto Read the RSA key into memory.
func InitializeCrypto() {
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
	ShakeHash = sha3.NewCShake256([]byte{}, []byte(TomlConfig.Crypto.HashSalt))
}

//SeedISAAC Initialize the ISAAC+ PRNG for use as a stream cipher for this client.
func (c *Client) SeedISAAC(clientSeed uint64, serverSeed uint64) *IsaacStream {
	if serverSeed != c.serverSeed {
		LogWarning.Printf("Session encryption key for command cipher received from client doesn't match the one we supplied it.\n")
		return nil
	}
	c.clientSeed = clientSeed
	decodingStream := isaac.New([]uint64{clientSeed, serverSeed})
	encodingStream := isaac.New([]uint64{clientSeed + 50, serverSeed + 50})

	return &IsaacStream{encodingStream, decodingStream}
}

//GenerateSessionID Generates a new 64-bit long using the systems CSPRNG.
//  For use as a seed with the ISAAC cipher (or similar secure stream cipher) used to encrypt packet data.
func GenerateSessionID() uint64 {
	return rscrand.Uint64S()
}

//HashPassword Takes a plaintext password as input, returns a hexidecimal string representation of the SHAKE256
//  hash of the input password.
func HashPassword(password string) string {
	if n, err := ShakeHash.Write([]byte(password)); n < len(password) || err != nil {
		LogWarning.Printf("HashPassword(string): Write failed:")
		if n < len(password) {
			LogWarning.Printf("Invalid length.  Expected %v, got %v\n", len(password), n)
			return "nil"
		}
		LogWarning.Printf("Error: %v", err.Error())
		return "nil"
	}
	dst := make([]byte, 64)
	if n, err := ShakeHash.Read(dst); n != 64 || err != nil {
		LogWarning.Println("HashPassword(string): Could not read hash back from shake function.", err)
		return "nil"
	}
	ShakeHash.Reset()

	return hex.EncodeToString(dst)
}
