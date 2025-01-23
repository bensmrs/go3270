// This file is part of https://github.com/racingmars/go3270/
// Copyright 2020 by Matthew R. Wilson, licensed under the MIT license. See
// LICENSE in the project root for license information.

package go3270

import (
	"strings"
)

const sub = 0xfffd
const subS = "�"

// Each index in this array is the UTF-8 value, and the value at the index is
// the corresponding EBCDIC (codepage 37) / APL (codepage 310) value.
var encoder = map[string][]byte{
	"\x00": {0x00}, "\x01": {0x01}, "\x02": {0x02}, "\x03": {0x03},
	"\x04": {0x37}, "\x05": {0x2d}, "\x06": {0x2e}, "\x07": {0x2f},
	"\x08": {0x16}, "\x09": {0x05}, "\x0a": {0x15}, "\x0b": {0x0b},
	"\x0c": {0x0c}, "\x0d": {0x0d}, "\x0e": {0x0e}, "\x0f": {0x0f},

	"\x10": {0x10}, "\x11": {0x11}, "\x12": {0x12}, "\x13": {0x13},
	"\x14": {0x3c}, "\x15": {0x3d}, "\x16": {0x32}, "\x17": {0x26},
	"\x18": {0x18}, "\x19": {0x19}, "\x1a": {0x3f}, "\x1b": {0x27},
	"\x1c": {0x1c}, "\x1d": {0x1d}, "\x1e": {0x1e}, "\x1f": {0x1f},

	"\x7f": {0x07},

	"\x80": {0x20}, "\x81": {0x21}, "\x82": {0x22}, "\x83": {0x23},
	"\x84": {0x24}, "\x85": {0x25}, "\x86": {0x06}, "\x87": {0x17},
	"\x88": {0x28}, "\x89": {0x29}, "\x8a": {0x2a}, "\x8b": {0x2b},
	"\x8c": {0x2c}, "\x8d": {0x09}, "\x8e": {0x0a}, "\x8f": {0x1b},

	"\x90": {0x30}, "\x91": {0x31}, "\x92": {0x1a}, "\x93": {0x33},
	"\x94": {0x34}, "\x95": {0x35}, "\x96": {0x36}, "\x97": {0x08},
	"\x98": {0x38}, "\x99": {0x39}, "\x9a": {0x3a}, "\x9b": {0x3b},
	"\x9c": {0x04}, "\x9d": {0x14}, "\x9e": {0x3e}, "\x9f": {0xff},

	 " ": {0x40}, " ": {0x41},    "â": {0x42},  "ä": {0x43},
	 "à": {0x44}, "á": {0x45},    "ã": {0x46},  "å": {0x47},
	 "ç": {0x48}, "ñ": {0x49},    "¢": {0x4a},  ".": {0x4b},
	 "<": {0x4c}, "(": {0x4d},    "+": {0x4e},  "|": {0x4f},

	 "&": {0x50}, "é": {0x51},    "ê": {0x52},  "ë": {0x53},
	 "è": {0x54}, "í": {0x55},    "î": {0x56},  "ï": {0x57},
	 "ì": {0x58}, "ß": {0x59},    "!": {0x5a},  "$": {0x5b},
	 "*": {0x5c}, ")": {0x5d},    ";": {0x5e},  "¬": {0x5f},

	 "-": {0x60}, "/": {0x61},    "Â": {0x62},  "Ä": {0x63},
	 "À": {0x64}, "Á": {0x65},    "Ã": {0x66},  "Å": {0x67},
	 "Ç": {0x68}, "Ñ": {0x69},    "¦": {0x6a},  ",": {0x6b},
	 "%": {0x6c}, "_": {0x6d},    ">": {0x6e},  "?": {0x6f},

	 "ø": {0x70}, "É": {0x71},    "Ê": {0x72},  "Ë": {0x73},
	 "È": {0x74}, "Í": {0x75},    "Î": {0x76},  "Ï": {0x77},
	 "Ì": {0x78}, "`": {0x79},    ":": {0x7a},  "#": {0x7b},
	 "@": {0x7c}, "'": {0x7d},    "=": {0x7e}, "\"": {0x7f},

	 "Ø": {0x80}, "a": {0x81},    "b": {0x82},  "c": {0x83},
	 "d": {0x84}, "e": {0x85},    "f": {0x86},  "g": {0x87},
	 "h": {0x88}, "i": {0x89},    "«": {0x8a},  "»": {0x8b},
	 "ð": {0x8c}, "ý": {0x8d},    "þ": {0x8e},  "±": {0x8f},

	 "°": {0x90}, "j": {0x91},    "k": {0x92},  "l": {0x93},
	 "m": {0x94}, "n": {0x95},    "o": {0x96},  "p": {0x97},
	 "q": {0x98}, "r": {0x99},    "ª": {0x9a},  "º": {0x9b},
	 "æ": {0x9c}, "¸": {0x9d},    "Æ": {0x9e},  "¤": {0x9f},

	 "µ": {0xa0}, "~": {0xa1},    "s": {0xa2},  "t": {0xa3},
	 "u": {0xa4}, "v": {0xa5},    "w": {0xa6},  "x": {0xa7},
	 "y": {0xa8}, "z": {0xa9},    "¡": {0xaa},  "¿": {0xab},
	 "Ð": {0xac}, "Ý": {0xad},    "Þ": {0xae},  "®": {0xaf},

	 "^": {0xb0}, "£": {0xb1},    "¥": {0xb2},  "·": {0xb3},
	 "©": {0xb4}, "§": {0xb5},    "¶": {0xb6},  "¼": {0xb7},
	 "½": {0xb8}, "¾": {0xb9},    "[": {0xba},  "]": {0xbb},
	 "¯": {0xbc}, "¨": {0xbd},    "´": {0xbe},  "×": {0xbf},

	 "{": {0xc0}, "A": {0xc1},    "B": {0xc2},  "C": {0xc3},
	 "D": {0xc4}, "E": {0xc5},    "F": {0xc6},  "G": {0xc7},
	 "H": {0xc8}, "I": {0xc9}, "\xad": {0xca},  "ô": {0xcb},
	 "ö": {0xcc}, "ò": {0xcd},    "ó": {0xce},  "õ": {0xcf},

	 "}": {0xd0}, "J": {0xd1},    "K": {0xd2},  "L": {0xd3},
	 "M": {0xd4}, "N": {0xd5},    "O": {0xd6},  "P": {0xd7},
	 "Q": {0xd8}, "R": {0xd9},    "¹": {0xda},  "û": {0xdb},
	 "ü": {0xdc}, "ù": {0xdd},    "ú": {0xde},  "ÿ": {0xdf},

	"\\": {0xe0}, "÷": {0xe1},    "S": {0xe2},  "T": {0xe3},
	 "U": {0xe4}, "V": {0xe5},    "W": {0xe6},  "X": {0xe7},
	 "Y": {0xe8}, "Z": {0xe9},    "²": {0xea},  "Ô": {0xeb},
	 "Ö": {0xec}, "Ò": {0xed},    "Ó": {0xee},  "Õ": {0xef},

	 "0": {0xf0}, "1": {0xf1},    "2": {0xf2},  "3": {0xf3},
	 "4": {0xf4}, "5": {0xf5},    "6": {0xf6},  "7": {0xf7},
	 "8": {0xf8}, "9": {0xf9},    "³": {0xfa},  "Û": {0xfb},
	 "Ü": {0xfc}, "Ù": {0xfd},    "Ú": {0xfe},

	                   "𝐴̲": {0x08, 0x41}, "𝐵̲": {0x08, 0x42}, "𝐶̲": {0x08, 0x43},
	"𝐷̲": {0x08, 0x44}, "𝐸̲": {0x08, 0x45}, "𝐹̲": {0x08, 0x46}, "𝐺̲": {0x08, 0x47},
	"𝐻̲": {0x08, 0x48}, "𝐼̲": {0x08, 0x49},
	
	                   "𝐽̲": {0x08, 0x51}, "𝐾̲": {0x08, 0x52}, "𝐿̲": {0x08, 0x53},
	"𝑀̲": {0x08, 0x54}, "𝑁̲": {0x08, 0x55}, "𝑂̲": {0x08, 0x56}, "𝑃̲": {0x08, 0x57},
	"𝑄̲": {0x08, 0x58}, "𝑅̲": {0x08, 0x59},
	
	                                      "𝑆̲": {0x08, 0x62}, "𝑇̲": {0x08, 0x63},
	"𝑈̲": {0x08, 0x64}, "𝑉̲": {0x08, 0x65}, "𝑊̲": {0x08, 0x66}, "𝑋̲": {0x08, 0x67},
	"𝑌̲": {0x08, 0x68}, "𝑍̲": {0x08, 0x69},
	
	"◊": {0x08, 0x71}, "∧": {0x08, 0x72},                    "⌻": {0x08, 0x74},
	"⍸": {0x08, 0x75}, "⍷": {0x08, 0x76}, "⊢": {0x08, 0x77}, "⊣": {0x08, 0x78},
	"∨": {0x08, 0x79},
	
	"∼": {0x08, 0x80}, "║": {0x08, 0x81}, "═": {0x08, 0x82}, "⎸": {0x08, 0x83},
	"⎹": {0x08, 0x84}, "│": {0x08, 0x85},
	                                      "↑": {0x08, 0x8a}, "↓": {0x08, 0x8b},
	"≤": {0x08, 0x8c}, "⌈": {0x08, 0x8d}, "⌊": {0x08, 0x8e}, "→": {0x08, 0x8f},
	
	"⎕": {0x08, 0x90}, "▌": {0x08, 0x91}, "▐": {0x08, 0x92}, "▀": {0x08, 0x93},
	"▄": {0x08, 0x94}, "█": {0x08, 0x95},
	                                      "⊃": {0x08, 0x9a}, "⊂": {0x08, 0x9b},
	                   "○": {0x08, 0x9d},                    "←": {0x08, 0x9f},
	
	                                      "─": {0x08, 0xa2}, "∙": {0x08, 0xa3},
	"ₙ": {0x08, 0xa4},
	                                      "∩": {0x08, 0xaa}, "∪": {0x08, 0xab},
	"⊥": {0x08, 0xac},                    "≥": {0x08, 0xae}, "∘": {0x08, 0xaf},
	
	"⍺": {0x08, 0xb0}, "∊": {0x08, 0xb1}, "⍳": {0x08, 0xb2}, "⍴": {0x08, 0xb3},
	"⍵": {0x08, 0xb4},                                       "∖": {0x08, 0xb7},
	                                      "∇": {0x08, 0xba}, "∆": {0x08, 0xbb},
	"⊤": {0x08, 0xbc},                    "≠": {0x08, 0xbe}, "∣": {0x08, 0xbf},
	
	                   "⁽": {0x08, 0xc1}, "⁺": {0x08, 0xc2}, "■": {0x08, 0xc3},
	"└": {0x08, 0xc4}, "┌": {0x08, 0xc5}, "├": {0x08, 0xc6}, "┴": {0x08, 0xc7},
	                                      "⍲": {0x08, 0xca}, "⍱": {0x08, 0xcb},
	"⌷": {0x08, 0xcc}, "⌽": {0x08, 0xcd}, "⍂": {0x08, 0xce}, "⍉": {0x08, 0xcf},
	
	                   "⁾": {0x08, 0xd1}, "⁻": {0x08, 0xd2}, "┼": {0x08, 0xd3},
	"┘": {0x08, 0xd4}, "┐": {0x08, 0xd5}, "┤": {0x08, 0xd6}, "┬": {0x08, 0xd7},
	                                      "⌶": {0x08, 0xda}, "ǃ": {0x08, 0xdb},
	"⍒": {0x08, 0xdc}, "⍋": {0x08, 0xdd}, "⍞": {0x08, 0xde}, "⍝": {0x08, 0xdf},
	
	"≡": {0x08, 0xe0}, "₁": {0x08, 0xe1}, "₂": {0x08, 0xe2}, "₃": {0x08, 0xe3},
	"⍤": {0x08, 0xe4}, "⍥": {0x08, 0xe5}, "⍪": {0x08, 0xe6}, "€": {0x08, 0xe7},
	                                      "⌿": {0x08, 0xea}, "⍀": {0x08, 0xeb},
	"∵": {0x08, 0xec}, "⊖": {0x08, 0xed}, "⌹": {0x08, 0xee}, "⍕": {0x08, 0xef},
	
	"⁰": {0x08, 0xf0},
	"⁴": {0x08, 0xf4}, "⁵": {0x08, 0xf5}, "⁶": {0x08, 0xf6}, "⁷": {0x08, 0xf7},
	"⁸": {0x08, 0xf8}, "⁹": {0x08, 0xf9},                    "⍫": {0x08, 0xfb},
	"⍙": {0x08, 0xfc}, "⍟": {0x08, 0xfd}, "⍎": {0x08, 0xfe},
}

// Each index in this array is the EBCDIC (codepage 37) value, and the value
// at the index is the corresponding UTF-8 value.
var ebcdicDecoder = []rune{
	0x00, 0x01, 0x02, 0x03,  sub, 0x09,  sub, 0x7f,  sub,  sub,  sub, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13,  sub, 0x0a, 0x08,  sub, 0x18, 0x19,  sub,  sub, 0x1c, 0x1d, 0x1e, 0x1f,
	 sub,  sub,  sub,  sub,  sub,  sub, 0x17, 0x1b,  sub,  sub,  sub,  sub,  sub, 0x05, 0x06, 0x07,
	 sub,  sub, 0x16,  sub,  sub,  sub,  sub, 0x04,  sub,  sub,  sub,  sub, 0x14, 0x15,  sub, 0x1a,
	 ' ',  ' ',  'â',  'ä',  'à',  'á',  'ã',  'å',  'ç',  'ñ',  '¢',  '.',  '<',  '(',  '+',  '|',
	 '&',  'é',  'ê',  'ë',  'è',  'í',  'î',  'ï',  'ì',  'ß',  '!',  '$',  '*',  ')',  ';',  '¬',
	 '-',  '/',  'Â',  'Ä',  'À',  'Á',  'Ã',  'Å',  'Ç',  'Ñ',  '¦',  ',',  '%',  '_',  '>',  '?',
	 'ø',  'É',  'Ê',  'Ë',  'È',  'Í',  'Î',  'Ï',  'Ì',  '`',  ':',  '#',  '@', '\'',  '=',  '"',
	 'Ø',  'a',  'b',  'c',  'd',  'e',  'f',  'g',  'h',  'i',  '«',  '»',  'ð',  'ý',  'þ',  '±',
	 '°',  'j',  'k',  'l',  'm',  'n',  'o',  'p',  'q',  'r',  'ª',  'º',  'æ',  '¸',  'Æ',  '¤',
	 'µ',  '~',  's',  't',  'u',  'v',  'w',  'x',  'y',  'z',  '¡',  '¿',  'Ð',  'Ý',  'Þ',  '®',
	 '^',  '£',  '¥',  '·',  '©',  '§',  '¶',  '¼',  '½',  '¾',  '[',  ']',  '¯',  '¨',  '´',  '×',
	 '{',  'A',  'B',  'C',  'D',  'E',  'F',  'G',  'H',  'I', 0xad,  'ô',  'ö',  'ò',  'ó',  'õ',
	 '}',  'J',  'K',  'L',  'M',  'N',  'O',  'P',  'Q',  'R',  '¹',  'û',  'ü',  'ù',  'ú',  'ÿ',
	'\\',  '÷',  'S',  'T',  'U',  'V',  'W',  'X',  'Y',  'Z',  '²',  'Ô',  'Ö',  'Ò',  'Ó',  'Õ',
	 '0',  '1',  '2',  '3',  '4',  '5',  '6',  '7',  '8',  '9',  '³',  'Û',  'Ü',  'Ù',  'Ú', 0x9f,
}

// Each index in this array is the APL (codepage 310) value, and the value at
// the index is the corresponding UTF-8 value.
var aplDecoder = []string{
	subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS,
	subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS,
	subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS,
	subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS, subS,
	 " ",  "𝐴̲",  "𝐵̲",  "𝐶̲",  "𝐷̲",  "𝐸̲",  "𝐹̲",  "𝐺̲",  "𝐻̲",  "𝐼̲", subS, subS, subS, subS, subS, subS,
	subS,  "𝐽̲",  "𝐾̲",  "𝐿̲",  "𝑀̲",  "𝑁̲",  "𝑂̲",  "𝑃̲",  "𝑄̲",  "𝑅̲", subS, subS, subS, subS, subS, subS,
	subS, subS,  "𝑆̲",  "𝑇̲",  "𝑈̲",  "𝑉̲",  "𝑊̲",  "𝑋̲",  "𝑌̲",  "𝑍̲", subS, subS, subS, subS, subS, subS,
	 "◊",  "∧",  "¨",  "⌻",  "⍸",  "⍷",  "⊢",  "⊣",  "∨", subS, subS, subS, subS, subS, subS, subS,
	 "∼",  "║",  "═",  "⎸",  "⎹",  "│", subS, subS, subS, subS,  "↑",  "↓",  "≤",  "⌈",  "⌊",  "→",
	 "⎕",  "▌",  "▐",  "▀",  "▄",  "█", subS, subS, subS, subS,  "⊃",  "⊂",  "¤",  "○",  "±",  "←",
	 "¯",  "°",  "─",  "∙",  "ₙ", subS, subS, subS, subS, subS,  "∩",  "∪",  "⊥",  "[",  "≥",  "∘",
	 "⍺",  "∊",  "⍳",  "⍴",  "⍵", subS,  "×",  "∖",  "÷", subS,  "∇",  "∆",  "⊤",  "]",  "≠",  "∣",
	 "{",  "⁽",  "⁺",  "■",  "└",  "┌",  "├",  "┴",  "§", subS,  "⍲",  "⍱",  "⌷",  "⌽",  "⍂",  "⍉",
	 "}",  "⁾",  "⁻",  "┼",  "┘",  "┐",  "┤",  "┬",  "¶", subS,  "⌶",  "ǃ",  "⍒",  "⍋",  "⍞",  "⍝",
	 "≡",  "₁",  "₂",  "₃",  "⍤",  "⍥",  "⍪",  "€", subS, subS,  "⌿",  "⍀",  "∵",  "⊖",  "⌹",  "⍕",
	 "⁰",  "¹",  "²",  "³",  "⁴",  "⁵",  "⁶",  "⁷",  "⁸",  "⁹", subS,  "⍫",  "⍙",  "⍟",  "⍎", subS,
}

// encode converts the UTF-8 input string, u, to an EBCDIC byte array.
func encode(u string) []byte {
	var result []byte
	var lastRune *rune

	for _, r := range u {
		if lastRune != nil {
			c, ok := encoder[string([]rune{*lastRune, r})]
			lastRune = nil
			if ok {
				result = append(result, c...)
				continue
			} else {
				result = append(result, '\x3f')
			}
		}

		c, ok := encoder[string(r)]
		if ok {
			result = append(result, c...)
		} else {
			lastRune = &r
		}
	}

	if lastRune != nil {
		result = append(result, '\x3f')
	}

	return result
}

// decode converts the EBCDIC input byte array, e, to a UTF-8 string.
func decode(e []byte) string {
	var sb strings.Builder
	var apl = false
	for _, c := range e {
		if c == '\x08' {
			apl = true
			continue
		}

		if apl {
			sb.WriteString(aplDecoder[c])
			apl = false
		} else {
			sb.WriteRune(ebcdicDecoder[c])
		}
	}

	return sb.String()
}
