// Package escape maps unneeded Unicode characters to an ASCII equivalent and
// removes dangerous unprinting characters.
//
// Many applications do not need/want to support Unicode, but still wish to be
// able to be fed most Unicode input from their respective language. For
// instance, programs using ncurses are incapable of using wide characters
// without compile-time hacking. In addition, lookalike homoglyph characters
// can be used in social engineering attacks or to circumvent naive filters.
// This package attempts to act as a "temporary" solution to all these problems
// by converting such characters to a sane, ASCII equivalent.
//
// As UTF-8 is backward compatible with ASCII, all outputted strings are still
// valid UTF-8.
//
// Data sourced from http://homoglyphs.net/
//
// Deprecated: This library is not for use. You should really just support
// UTF-8 instead.
package escape

// Convertibles are runes which should be converted from the map key to the map
// value. Eg:
//	Convertibles['″'] --> '"'
//
// Although values may appear to be the same (as they are homoglyphs), they
// have a different internal representation and are NOT the same value.
var Convertibles = map[rune]rune{
	// Punctuation
	'″': '"',
	'‘': '\'',
	'’': '\'',
	'‚': ',',
	'⁎': '*',
	'‐': '-',
	'⁄': '/',
	'‹': '<',
	'›': '>',
	'⁓': '~',
	';': ';',

	// Cyrillic lookalikes (upper)
	'А': 'A',
	'В': 'B',
	'С': 'C',
	'Е': 'E',
	'Ԍ': 'G',
	'Н': 'H',
	'І': 'I',
	'Ј': 'J',
	'К': 'K',
	'М': 'M',
	'О': 'O',
	'Р': 'P',
	'Ѕ': 'S',
	'Т': 'T',
	'Ѵ': 'V',
	'Х': 'X',
	'Ү': 'Y',
	// Cyrillic lookalikes (lower)
	'а': 'a',
	'Ь': 'b',
	'с': 'c',
	'ԁ': 'd',
	'е': 'e',
	'һ': 'h',
	'і': 'i',
	'ј': 'j',
	'о': 'o',
	'р': 'p',
	'ѕ': 's',
	'ѵ': 'v',
	'ѡ': 'w',
	'х': 'x',
	'у': 'y',

	// Greek lookalikes (upper)
	'Α': 'A',
	'Β': 'B',
	'β': 'b',
	'Ϲ': 'C',
	'Ε': 'E',
	'Ϝ': 'F',
	'Η': 'H',
	'Ι': 'I',
	'Κ': 'K',
	'Μ': 'M',
	'Ν': 'N',
	'Ο': 'O',
	'Ρ': 'P',
	'Τ': 'T',
	'Χ': 'X',
	'Υ': 'Y',
	'Ζ': 'Z',
	// Greek lookalikes (lower)
	'ϲ': 'c',
	'ο': 'o',
	'ν': 'v',

	// Roman numerals (upper)
	'Ⅽ': 'C',
	'Ⅾ': 'D',
	'Ⅰ': 'I',
	'Ⅼ': 'L',
	'Ⅿ': 'M',
	'Ⅴ': 'V',
	'Ⅹ': 'X',
	// Roman numerals (lower)
	'ⅽ': 'c',
	'ⅾ': 'd',
	'ⅰ': 'i',
	'ⅼ': 'l',
	'ⅿ': 'm',
	'ⅴ': 'v',
	'ⅹ': 'x',

	// Fullwidth (punctuation)
	'！': '!',
	'＂': '"',
	'＄': '$',
	'％': '%',
	'＆': '&',
	'＇': '\'',
	'（': '(',
	'）': ')',
	'＊': '*',
	'＋': '+',
	'，': ',',
	'－': '-',
	'．': ',',
	'／': '/',
	'＼': '\\',
	'：': ':',
	'；': ';',
	'＜': '<',
	'＝': '=',
	'＞': '>',
	'？': '?',
	'＠': '@',
	'［': '[',
	'］': ']',
	'＾': '^',
	'＿': '_',
	'｀': '`',
	'｛': '{',
	'｝': '}',
	'～': '~',
	'｜': '|',
	// Fullwidth (numbers)
	'０': '0',
	'１': '1',
	'２': '2',
	'３': '3',
	'４': '4',
	'５': '5',
	'６': '6',
	'７': '7',
	'８': '8',
	'９': '9',
	// Fullwidth (upper)
	'Ａ': 'A',
	'Ｂ': 'B',
	'Ｃ': 'C',
	'Ｄ': 'D',
	'Ｅ': 'E',
	'Ｆ': 'F',
	'Ｇ': 'G',
	'Ｈ': 'H',
	'Ｉ': 'I',
	'Ｊ': 'J',
	'Ｋ': 'K',
	'Ｌ': 'L',
	'Ｍ': 'M',
	'Ｎ': 'N',
	'Ｏ': 'O',
	'Ｐ': 'P',
	'Ｑ': 'Q',
	'Ｒ': 'R',
	'Ｓ': 'S',
	'Ｔ': 'T',
	'Ｕ': 'U',
	'Ｖ': 'V',
	'Ｗ': 'W',
	'Ｘ': 'X',
	'Ｙ': 'Y',
	'Ｚ': 'Z',
	// Fullwidth (lower)
	'ａ': 'a',
	'ｂ': 'b',
	'ｃ': 'c',
	'ｄ': 'd',
	'ｅ': 'e',
	'ｆ': 'f',
	'ｇ': 'g',
	'ｈ': 'h',
	'ｉ': 'i',
	'ｊ': 'j',
	'ｋ': 'k',
	'ｌ': 'l',
	'ｍ': 'm',
	'ｎ': 'n',
	'ｏ': 'o',
	'ｐ': 'p',
	'ｑ': 'q',
	'ｒ': 'r',
	'ｓ': 's',
	'ｔ': 't',
	'ｕ': 'u',
	'ｖ': 'v',
	'ｗ': 'w',
	'ｘ': 'x',
	'ｙ': 'y',
	'ｚ': 'z',
}

// Escape returns s with all occurrences of unnecessary or confusing Unicode
// replaced with safe ASCII equivalents.
func Escape(s string) string {
	arr := []rune(s)
	for i, r := range arr {
		arr[i] = EscapeRune(r)
	}

	return string(arr)
}

func EscapeBytes(buf []byte) []byte {
	return []byte(Escape(string(buf)))
}

// EscapeRune returns r as a safe rune, meaning that it has been converted if
// an unsafe, unnecessary or confusing Unicode character to a safe ASCII
// equivalent.
func EscapeRune(r rune) rune {
	rep, ok := Convertibles[r]
	if ok {
		return rep
	}

	return r
}
