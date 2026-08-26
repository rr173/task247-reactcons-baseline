package chem

import (
	"fmt"
	"strconv"
	"strings"
)

// Term is one side of a parsed reaction equation.
type Term struct {
	Symbol string
	Coeff  float64
	Role   string // "reactant" | "product"
}

// ParseEquation parses a reaction equation such as "2 H2 + O2 -> 2 H2O" or
// "A + B <=> C". It returns the reactant and product terms and whether the
// reaction is reversible.
func ParseEquation(eq string) (reactants, products []Term, reversible bool, err error) {
	lhs, rhs, rev, ok := findArrow(eq)
	if !ok {
		return nil, nil, false, fmt.Errorf("no reaction arrow (->, <=>, =>, →) found in %q", eq)
	}
	reactants, err = parseSide(lhs, "reactant")
	if err != nil {
		return nil, nil, false, err
	}
	products, err = parseSide(rhs, "product")
	if err != nil {
		return nil, nil, false, err
	}
	if len(reactants) == 0 || len(products) == 0 {
		return nil, nil, false, fmt.Errorf("reaction must have reactants and products in %q", eq)
	}
	return reactants, products, rev, nil
}

func findArrow(s string) (string, string, bool, bool) {
	type cand struct {
		token string
		rev   bool
	}
	order := []cand{
		{"<=>", true},
		{"<->", true},
		{"=>", false},
		{"→", false},
		{"⇌", true},
		{"->", false},
	}
	for _, c := range order {
		if i := strings.Index(s, c.token); i >= 0 {
			return s[:i], s[i+len(c.token):], c.rev, true
		}
	}
	return "", "", false, false
}

func parseSide(side string, role string) ([]Term, error) {
	var terms []Term
	for _, raw := range splitSide(side) {
		t, err := parseTerm(raw, role)
		if err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	return terms, nil
}

func splitSide(side string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(side); i++ {
		switch side[i] {
		case '(':
			depth++
		case ')':
			depth--
		case '+':
			if depth == 0 {
				parts = append(parts, side[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, side[start:])
	return parts
}

func parseTerm(term string, role string) (Term, error) {
	term = strings.TrimSpace(term)
	if term == "" {
		return Term{}, fmt.Errorf("empty reaction term")
	}
	i := 0
	for i < len(term) && (term[i] == '.' || (term[i] >= '0' && term[i] <= '9')) {
		i++
	}
	coeff := 1.0
	if i > 0 {
		c, err := strconv.ParseFloat(term[:i], 64)
		if err != nil {
			return Term{}, fmt.Errorf("invalid coefficient in %q", term)
		}
		coeff = c
	}
	symbol := strings.TrimSpace(term[i:])
	if symbol == "" {
		return Term{}, fmt.Errorf("missing species symbol in %q", term)
	}
	if _, _, err := ParseFormula(symbol); err != nil {
		return Term{}, fmt.Errorf("invalid species %q: %w", symbol, err)
	}
	return Term{Symbol: symbol, Coeff: coeff, Role: role}, nil
}
