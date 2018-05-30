// Package wordwrap provides methods for wrapping the contents of a string
package textCvtImge


import (
	"bytes"
	"unicode"

	"github.com/mattn/go-runewidth"
)

// WrapString wraps the given string within lim width in characters.
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word.
func WrapString(s string, lim uint) string {
	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	var current uint
	var wordBuf, spaceBuf bytes.Buffer
	var wordWidth, spaceWidth int

	for _, char := range s {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+uint(spaceWidth) > lim {
					current = 0
				} else {
					current += uint(spaceWidth)
					spaceBuf.WriteTo(buf)
					spaceWidth += runewidth.StringWidth(buf.String())
				}
				spaceBuf.Reset()
				spaceWidth = 0
			} else {
				current += uint(spaceWidth + wordWidth)
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				spaceWidth = 0
				wordWidth = 0
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += uint(spaceWidth + wordWidth)
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				spaceWidth = 0
				wordWidth = 0
			}

			spaceBuf.WriteRune(char)
			spaceWidth += runewidth.RuneWidth(char)
		} else if runewidth.RuneWidth(char) > 1 {

			// 双字节处理
			//fmt.Println(current + uint(spaceWidth+wordWidth))
			if current+uint(spaceWidth+wordWidth) >= lim {
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordWidth = 0
				current = 0
				spaceBuf.Reset()
				spaceWidth = 0
				buf.WriteRune('\n')
			} else {
				wordBuf.WriteRune(char)
				wordWidth += runewidth.RuneWidth(char)
				//fmt.Println(wordBuf.String())
				//fmt.Println(wordBuf.Len())
			}

			//fmt.Println(wordWidth)
		} else {

			wordBuf.WriteRune(char)
			wordWidth += runewidth.RuneWidth(char)

			if current+uint(spaceWidth+wordWidth) > lim && uint(wordWidth) < lim {
				buf.WriteRune('\n')
				current = 0
				spaceBuf.Reset()
				spaceWidth = 0
			}

		}
	}

	if wordBuf.Len() == 0 {
		if current+uint(spaceWidth) <= lim {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}

	return buf.String()
}
