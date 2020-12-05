package strutil

import (
	"math"
	"net"
	"strings"
	"unicode"
	// "github.com/spkaeros/rscgo/pkg/log"
)

//MaxBase37 Max base37 string hash for 12-rune usernames. (999999999999)
const MaxBase37 = 6582952005840035281

// this shit is whack.
// Not even positive I'm using it in any active codepath really,
// seems to alter exotic chars to other exotic-er chars, and accessed
// in the chat crypto routines
var specialMap = map[rune]rune{
	'€': 'ﾀ',
	'?': 'ﾝ',
	'‚': 'ﾂ',
	'ƒ': 'ﾃ',
	'„': 'ﾄ',
	'…': 'ﾅ',
	'†': 'ﾆ',
	'‡': 'ﾇ',
	'ˆ': 'ﾈ',
	'‰': 'ﾉ',
	'Š': 'ﾊ',
	'‹': 'ﾋ',
	'Œ': 'ﾌ',
	'Ž': 'ﾎ',
	'‘': 'ﾑ',
	'’': 'ﾒ',
	'“': 'ﾓ',
	'”': 'ﾔ',
	'•': 'ﾕ',
	'–': 'ﾖ',
	'—': 'ﾗ',
	'˜': 'ﾘ',
	'™': 'ﾙ',
	'š': 'ﾚ',
	'›': 'ﾛ',
	'œ': 'ﾜ',
	'ž': 'ﾞ',
	'Ÿ': 'ﾟ',
}

// These two slices contain data which apparently ends up mixing up the chat strings
// It gets generated at runtime via the init routine below
var cipherDictionary []int
var cipherData []int

// This data here is the number of bits each char from 0-255 corresponds to in this super strange string encryption
var shiftCounts = []int{22, 22, 22, 22, 22, 22, 21, 22, 22, 20, 22, 22, 22, 21, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 3, 8, 22, 16, 22, 16, 17, 7, 13, 13, 13, 16,
	7, 10, 6, 16, 10, 11, 12, 12, 12, 12, 13, 13, 14, 14, 11, 14, 19, 15, 17, 8, 11, 9, 10, 10, 10, 10, 11, 10,
	9, 7, 12, 11, 10, 10, 9, 10, 10, 12, 10, 9, 8, 12, 12, 9, 14, 8, 12, 17, 16, 17, 22, 13, 21, 4, 7, 6, 5, 3,
	6, 6, 5, 4, 10, 7, 5, 6, 4, 4, 6, 10, 5, 4, 4, 5, 7, 6, 10, 6, 10, 22, 19, 22, 14, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 21, 22, 21, 22, 22, 22, 21, 22, 22}

//IPToInteger Converts a string representation of an IPv4 address(e.g 127.0.0.1) to a 4-byte integer, each byte containing the information from one octet.
func IPToInteger(s string) (ip net.IP) {
	return net.ParseIP(s)
	/*	if strings.HasPrefix(s, "[") {
			// IPv6, I think
			if octets := strings.Split(strings.Replace(strings.Replace(s, "[", "", -1), "]", "", -1), ":"); len(octets) > 0 {
				for index, octet := range octets {
					numericOctet, err := strconv.Atoi(strings.Split(octet, ":")[0])
					if err != nil {
						fmt.Println("Error parsing IP address to integer:", err)
						return -1
					}
					ip |= numericOctet << uint((len(octets)-index)*8)
				}
			}
		}
		if octets := strings.Split(s, "."); len(octets) > 0 {
			for index, octet := range octets {
				numericOctet, err := strconv.Atoi(strings.Split(octet, ":")[0])
				if err != nil {
					fmt.Println("Error parsing IP address to integer:", err)
					return -1
				}
				ip |= numericOctet << uint((3-index)*8)
			}
		}
		return ip
	*/
}

func IPToHexidecimal(s string) string {
	return net.ParseIP(s).String()
}

//JagHash implements a string hashing function for file names contained in a jagball archive.
func JagHash(s string) int {
	ident := 0
	for _, c := range strings.ToUpper(s) {
		ident = int((rune(ident)*61 + c) - 32)
	}
	return ident
}

//ParseArgs Neat command argument parsing function with support for single-quotes, ported from Java
func ParseArgs(s string) []string {
	str := strings.Builder{}
	escaped := false
	quoted := false
	var out []string
	for _, c := range s {
		if c == '\\' {
			escaped = !escaped
			continue
		}
		if c == '\'' && !escaped {
			if quoted && str.Len() > 0 {
				out = append(out, str.String())
				str.Reset()
			}
			quoted = !quoted
			continue
		}
		if c == ' ' && !escaped && !quoted {
			if str.Len() > 0 {
				out = append(out, str.String())
				str.Reset()
			}
			continue
		}
		if escaped {
			escaped = false
		}
		str.WriteRune(c)
	}

	if str.Len() > 0 {
		out = append(out, str.String())
	}
	return out
}

//CombatPrefix Returns the chat prefix to colorize combat levels in right click menus and such.
// The color fades red as the target compares better than you, or fades green as the target compares worse than you.
// White indicates an equal target.
func CombatPrefix(delta int) string {
	// They're stronger
	if delta < -9 {
		return "@red@"
	}
	if delta < -6 {
		return "@or3@"
	}
	if delta < -3 {
		return "@or2@"
	}
	if delta < 0 {
		return "@or1@"
	}

	// They're weaker
	if delta > 9 {
		return "@gre@"
	}
	if delta > 6 {
		return "@gr3@"
	}
	if delta > 3 {
		return "@gr2@"
	}
	if delta > 0 {
		return "@gr1@"
	}

	// They're the same
	return "@whi@"
}

//ChatFilter Represents a single API access point for encoding and decoding chat messages using
var ChatFilter struct {
	//Pack Takes a string as input, and returns a packed bitstream as output.  It can fit any of the first 13 runes in the alphabet into 4 bits per rune, but for the rest of the alphabet it's 8 bits.
	Pack func(string) []byte
	//Unpack Takes a byte slice as input, and decodes it into a human readable chat message.
	//  Works by unpacking the bits into 2 4-bit values.  86% of the time, this is enough information for the output rune.
	//  14% of the time, it will cache this value and use the next 4 bits to help decode the rune.
	//  The first 13 runes in the charset map 1 to 1 for 13/15 possibilities within each 4 bits.
	Unpack func([]byte) string
	//Format Format a given string for use with the in-game chat systems.
	// Will replace certain symbols and auto-capitalize sentences.  Maybe more later.
	Format func(string) string
}

//Base37 Represents a single API access point for encoding strings to their base 37 integer hash and converting base 37 hashes back to strings.
var Base37 struct {
	//Encode Encodes a string into its base 37 integer hash and returns the hash.  Input will be filtered appropriately.
	Encode func(string) uint64
	//Decode Decodes a string from its base 37 integer hash form back into a Go string and returns it, placing capital letters where appropriate.
	Decode func(uint64) string
}

var Base16 struct {
	//Int returns an integer representation of the provided base-16 string
	Int func(string) uint64
	//String returns the base-16 string representation of the provided hexidecimal integer
	String func(uint64) string
}

var Base2 struct {
	//Int returns an integer representation of the provided base-2 string
	Int func(string) uint64
	//String returns the base-2 string representation of the provided hexidecimal integer
	String func(uint64) string
}

var BaseConversion struct {
	//Int returns an integer representation of the provided string using the provided base.
	Int func(int, string) uint64
	//String returns a string encoding of the provided integer using the provided base.
	String func(int, uint64) string
}

func init() {
	// Presumably this charset is optimized to be in order of most-used in the English language, as I think I've encountered this character array before elsewhere and that was its stated design
	validChar := func(c byte) bool {
		charset := []rune{
				' ', 'e', 't', 'a', 'o', 'i', 'h', 'n', 's', 'r', 'd', 'l', 'u', 'm', 'w',
				'c', 'y', 'f', 'g', 'p', 'b', 'v', 'k', 'x', 'j', 'q', 'z', '0', '1', '2',
				'3', '4', '5', '6', '7', '8', '9', ' ', '!', '?', '.', ',', ':', ';', '(',
				')', '-', '&', '*', '\\', '\'', '@', '#', '+', '=', '\243', '$', '%', '"',
				'[', ']' }
		for _, cs := range charset {
			if cs == rune(c) {
				return true
			}
		}
		return false
	}
	endsSentence := func(c byte) bool {
		punct := []rune{'!', '?', '.', ':'}
		for _, cs := range punct {
			if cs == rune(c) {
				return true
			}
		}
		return false
	}
	// ChatFilter.Pack = func(msg string) []byte {
		// var buf []byte
		// if len(msg) > 80 {
			// msg = msg[:80]
		// }
		// msg = strings.ToLower(msg)
		// cachedValue := -1
		// for _, c := range msg {
			// code := getCharCode(c)
			// if code > 12 {
				// code += 195
			// }
			// if cachedValue == -1 {
				// if code < 13 {
					// cachedValue = code
				// } else {
					// buf = append(buf, byte(code))
				// }
			// } else if code < 13 {
				// buf = append(buf, byte((cachedValue<<4)|code)) // little end
				// cachedValue = -1
			// } else {
				// buf = append(buf, byte((cachedValue<<4)|(code>>4)))
				// cachedValue = code & 0xF // big end
			// }
		// }
		// if cachedValue != -1 {
			// buf = append(buf, byte(cachedValue<<4))
		// }
// 
		// return buf
	// }
	// ChatFilter.Unpack = func(data []byte) string {
		// // deprecated
		// var buf []rune
		// dataOffset := 0
		// cachedValue := -1
		// for i1 := 0; i1 < len(data); i1++ {
			// nextChar := data[dataOffset] & 0xFF
			// dataOffset++
// 
			// upperHalf := nextChar & 0xF      // Mask out first half of byte
			// lowerHalf := nextChar >> 4 & 0xF // mask out last half of byte
			// if cachedValue == -1 {
				// if lowerHalf < 13 {
					// buf = append(buf, charset[lowerHalf])
				// } else {
					// cachedValue = int(lowerHalf)
				// }
			// } else {
				// buf = append(buf, charset[byte(cachedValue<<4)|lowerHalf-195])
				// cachedValue = -1
			// }
// 
			// if cachedValue == -1 {
				// if upperHalf < 13 {
					// buf = append(buf, charset[upperHalf])
				// } else {
					// cachedValue = int(upperHalf)
				// }
			// } else {
				// buf = append(buf, charset[byte(cachedValue<<4)|upperHalf-195])
				// cachedValue = -1
			// }
		// }
		// return string(buf)
	// }
	ChatFilter.Format = func(msg string) string {
		builder := &strings.Builder{}
		startingSentence := true
		caret := 0
		msg = strings.ToLower(msg)
		for i := 0; i < len(msg); i += 1 {
			caret += 1
			if validChar(msg[i]) {
				if msg[i] == '@' {
					if caret == 1 && msg[i+4] == '@' {
						builder.WriteString(msg[i:i+5])
						i += 4
						startingSentence = true
					} else {
						builder.WriteByte(' ')
					}
				} else if msg[i] == '~' {
					if i == 0 && msg[i+4] == '~' {
						i += 4
						caret = 0
					} else {
						builder.WriteByte(' ')
					}
				} else if endsSentence(msg[i]) {
					startingSentence = true
					builder.WriteByte(msg[i])
				} else if unicode.IsSpace(rune(msg[i])) || unicode.IsNumber(rune(msg[i])) {
					builder.WriteByte(msg[i])
				} else if startingSentence && !unicode.IsSpace(rune(msg[i])) {
					builder.WriteRune(unicode.ToUpper(rune(msg[i])))
					startingSentence = false
				} else {
					builder.WriteRune(unicode.ToLower(rune(msg[i])))
					startingSentence = false
				}
			}
		}

		return builder.String()
	}
	Base37.Encode = func(s string) uint64 {
		s = strings.ToLower(s)
		var buf []rune
		for _, c := range s {
			if c >= 'a' && c <= 'z' {
				buf = append(buf, c)
			} else if c >= '0' && c <= '9' {
				buf = append(buf, c)
			} else {
				buf = append(buf, ' ')
			}
		}

		s1 := strings.TrimSpace(string(buf))
		if len(s1) > 12 {
			s1 = s1[:12]
		}
		var l uint64
		for _, c := range s1 {
			l *= 37
			if c >= 'a' && c <= 'z' {
				l += 1 + uint64(c) - 97
			} else if c >= '0' && c <= '9' {
				l += 27 + uint64(c) - 48
			}
			if l >= MaxBase37 {
				return MaxBase37
			}
		}
		return l
	}

	Base37.Decode = func(i uint64) string {
		if i < 0 || i >= math.MaxUint64 {
			return "invalid_name"
		}
		var s string
		for i != 0 {
			remainder := i % 37
			i /= 37
			if remainder == 0 {
				s = " " + s
			} else if remainder < 27 {
				if i%37 == 0 {
					s = string(remainder+64) + s
				} else {
					s = string(remainder+96) + s
				}
			} else {
				s = string(remainder+21) + s
			}
		}

		return s
	}

	Base16.Int = func(s string) uint64 {
		if s[0] == '0' && s[1] == 'x' {
			s = s[2:]
		}
		return BaseConversion.Int(16, s)
	}

	Base16.String = func(i uint64) string {
		return "0x" + BaseConversion.String(16, i)
	}

	Base2.Int = func(s string) uint64 {
		if s[0] == '0' && s[1] == 'b' {
			s = s[2:]
		}
		return BaseConversion.Int(2, s)
	}

	Base2.String = func(i uint64) string {
		return "0b" + BaseConversion.String(2, i)
	}

	BaseConversion.Int = func(base int, s string) (l uint64) {
		for _, c := range s {
			l *= uint64(base)
			if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
				l += uint64(c - 'A' + 10)
			} else if c >= '0' && c <= '9' {
				l += uint64(c - '0')
			}
		}

		return
	}

	BaseConversion.String = func(base int, i uint64) (s string) {
		if i < 0 {
			return "invalid_integer_to_string (enc failure)"
		}
		for i != 0 {
			remainder := i % uint64(base)
			i /= uint64(base)
			if remainder >= 10 {
				s = string(remainder+'A'-10) + s
			} else {
				s = string(remainder+'0') + s
			}
		}
		return s
	}
	blockBuilder := make([]int, 33)
	cipherDictIndexTemp := 0
	for initPos := 0; initPos < len(shiftCounts); initPos++ {
		initValue := shiftCounts[initPos]
		builderBitSelector := 1 << (32 - uint(initValue))
		builderValue := blockBuilder[initValue]
		cipherData = append(cipherData, builderValue)
		builderValueBit := 0
		if (builderValue & builderBitSelector) == 0 {
			builderValueBit = builderValue | builderBitSelector
			for initValueCounter := initValue - 1; initValueCounter > 0; initValueCounter-- {
				builderValue2 := blockBuilder[initValueCounter]
				if builderValue != builderValue2 {
					break
				}
				builderValue2BitSelector := 1 << (32 - uint(initValueCounter))
				if (builderValue2 & builderValue2BitSelector) == 0 {
					blockBuilder[initValueCounter] = builderValue2BitSelector | builderValue2
				} else {
					blockBuilder[initValueCounter] = blockBuilder[initValueCounter-1]
					break
				}
			}
		} else {
			builderValueBit = blockBuilder[initValue-1]
		}
		blockBuilder[initValue] = builderValueBit
		for initValueCounter := initValue + 1; initValueCounter <= 32; initValueCounter++ {
			if builderValue == blockBuilder[initValueCounter] {
				blockBuilder[initValueCounter] = builderValueBit
			}
		}
		cipherDictIndex := 0
		for initValueCounter := 0; initValueCounter < initValue; initValueCounter++ {
			builderBitSelector2 := 0x80000000 >> uint(initValueCounter)
			if (builderValue & builderBitSelector2) == 0 {
				cipherDictIndex++
			} else {
				if cipherDictionary[cipherDictIndex] == 0 {
					cipherDictionary[cipherDictIndex] = cipherDictIndexTemp
				}
				cipherDictIndex = cipherDictionary[cipherDictIndex]
			}
			for len(cipherDictionary) <= cipherDictIndex {
				cipherDictionary = append(cipherDictionary, 0)
			}
		}
		cipherDictionary[(cipherDictIndex)] = ^initPos
		if cipherDictIndex >= cipherDictIndexTemp {
			cipherDictIndexTemp = cipherDictIndex + 1
		}
	}
}

func convertMsg(txt string) []rune {
	buf := []rune{}
	for _, b := range []rune(txt) {
		if b >= 128 && b < 160 {
			buf = append(buf, rune(specialMap[b]))
			continue
		}
		buf = append(buf, rune(b))
	}
	return buf
}

func Decipher(msg []byte, decipheredLength int) string {
	bufferIndex := 0
	off := 0
	decipherIndex := 0
	var cipherDictValue int
	buffer := make([]int, decipheredLength)

	for bufferIndex < decipheredLength && len(msg)-off > 0 {
		encipheredByte := int8(msg[off])
		off++
		if encipheredByte < 0 {
			decipherIndex = cipherDictionary[decipherIndex]
		} else {
			decipherIndex = decipherIndex + 1
		}

		if cipherDictValue = cipherDictionary[decipherIndex]; 0 > (cipherDictValue) {
			buffer[bufferIndex] = ^cipherDictValue
			bufferIndex += 1
			if bufferIndex >= decipheredLength {
				break
			}

			decipherIndex = 0
		}

		for andVal := 0x40; andVal > 0; andVal >>= 1 {
			if encipheredByte&int8(andVal) == 0 {
				decipherIndex = decipherIndex + 1
			} else {
				decipherIndex = cipherDictionary[decipherIndex]
			}
			if cipherDictValue = cipherDictionary[decipherIndex]; cipherDictValue < 0 {
				buffer[bufferIndex] = ^cipherDictValue
				bufferIndex += 1
				if bufferIndex >= decipheredLength {
					break
				}

				decipherIndex = 0
			}
		}
	}

	s := make([]byte, decipheredLength)
	//Swap bytes for unicode characters
	for bufferIndex = 0; bufferIndex < decipheredLength; bufferIndex += 1 {
		bufferValue := buffer[bufferIndex] & 0xFF
		if bufferValue != 0 {
			if bufferValue >= 128 && bufferValue < 160 {
				bufferValue = int(specialMap[rune(bufferValue-128)])
			}

			s[bufferIndex] = byte(bufferValue)
		}
	}
	return string(s[:])
}

func Encipher(txt string) ([]byte, int) {
	buf := convertMsg(txt)
	output := make([]int, 0, len(txt))
	enciphered := 0
	bitOffset := 0
	for _, c := range buf {
		v := cipherData[c]
		i1 := shiftCounts[c]

		off := bitOffset >> 3
		shift := bitOffset & 7
		enciphered &= -shift >> 31
		endOff := off + ((shift + i1 - 1) >> 3)
		bitOffset += i1
		shift += 24
		enciphered |= v >> uint(shift)
		for len(output) <= off {
			output = append(output, 0)
		}
		output[off] = enciphered
		if endOff > off {
			off++
			shift -= 8
			enciphered = v >> uint(shift)
			for len(output) <= off {
				output = append(output, 0)
			}
			output[off] = enciphered & 0xFF
			if off < endOff {
				shift -= 8
				off++
				enciphered = v >> uint(shift)
				for len(output) <= off {
					output = append(output, 0)
				}
				output[off] = enciphered & 0xFF
				if endOff > off {
					shift -= 8
					off++
					enciphered = v >> uint(shift)
					for len(output) <= off {
						output = append(output, 0)
					}
					output[off] = enciphered & 0xFF
					if endOff > off {
						shift -= 8
						off++
						enciphered = v << uint(-shift)
						for len(output) <= off {
							output = append(output, 0)
						}
						output[off] = enciphered & 0xFF
					}
				}
			}
		}
	}

	return toBytes(output[:(7+bitOffset)>>3]), len(txt)
}

func toBytes(arr []int) (out []byte) {
	out = make([]byte, len(arr))
	for i, v := range arr {
		out[i] = byte(v)
	}
	return out
}
