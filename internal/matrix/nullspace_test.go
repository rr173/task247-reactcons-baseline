package matrix

import (
	"math/big"
	"testing"

	"task247-reactcons/internal/model"
)

func TestConservedPoolsMassConservation(t *testing.T) {
	// H2 + 0.5 O2 -> H2O is atom-balanced; total H and total O must be conserved.
	species := []model.SpeciesRef{
		{ID: 1, Symbol: "H2", Composition: map[string]int{"H": 2}, Charge: 0},
		{ID: 2, Symbol: "O2", Composition: map[string]int{"O": 2}, Charge: 0},
		{ID: 3, Symbol: "H2O", Composition: map[string]int{"H": 2, "O": 1}, Charge: 0},
	}
	reactions := []model.ReactionRef{
		{ID: 1, Equation: "H2 + 0.5 O2 -> H2O", Participants: []model.ParticipantRef{
			{SpeciesID: 1, Role: "reactant", Coeff: 1},
			{SpeciesID: 2, Role: "reactant", Coeff: 0.5},
			{SpeciesID: 3, Role: "product", Coeff: 1},
		}},
	}
	pools, err := ConservedPools(species, reactions, map[int64]bool{})
	if err != nil {
		t.Fatalf("ConservedPools error: %v", err)
	}
	labels := map[string]bool{}
	for _, p := range pools {
		labels[p.Label] = true
	}
	if !labels["element:H"] {
		t.Errorf("expected conserved pool element:H, got %v", labels)
	}
	if !labels["element:O"] {
		t.Errorf("expected conserved pool element:O, got %v", labels)
	}
}

func TestNullspaceBasisIsExact(t *testing.T) {
	// Matrix [[2,0],[0,2]] has trivial (zero) nullspace.
	A := [][]big.Rat{
		{*big.NewRat(2, 1), *big.NewRat(0, 1)},
		{*big.NewRat(0, 1), *big.NewRat(2, 1)},
	}
	basis := NullspaceBasis(A)
	if len(basis) != 0 {
		t.Errorf("expected empty nullspace, got %d vectors", len(basis))
	}
}
