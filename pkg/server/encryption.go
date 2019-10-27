package server

import (
	"bitbucket.org/zlacki/rscgo/pkg/isaac"
	rscrand "bitbucket.org/zlacki/rscgo/pkg/rand"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
)

//IsaacStream Container struct for 2 instances of the ISAAC+ CSPRNG, one for incoming data, the other outgoing data.
type IsaacStream struct {
	encoder, decoder *isaac.ISAAC
}

//SeedOpcodeCipher Initialize the ISAAC+ PRNG for use as a stream cipher for this client.
func (c *Client) SeedOpcodeCipher(clientSeed uint64, serverSeed uint64) *IsaacStream {
	if serverSeed != c.player.ServerSeed() {
		log.Warning.Printf("Session encryption key for command cipher received from client doesn't match the one we supplied it.\n")
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
