package bais

import (
	"fmt"
)

type IndexOutOfBoundsError struct {
	index       int
	arrayLength int
}

// Encode encodes a []byte into an alternative Base64 encoding as described above.
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
						sb = append(sb, byte(encode(accum>>18)))
						sb = append(sb, byte(encode(accum>>12)))
						sb = append(sb, byte(encode(accum>>6)))
						sb = append(sb, byte(encode(accum)))
						if i >= arrayLength {
							break
						}
					} else {
						sb = append(sb, byte(encode(accum>>10)))
						sb = append(sb, byte(encode(accum>>4)))
						sb = append(sb, byte(encode(accum<<2)))
						break
					}
				} else {
					sb = append(sb, byte(encode(accum>>2)))
					sb = append(sb, byte((encode(accum << 4))))
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

func encode(b int32) int32 {
	return (b+1)&63 + 63
}

func Decode(s string) ([]byte, error) {
	bb := make([]byte, 0, len(s))
	sb := []byte(s)
	for i := 0; i < len(sb); i++ {
		if sb[i] == byte('\b') {
			i++ // skip \b (backspace)
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
				if i < len(sb) && sb[i] != byte('!') {
					return []byte{}, fmt.Errorf(`expecting '!' got %d`, sb[i])
				}
				i++

				for i < len(sb) && sb[i] != byte('\b') {
					bb = append(bb, sb[i])
					i++
				}
				if i >= len(sb) {
					return bb[0:], nil
				}
				i++
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
