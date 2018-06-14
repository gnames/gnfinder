package token

// Tokenize creates a slice containing every word in the document tokenized.
func Tokenize(text []rune) []Token {
	var res []Token
	space := spacesMap()
	start := 0
	dashToken := Token{}
	for i, v := range text {
		if _, ok := space[v]; ok {
			if _, isSpace := space[text[start]]; !isSpace {
				t := NewToken(text, start, i)
				if dashToken.Start > 0 {
					t = concatenateTokens(dashToken, t)
					dashToken.Start = 0
					res = append(res, t)
				} else {
					if lineEndsWithDash(text, i, t, space) {
						dashToken = t
					} else {
						res = append(res, t)
					}
				}
			}
			start = i + 1
		}
	}
	if len(text)-start > 1 {
		t := NewToken(text, start, len(text))
		res = append(res, t)
	}
	return res
}

func spacesMap() map[rune]byte {
	m := make(map[rune]byte)
	spaces := []rune("\n\r\v\t ")
	m[spaces[0]] = byte('\n')
	m[spaces[1]] = byte('\n')
	for i := 2; i < len(spaces); i++ {
		m[spaces[i]] = '\x00'
	}
	return m
}

func lineEndsWithDash(text []rune, i int, t Token, space map[rune]byte) bool {
	dash := rune('-')
	l := len(t.Raw)
	if l > 1 && t.Raw[l-1] == dash && lastWordForLine(text, i, space) {
		return true
	}
	return false
}

func lastWordForLine(text []rune, i int, space map[rune]byte) bool {
	for {
		if v, ok := space[text[i]]; ok {
			if v == '\n' {
				return true
			}
		} else {
			return false
		}
		i++
	}
}

func concatenateTokens(t1 Token, t2 Token) Token {
	var v []rune
	if t2.Raw[0] >= rune('a') && t2.Raw[0] <= rune('z') {
		v = append(t1.Raw[0:len(t1.Raw)-1], t2.Raw...)
	} else {
		v = append(t1.Raw, t2.Raw...)
	}
	t := Token{
		Raw:   v,
		Start: t1.Start,
		End:   t2.End,
	}
	t.Clean()
	return t
}
