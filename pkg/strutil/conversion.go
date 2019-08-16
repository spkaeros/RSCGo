package strutil

import (
	"fmt"
	"strings"
	"unicode"
)

func init() {
	fmt.Printf("msg: %v, %v\n", PackChatMessage("lol cool story"), UnpackChatMessage(PackChatMessage("lol cool story")))
}

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
			cachedValue = code & 0xF // Bottom 4 bits
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
