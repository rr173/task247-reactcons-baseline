package chem

import "testing"

func TestParseFormula(t *testing.T) {
	cases := []struct {
		in        string
		wantComp  Composition
		wantCharge int
	}{
		{"H2O", Composition{"H": 2, "O": 1}, 0},
		{"CO2", Composition{"C": 1, "O": 2}, 0},
		{"SO4^2-", Composition{"S": 1, "O": 4}, -2},
		{"Fe3+", Composition{"Fe": 1}, 3},
		{"Ca(OH)2", Composition{"Ca": 1, "O": 2, "H": 2}, 0},
		{"(NH4)2SO4", Composition{"N": 2, "H": 8, "S": 1, "O": 4}, 0},
		{"PO4^3-", Composition{"P": 1, "O": 4}, -3},
		{"H+", Composition{"H": 1}, 1},
	}
	for _, c := range cases {
		comp, charge, err := ParseFormula(c.in)
		if err != nil {
			t.Fatalf("ParseFormula(%q) error: %v", c.in, err)
		}
		if charge != c.wantCharge {
			t.Errorf("ParseFormula(%q) charge = %d, want %d", c.in, charge, c.wantCharge)
		}
		for k, v := range c.wantComp {
			if comp[k] != v {
				t.Errorf("ParseFormula(%q) %s = %d, want %d", c.in, k, comp[k], v)
			}
		}
		for k, v := range comp {
			if c.wantComp[k] != v {
				t.Errorf("ParseFormula(%q) unexpected element %s=%d", c.in, k, v)
			}
		}
	}
}

func TestParseFormulaInvalid(t *testing.T) {
	for _, in := range []string{"", "Xx", "H2O)", "(CO3"} {
		if _, _, err := ParseFormula(in); err == nil {
			t.Errorf("ParseFormula(%q) expected error", in)
		}
	}
}

func TestParseEquation(t *testing.T) {
	reactants, products, rev, err := ParseEquation("2 H2 + O2 -> 2 H2O")
	if err != nil {
		t.Fatalf("ParseEquation error: %v", err)
	}
	if rev {
		t.Errorf("expected irreversible")
	}
	if len(reactants) != 2 || len(products) != 1 {
		t.Fatalf("unexpected term counts r=%d p=%d", len(reactants), len(products))
	}
	if reactants[0].Symbol != "H2" || reactants[0].Coeff != 2 || reactants[0].Role != "reactant" {
		t.Errorf("reactant[0] = %+v", reactants[0])
	}
	if products[0].Symbol != "H2O" || products[0].Coeff != 2 || products[0].Role != "product" {
		t.Errorf("product[0] = %+v", products[0])
	}

	if _, _, _, err := ParseEquation("A + B C"); err == nil {
		t.Errorf("expected missing-arrow error")
	}
}
