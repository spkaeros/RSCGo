package strutil

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

//MaxBase37 Max base37 string hash for 12-rune usernames. (999999999999)
const MaxBase37 = 6582952005840035281

// Presumably this charset is optimized to be in order of most-used in the English language.
// Not sure how Jagex came up with this arrangement, along with their interesting bit-packing methods involved here
var charset = []rune{' ', 'e', 't', 'a', 'o', 'i', 'h', 'n', 's', 'r', 'd', 'l', 'u', 'm', 'w', 'c', 'y', 'f', 'g', 'p', 'b', 'v', 'k', 'x', 'j', 'q', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ' ', '!', '?', '.', ',', ':', ';', '(', ')', '-', '&', '*', '\\', '\'', '@', '#', '+', '=', '\243', '$', '%', '"', '[', ']'}

func getCharCode(c rune) int {
	for i, cs := range charset {
		if cs == c {
			return i
		}
	}

	return 0
}

//IPToInteger Converts a string representation of an IPv4 address(e.g 127.0.0.1) to a 4-byte integer, each byte containing the information from one octet.
func IPToInteger(s string) (ip int) {
	if octets := strings.Split(s, "."); len(octets) > 0 {
		for index, octet := range octets {
			numericOctet, err := strconv.Atoi(octet)
			if err != nil {
				fmt.Println("Error parsing IP address to integer:", err)
				return -1
			}
			ip |= numericOctet << uint((3-index)*8)
		}
	}
	return ip
}

//ModalParse Neat command argument parsing function with support for single-quotes, ported from Java
func ModalParse(s string) []string {
	var cur string
	escaped := false
	quoted := false
	var out []string
	for _, c := range s {
		if c == '\\' {
			escaped = !escaped
			continue
		}
		if c == '\'' && !escaped {
			if quoted {
				if len(cur) > 0 {
					out = append(out, cur)
				}
				cur = ""
			}
			quoted = !quoted
			continue
		}
		if c == ' ' && !escaped && !quoted {
			if len(cur) > 0 {
				out = append(out, cur)
			}
			cur = ""
			continue
		}
		if escaped {
			escaped = false
		}
		cur += string(c)
	}
	if len(cur) > 0 {
		out = append(out, cur)
	}

	return out
}

var Base37 struct {
	Encode func(string) uint64
	Decode func(uint64) string
}

func init() {
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
				return 0xDEADBEEF
			}
		}
		return l
	}
	Base37.Decode = func(l uint64) string {
		if l < 0 || l >= MaxBase37 {
			return "invalid_name"
		}
		var s string
		for l != 0 {
			i := l % 37
			l /= 37
			if i == 0 {
				s = " " + s
			} else if i < 27 {
				if l%37 == 0 {
					s = string(i+64) + s
				} else {
					s = string(i+96) + s
				}
			} else {
				s = string(i+21) + s
			}
		}

		return s
	}
}

//PackChatMessage Takes a string as input, and returns a packed bitstream as output.  It can fit any of the first 13 runes in the alphabet into 4 bits per rune, but for the rest of the alphabet it's 8 bits.
func PackChatMessage(msg string) []byte {
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
			buf = append(buf, byte((cachedValue<<4)+code))
			cachedValue = -1
		} else {
			buf = append(buf, byte((cachedValue<<4)+(code>>4)))
			cachedValue = code & 0xF // LegsColor 4 bits
		}
	}
	if cachedValue != -1 {
		buf = append(buf, byte(cachedValue<<4))
	}

	return buf
}

//UnpackChatMessage Takes a byte slice as input, and decodes it into a human readable chat message.
//  Works by unpacking the bits into 2 4-bit values.  86% of the time, this is enough information for the output rune.
//  14% of the time, it will cache this value and use the next 4 bits to help decode the rune.
//  The first 13 runes in the charset map 1 to 1 for 13/15 possibilities within each 4 bits.
func UnpackChatMessage(data []byte) string {
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
			buf = append(buf, charset[byte(cachedValue<<4)+lowerHalf-195])
			cachedValue = -1
		}

		if cachedValue == -1 {
			if upperHalf < 13 {
				buf = append(buf, charset[upperHalf])
			} else {
				cachedValue = int(upperHalf)
			}
		} else {
			buf = append(buf, charset[byte(cachedValue<<4)+upperHalf-195])
			cachedValue = -1
		}
	}
	return string(buf)
}

//FormatChatMessage Format a given string for use with the in-game chat systems.
//Will replace certain symbols and auto-capitalize sentences.  Maybe more later.
func FormatChatMessage(msg string) string {
	buf := []rune(msg)
	flag := true
	for i, c := range msg {
		if (i > 4 && c == '@') || c == '%' {
			buf[i] = ' '
		}
		if i == 4 && c == '@' {
			flag = true
		}
		if flag && c >= 'a' && c <= 'z' {
			buf[i] = unicode.ToUpper(c)
			flag = false
		}
		if c == '.' || c == '!' || c == ':' {
			flag = true
		}
	}

	return string(buf)
}
