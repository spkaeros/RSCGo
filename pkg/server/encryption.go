package server

import (
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"os"

	"bitbucket.org/zlacki/rscgo/pkg/isaac"
	rscrand "bitbucket.org/zlacki/rscgo/pkg/rand"
)

var RsaKey *rsa.PrivateKey

type IsaacSeed struct {
	encoder, decoder *isaac.ISAAC
}

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

func (c *Client) SeedISAAC(seed []uint32) *IsaacSeed {
	if seed[2] != c.isaacSeed[2] || seed[3] != c.isaacSeed[3] {
		LogWarning.Printf("Session encryption key for command cipher received from client doesn't match the one we supplied it.\n")
		return nil
	}
	for i := 0; i < 2; i++ {
		c.isaacSeed[i] = seed[i]
	}
	for i := 4; i < 256; i += 4 {
		if i%2 == 0 {
			seed = append(seed, seed[2:4]...)
			seed = append(seed, seed[:2]...)
		} else {
			seed = append(seed, seed[:2]...)
			seed = append(seed, seed[2:4]...)
		}
	}
	decodingStream := isaac.New(seed)
	for i := 0; i < 256; i++ {
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
