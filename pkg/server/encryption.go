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

//InitializeHashing Initializes global hashing function reference for future usage, with `salt` as the salt.
func InitializeHashing(salt string) {
	ShakeHash = sha3.NewCShake256([]byte{}, []byte(salt))
}

//IsaacSeed Container struct for 2 instances of the ISAAC+ CSPRNG, one for incoming data, the other outgoing data.
type IsaacSeed struct {
	encoder, decoder *isaac.ISAAC
}

//ReadRSAKeyFile Read the RSA key from the file specified, within the DataDirectory.
func ReadRSAKeyFile(file string) {
	buf, err := ioutil.ReadFile(DataDirectory + string(os.PathSeparator) + file)
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

//SeedISAAC Initialize the ISAAC+ PRNG for use as a stream cipher for this client.
func (c *Client) SeedISAAC(seed []uint64) *IsaacSeed {
	if seed[1] != c.isaacSeed[1] {
		LogWarning.Printf("Session encryption key for command cipher received from client doesn't match the one we supplied it.\n")
		return nil
	}
	for i := 0; i < 2; i++ {
		c.isaacSeed[i] = seed[i]
	}
	decodingStream := isaac.New(seed)
	for i := 0; i < 2; i++ {
		seed[i] += 50
	}
	encodingStream := isaac.New(seed)

	return &IsaacSeed{encodingStream, decodingStream}
}

//GenerateSessionID Generates a new 64-bit long using the systems CSPRNG.
//  For use as a seed with the ISAAC cipher (or similar secure stream cipher) used to encrypt packet data.
func GenerateSessionID() uint64 {
	return rscrand.GetSecureRandomLong()
}

//HashPassword Takes a plaintext password as input, returns a hexidecimal string representation of the SHAKE256
//  hash of the input password.
func HashPassword(password string) string {
	ShakeHash.Reset()
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
