package strutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

//MaxBase37 Max base37 string hash for 12-rune usernames. (999999999999)
const MaxBase37 = 6582952005840035281

//IPToInteger Converts a string representation of an IPv4 address(e.g 127.0.0.1) to a 4-byte integer, each byte containing the information from one octet.
func IPToInteger(s string) (ip int) {
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
}

func IPToHexidecimal(s string) string {
	return Base16.String(uint64(IPToInteger(s)))
}

func JagHash(s string) int {
	ident := 0
	s = strings.ToUpper(s)
	for _, c := range s {
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
	charset := []rune{' ', 'e', 't', 'a', 'o', 'i', 'h', 'n', 's', 'r', 'd', 'l', 'u', 'm', 'w', 'c', 'y', 'f', 'g', 'p', 'b', 'v', 'k', 'x', 'j', 'q', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ' ', '!', '?', '.', ',', ':', ';', '(', ')', '-', '&', '*', '\\', '\'', '@', '#', '+', '=', '\243', '$', '%', '"', '[', ']'}
	getCharCode := func(c rune) int {
		for i, cs := range charset {
			if cs == c {
				return i
			}
		}
		return 0
	}
	ChatFilter.Pack = func(msg string) []byte {
		var buf []byte
		if len(msg) > 80 {
			msg = msg[:80]
		}
		msg = strings.ToLower(msg)
		cachedValue := -1
		for _, c := range msg {
			code := getCharCode(c)
			if code > 12 {
				code += 195
			}
			if cachedValue == -1 {
				if code < 13 {
					cachedValue = code
				} else {
					buf = append(buf, byte(code))
				}
			} else if code < 13 {
				buf = append(buf, byte((cachedValue<<4)|code)) // little end
				cachedValue = -1
			} else {
				buf = append(buf, byte((cachedValue<<4)|(code>>4)))
				cachedValue = code & 0xF // big end
			}
		}
		if cachedValue != -1 {
			buf = append(buf, byte(cachedValue<<4))
		}

		return buf
	}
	ChatFilter.Unpack = func(data []byte) string {
		// deprecated
		var buf []rune
		dataOffset := 0
		cachedValue := -1
		for i1 := 0; i1 < len(data); i1++ {
			nextChar := data[dataOffset] & 0xFF
			dataOffset++

			upperHalf := nextChar & 0xF      // Mask out first half of byte
			lowerHalf := nextChar >> 4 & 0xF // mask out last half of byte
			if cachedValue == -1 {
				if lowerHalf < 13 {
					buf = append(buf, charset[lowerHalf])
				} else {
					cachedValue = int(lowerHalf)
				}
			} else {
				buf = append(buf, charset[byte(cachedValue<<4)|lowerHalf-195])
				cachedValue = -1
			}

			if cachedValue == -1 {
				if upperHalf < 13 {
					buf = append(buf, charset[upperHalf])
				} else {
					cachedValue = int(upperHalf)
				}
			} else {
				buf = append(buf, charset[byte(cachedValue<<4)|upperHalf-195])
				cachedValue = -1
			}
		}
		return string(buf)
	}
	ChatFilter.Format = func(msg string) string {
		builder := &strings.Builder{}
		startingSentence := true
		for i, c := range msg {
			if unicode.IsGraphic(c) {
				if c == '@' {
					if i == 4 && msg[i-4] == '@' {
						startingSentence = true
					} else if i == 0 && msg[i+4] == '@' {
						startingSentence = false
					} else {
						c = ' '
					}
				} else if unicode.IsPunct(c) {
					startingSentence = true
				}
				startingSentence = false
				if startingSentence {
					c = unicode.ToUpper(c)
				} else {
					c = unicode.ToLower(c)
				}

				builder.WriteRune(c)
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
	/*
		BaseConversion.Encode = func(base int, s string) (l uint64) {
			for _, c := range s {
				l *= uint64(base)
				if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
					l += uint64(c-'a'+1)
				} else if c >= '0' && c <= '9' {
					l += uint64(c-'0'+27)
				}
			}

			return
		}

		BaseConversion.Decode = func(base int, i uint64) (s string) {
			if i < 0 {
				return "invalid_integer_to_string (enc failure)"
			}
			upper := true
			for i != 0 {
				remainder := i%uint64(base)
				i /= uint64(base)
				if remainder >= 11 {
					if upper {
						s = string(remainder + 'a'-1) + s
						upper = false
					} else {
						s = string(remainder + 'A'-1) + s
					}
				} else if remainder > 0 {
					s = string(remainder + '0' - 27) + s
				} else {
					s = string(' ') + s
					upper = true
				}
			}
			return s
		}
		fmt.Println(Base37.Decode(418444))
		fmt.Println(Base37.Encode(Base37.Decode(418444)))
		fmt.Println(Base16.String(418444))
		fmt.Println(Base16.Int(Base16.String(418444)))
	*/
}
