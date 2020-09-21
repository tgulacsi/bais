// <summary>Encodes and decodes BAIS (Byte Array In String) encoding,
// which preserves runs of ASCII characters unchanged. This encoding is
// useful for debugging (since ASCII runs are visible) and for conversion
// of bytes to JSON.</summary>
// <remarks>
// Arrays encoded with <see cref="ByteArrayInString.Convert(ArraySlice{byte}, bool)"/>
// tend to be slightly more compact than standard Uuencoding or Base64,
// and when you use this encoding in JSON with UTF-8, the output is
// typically also more compact than yEnc since double-byte characters
// above 127 are avoided.
// <para/>
// A BAIS string alternates between runs of "direct" bytes (usually bytes
// in the ASCII range that are represented as themselves) and runs of a
// special base-64 encoding. The base-64 encoding is a sequence of 6-bit
// digits with 64 added to them, except for 63 which is mapped to itself.
// This is easier and faster to encode and decode than standard Base64
// and has an interesting property described below.
// <para/>
// A BAIS string begins in ASCII mode and switches to base 64 when the '\b'
// character is encountered. Base-64 mode ends, returning to ASCII, when a
// '!' character is encountered.
// <para/>
// For example:
// <pre>
//   //                    C   a    t       \n  E        A   B   C   D
//   var b = new byte[] { 67, 97, 116, 128, 10, 69, 255, 65, 66, 67, 68 };
//   Assert.AreEqual(ByteArrayInString.Convert(b), "Cat\b`@iE?tEB!CD");
// </pre>
// A byte sequence such as 128, 10, 69, 255 can be encoded in base 64 as
// illustrated:
// <pre>
//              ---128---    ---10----    ---69----  ---255---
//   Bytes:     1000 0000    0000 1010    0100 0101  1111 1111
//   Base 64:   100000   000000   101001    000101   111111   110000
//   Encoded: 01100000 01000000 01101001  01000101 01111111 01110000
//            ---96--- ---64--- --105---  ---69--- --127--- --112---
//               `        @        i         E        ~        p
// </pre>
// <para/>
// An interesting property of this base-64 encoding is that when it encodes
// bytes between 63 and 126, those bytes appear unchanged at certain
// offsets (specifically the third, sixth, ninth, etc.) In this example,
// since the third byte is 'E' (69), it also appears as 'E' in the
// output.
// <para/>
// When viewing BAIS strings, another thing to keep in mind is that
// runs of zeroes ('\0') will tend to appear as runs of `@` characters
// in the base 64 encoding, although a single zero is not always enough
// to make a `@` appear. Runs of 255 will tend to appear as runs of `?`.
// <para/>
// There are many ways to encode a given byte array as BAIS.
// </remarks>
package bais

import (
	"fmt"
)

type IndexOutOfBoundsError struct {
	index       int
	arrayLength int
}

func (ioobe IndexOutOfBoundsError) Error() string {
	panic(fmt.Sprintf("%d is not a valid index[0,%d] for the array", ioobe.index, ioobe.arrayLength))
}

func Encode(ba *[]byte, allowControlCharacters bool) string {
	bytes := *ba
	arrayLength := len(bytes)
	sb := make([]byte, 0, arrayLength*3)
	i := 0
	var b byte
	for i < arrayLength {
		if isAscii(bytes, i, allowControlCharacters) {
			b, i = getByte(bytes, i)
			sb = append(sb, b)
		} else {
			b, i = getByte(bytes, i)
			sb = append(sb, byte('\b'))
			for ; ; b, i = getByte(bytes, i) {
				var accum int32
				accum = int32(b)
				if i < arrayLength {
					b, i = getByte(bytes, i)
					accum = (accum << 8) | int32(b)
					if i < arrayLength {
						b, i = getByte(bytes, i)
						accum = accum<<8 | int32(b)
						sb = append(sb, byte(encodeBase64Digit(accum>>18)))
						sb = append(sb, byte(encodeBase64Digit(accum>>12)))
						sb = append(sb, byte(encodeBase64Digit(accum>>6)))
						sb = append(sb, byte(encodeBase64Digit(accum)))
						if i >= arrayLength {
							break
						}
					} else {
						sb = append(sb, byte(encodeBase64Digit(accum>>10)))
						sb = append(sb, byte(encodeBase64Digit(accum>>4)))
						sb = append(sb, byte(encodeBase64Digit(accum<<2)))
						break
					}
				} else {
					sb = append(sb, byte(encodeBase64Digit(accum>>2)))
					sb = append(sb, byte((encodeBase64Digit(accum << 4))))
					break
				}
				if isAscii(bytes, i, allowControlCharacters) &&
					(i+1 >= len(bytes) || isAscii(bytes, i, allowControlCharacters)) &&
					(i+2 >= len(bytes) || isAscii(bytes, i, allowControlCharacters)) {
					sb = append(sb, byte('!'))
					break
				}
			}
		}
	}
	return string(sb[0:])
}

func encodeBase64Digit(b int32) int32 {
	return (b+1)&63 + 63
}

func Decode(s string) ([]byte, error) {
	bb := make([]byte, 0, len(s))
	sb := []byte(s)
	for i := 0; i < len(sb); i++ {
		if sb[i] == '\b' {
			i++ // skip \b
			for {
				cur := int32(sb[i]) & 0xFF
				if i >= len(sb) {
					return []byte{}, fmt.Errorf("index %d is greater than number of bytes in %s (%d)", i, s, len(s))
				}
				if cur < 63 || cur > 126 {
					return []byte{}, fmt.Errorf("current byte (%d) is not ASCII, was expecting %[1]d to be > 63 and < 126", cur)
				}
				digit := (cur - 64) & 63
				zeros := 16 - 6
				accum := digit << zeros
				i++
				for i < len(sb) {
					if i >= len(sb) {
						break
					}
					cur = int32(sb[i]) & 0xFF
					if cur < 63 || cur > 126 {
						break
					}
					digit = (cur - 64) & 63
					zeros = zeros - 6
					accum = accum | digit<<zeros
					if zeros <= 8 {
						bb = append(bb, byte(accum>>8))
						accum = accum << 8
						zeros = zeros + 8
					}
					i++
				}

				if accum&0xFF00 != 0 {
					return []byte{}, fmt.Errorf(`%d & 0xFF00 != 0`, accum)
				}
				if i < len(sb) && sb[i] != '!' {
					return []byte{}, fmt.Errorf(`%s expecting '!' got %d`, sb[i])
				}
				i++

				for i < len(sb) && sb[i] != '\b' {
					bb = append(bb, sb[i])
					i++
				}
				if i >= len(sb) {
					return bb[0:], nil
				}
			}
		}
		bb = append(bb, sb[i])
	}
	return bb[0:], nil
}

func getByte(bytes []byte, index int) (b byte, nextIndex int) {
	if index < len(bytes) {
		return bytes[index] & 0xFF, index + 1
	} else {
		panic(&IndexOutOfBoundsError{index, len(bytes)})
	}
}

func isAscii(bytes []byte, index int, allowControlChars bool) bool {
	b := bytes[index]
	return b < 127 && (b >= 32 || (allowControlChars && b != '\b'))
}
