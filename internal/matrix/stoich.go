package matrix

import (
	"math"
	"sort"

	"task247-reactcons/internal/model"
)

// BuildStoich builds the species x reactions stoichiometric matrix S where
// S[i][j] = (sum of product coefficients - sum of reactant coefficients) of
// species i in reaction j. It returns S, a speciesID->row index map, and the
// ordered reaction IDs.
func BuildStoich(species []model.SpeciesRef, reactions []model.ReactionRef) ([][]float64, map[int64]int, []int64) {
	rowOf := make(map[int64]int, len(species))
	for i, sp := range species {
		rowOf[sp.ID] = i
	}
	S := make([][]float64, len(species))
	for i := range S {
		S[i] = make([]float64, len(reactions))
	}
	for j, rx := range reactions {
		for _, p := range rx.Participants {
			i, ok := rowOf[p.SpeciesID]
			if !ok {
				continue
			}
			if p.Role == "product" {
				S[i][j] += p.Coeff
			} else {
				S[i][j] -= p.Coeff
			}
		}
	}
	ids := make([]int64, len(reactions))
	for j, rx := range reactions {
		ids[j] = rx.ID
	}
	return S, rowOf, ids
}

func elementList(species []model.SpeciesRef) []string {
	set := map[string]struct{}{}
	for _, sp := range species {
		for e := range sp.Composition {
			set[e] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for e := range set {
		out = append(out, e)
	}
	sort.Strings(out)
	return out
}

// ReactionsConservation evaluates element and charge balance for every
// non-exempt reaction. A reaction is conserved only when all element counts and
// the total charge are unchanged across the reaction.
func ReactionsConservation(species []model.SpeciesRef, reactions []model.ReactionRef, exempt map[int64]bool) []model.ReactionConservation {
	S, _, ids := BuildStoich(species, reactions)
	byID := make(map[int64]int, len(ids))
	for j, rid := range ids {
		byID[rid] = j
	}
	elems := elementList(species)
	results := make([]model.ReactionConservation, 0, len(reactions))
	for _, rx := range reactions {
		if exempt[rx.ID] {
			continue
		}
		j := byID[rx.ID]
		rc := model.ReactionConservation{ReactionID: rx.ID, Equation: rx.Equation, Conserved: true}
		for _, e := range elems {
			v := make([]float64, len(species))
			for i, sp := range species {
				if !sp.Open {
					v[i] = float64(sp.Composition[e])
				}
			}
			net := 0.0
			for i := range species {
				net += S[i][j] * v[i]
			}
			if math.Abs(net) > 1e-9 {
				rc.Conserved = false
				rc.Violations = append(rc.Violations, model.Violation{Kind: "element", Target: e, Net: net})
			}
		}
		v := make([]float64, len(species))
		for i, sp := range species {
			if !sp.Open {
				v[i] = float64(sp.Charge)
			}
		}
		net := 0.0
		for i := range species {
			net += S[i][j] * v[i]
		}
		if math.Abs(net) > 1e-9 {
			rc.Conserved = false
			rc.Violations = append(rc.Violations, model.Violation{Kind: "charge", Target: "charge", Net: net})
		}
		results = append(results, rc)
	}
	return results
}

// ColumnNet returns, for each non-exempt reaction, the dot product of the
// stoichiometric column with the supplied coefficient vector v (length = len(species)).
func ColumnNet(species []model.SpeciesRef, reactions []model.ReactionRef, exempt map[int64]bool, v []float64) map[int64]float64 {
	S, _, ids := BuildStoich(species, reactions)
	byID := make(map[int64]int, len(ids))
	for j, rid := range ids {
		byID[rid] = j
	}
	out := make(map[int64]float64)
	for _, rx := range reactions {
		if exempt[rx.ID] {
			continue
		}
		j := byID[rx.ID]
		net := 0.0
		for i := range species {
			net += S[i][j] * v[i]
		}
		out[rx.ID] = net
	}
	return out
}
