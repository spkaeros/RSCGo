package crypto

import (
	"fmt"
	"math/big"
)

func init() {

//	test := []byte { byte(1000 >> 8) & 0xFF, byte(1000&0xFF)}
//	fmt.Println(test)
//	test = EncryptRSA(test)
//	fmt.Println(test)
//	test = DecryptRSA(test)
//	fmt.Println(test)
//	fmt.Println(Encrypt(test), "\n", Decrypt(test))
}

var mod = checkErr(new(big.Int).SetString("121727957757863576101561860005285292626079113266451124223351027655156011238254177877652729098983576274837395085392103662934978533548660507677480253506715648449246069310428873797293036242698272731265802720055413141023411018398284944110799347717001885969188133010830020603318079626459849229187149790609728348667", 10))
var priv = checkErr(new(big.Int).SetString("39574754648544995073200949288114212279763356643380221180569460664855696662160882512059875754822652518160202239533245623544435635274825672366080266671146999893044984788215998341126495384938227705193172944227710071442663553068050822068193588197727598992397886903673557185565299947088350518351942016798869358001", 10))
var pub = checkErr(new(big.Int).SetString("12256971504525176577999115521306614075749098639988274452692554670619288210288814203087336665303501555493198422881032409199392946347224070978354126295353401", 10))

func EncryptRSA(data []byte) []byte {
	b := new(big.Int).SetBytes(data)
	b = b.Exp(b, pub, mod)
	return b.Bytes()
}

func DecryptRSA(data []byte) []byte {
	b := new(big.Int).SetBytes(data)
	
	b = b.Exp(b, priv, mod)
	return b.Bytes()
}

func checkErr(b *big.Int, ok bool) *big.Int {
	if !ok {
		fmt.Println("Ran into some error initializing RSA:")
		return nil
	}
	return b
}