package bais

import (
	"fmt"
)

// Encode encodes a []byte into an alternative Base64 encoding as described above.
func Encode(p []byte, allowControlCharacters bool) string {
	sb := make([]byte, 0, 3*len(p))

	for i := 0; i < len(p); i++ {
		if isAscii(p[i], allowControlCharacters) {
			sb = append(sb, p[i])
			continue
		}

		sb = append(sb, byte('\b'))
		for ; i < len(p); i++ {
			var accum int32
			accum = int32(p[i])
			if i == len(p)-1 { // last byte
				sb = append(sb, encode(accum>>2))
				sb = append(sb, (encode(accum << 4)))
				break
			} else {
				i++
				accum = accum<<8 | int32(p[i])
				if i == len(p)-1 {
					sb = append(sb, encode(accum>>10))
					sb = append(sb, encode(accum>>4))
					sb = append(sb, encode(accum<<2))
					break
				}
				i++
				accum = accum<<8 | int32(p[i])
				sb = append(sb, encode(accum>>18))
				sb = append(sb, encode(accum>>12))
				sb = append(sb, encode(accum>>6))
				sb = append(sb, encode(accum))
				if i == len(p)-1 {
					break
				}
			}
			if isAscii(p[i+1], allowControlCharacters) &&
				(i+2 >= len(p) || isAscii(p[i+1], allowControlCharacters)) &&
				(i+3 >= len(p) || isAscii(p[i+1], allowControlCharacters)) {
				sb = append(sb, '!')
				break
			}
		}
	}
	return string(sb)
}

func encode(b int32) byte {
	return byte((b+1)&63 + 63)
}
func isAscii(b byte, allowControlChars bool) bool {
	return b < 127 && (b >= 32 || (allowControlChars && b != '\b'))
}

func Decode(bb []byte, s string) ([]byte, error) {
	sb := []byte(s)
	for i := 0; i < len(sb); i++ {
		if sb[i] == '\b' {
			i++ // skip \b (backspace)
			for {
				cur := int32(sb[i])
				if i >= len(sb) {
					return bb, fmt.Errorf("index %d is greater than number of bytes in %s (%d)", i, s, len(s))
				}
				if cur < 63 || cur > 126 {
					return bb, fmt.Errorf("current byte (%d) is not ASCII, was expecting %[1]d to be > 63 and < 126", cur)
				}
				digit := (cur - 64) & 63
				zeros := 16 - 6
				accum := digit << zeros
				i++
				for i < len(sb) {
					if i >= len(sb) {
						break
					}
					cur = int32(sb[i])
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
					return bb, fmt.Errorf(`%d & 0xFF00 != 0`, accum)
				}
				if i < len(sb) && sb[i] != '!' {
					return bb, fmt.Errorf(`expecting '!' got %d`, sb[i])
				}
				i++

				for i < len(sb) && sb[i] != '\b' {
					bb = append(bb, sb[i])
					i++
				}
				if i >= len(sb) {
					return bb, nil
				}
				i++
			}
		}
		bb = append(bb, sb[i])
	}
	return bb, nil
}
