package chem

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Composition maps an element symbol to its atom count in a formula.
type Composition map[string]int

var (
	trailingCharge = regexp.MustCompile(`[0-9]*[+\-]$`)
	chargeRe      = regexp.MustCompile(`^([0-9]*)([+\-])([0-9]*)$`)
)

// ParseFormula parses a chemical formula such as "H2O", "SO4^2-", "Fe2(SO4)3",
// "Ca(OH)2", "NH4+", "(NH4)2SO4" into an element composition and a net charge.
// Charge may appear after a caret ("^2-") or as a trailing token ("3+", "-").
func ParseFormula(s string) (Composition, int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, 0, fmt.Errorf("empty formula")
	}
	var compStr, chargeTok string
	if idx := strings.IndexByte(s, '^'); idx >= 0 {
		compStr = s[:idx]
		chargeTok = s[idx+1:]
	} else if m := trailingCharge.FindString(s); m != "" {
		compStr = s[:len(s)-len(m)]
		chargeTok = m
	} else {
		compStr = s
	}
	comp, err := parseComposition(compStr)
	if err != nil {
		return nil, 0, err
	}
	return comp, parseCharge(chargeTok), nil
}

func parseCharge(tok string) int {
	tok = strings.TrimSpace(tok)
	if tok == "" {
		return 0
	}
	m := chargeRe.FindStringSubmatch(tok)
	if m == nil {
		return 0
	}
	mag := 1
	if m[1] != "" {
		if v, err := strconv.Atoi(m[1]); err == nil {
			mag = v
		}
	} else if m[3] != "" {
		if v, err := strconv.Atoi(m[3]); err == nil {
			mag = v
		}
	}
	if mag == 0 {
		mag = 1
	}
	if m[2] == "-" {
		return -mag
	}
	return mag
}

func parseComposition(s string) (Composition, error) {
	pos := 0
	comp, err := parseBody(s, &pos)
	if err != nil {
		return nil, err
	}
	if pos != len(s) {
		return nil, fmt.Errorf("unexpected character %q at position %d in %q", s[pos], pos, s)
	}
	return comp, nil
}

func parseBody(s string, pos *int) (Composition, error) {
	result := Composition{}
	for *pos < len(s) {
		c := s[*pos]
		switch {
		case c == '(':
			*pos++
			sub, err := parseBody(s, pos)
			if err != nil {
				return nil, err
			}
			if *pos >= len(s) || s[*pos] != ')' {
				return nil, fmt.Errorf("missing closing parenthesis in %q", s)
			}
			*pos++ // consume ')'
			mult, err := parseCount(s, pos)
			if err != nil {
				return nil, err
			}
			for k, v := range sub {
				result[k] += v * mult
			}
		case c == ')':
			return result, nil
		case c >= 'A' && c <= 'Z':
			start := *pos
			*pos++
			if *pos < len(s) && s[*pos] >= 'a' && s[*pos] <= 'z' {
				*pos++
			}
			sym := s[start:*pos]
			if !IsElement(sym) {
				return nil, fmt.Errorf("unknown element symbol %q in %q", sym, s)
			}
			mult, err := parseCount(s, pos)
			if err != nil {
				return nil, err
			}
			result[sym] += mult
		default:
			return nil, fmt.Errorf("unexpected character %q in %q", c, s)
		}
	}
	return result, nil
}

func parseCount(s string, pos *int) (int, error) {
	start := *pos
	for *pos < len(s) && s[*pos] >= '0' && s[*pos] <= '9' {
		*pos++
	}
	if *pos == start {
		return 1, nil
	}
	n, err := strconv.Atoi(s[start:*pos])
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Elements returns the sorted union of element symbols present in the given compositions.
func Elements(comps ...Composition) []string {
	set := map[string]struct{}{}
	for _, c := range comps {
		for e := range c {
			set[e] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for e := range set {
		out = append(out, e)
	}
	return out
}
