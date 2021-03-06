package largo

const (
	space       = ' '
	singlequote = '\''
	doublequote = '"'
)

var nobreakspace = [2]byte{194, 160}

// Parse ...
func Tokenize(s string) (tokens []string) {
	return TokenizeBytes([]byte(s))
}

func TokenizeBytes(b []byte) (tokens []string) {
	t := new(parser)
	for i, c := range b {
		// See https://github.com/otiai10/largo/issues/1
		if c == nobreakspace[0] && i < len(b)-1 && b[i+1] == nobreakspace[1] {
			b[i+1] = space
			continue
		}
		t.push(c)
		if t.Closed {
			tokens = append(tokens, t.flush())
		}
	}
	if len(t.pool) != 0 {
		tokens = append(tokens, t.flush())
	}
	return
}

// token
type parser struct {
	delim  byte
	pool   []byte
	Closed bool
}

func (p *parser) flush() string {
	s := string(p.pool)
	p.delim = 0
	p.pool = nil
	p.Closed = false
	return s
}

func (p *parser) push(c byte) {
	switch c {
	case space:
		switch {
		case p.delim != 0:
			p.pool = append(p.pool, c)
		case len(p.pool) != 0:
			p.Closed = true
		}
	case singlequote:
		switch {
		case p.delim == singlequote: // Should close
			p.Closed = true
		case len(p.pool) == 0: // Should open
			p.delim = singlequote
		default: // This singlequote should be pushed as a value
			p.pool = append(p.pool, c)
		}
	case doublequote:
		switch {
		case p.delim == doublequote: // Should close
			p.Closed = true
		case len(p.pool) == 0: // Should open
			p.delim = doublequote
		default: // This doublequote should be pushed as a value
			p.pool = append(p.pool, c)
		}
	default:
		p.pool = append(p.pool, c)
	}
}
