package matrix

import (
	"fmt"
	"math/big"
	"sort"

	"task247-reactcons/internal/model"
)

const conservationTol = 1e-9

type poolVec struct {
	vec   []big.Rat
	label string
}

// NullspaceBasis returns a basis of the nullspace of A (rows x cols) using exact
// rational Gaussian elimination. Each returned vector has length cols.
func NullspaceBasis(A [][]big.Rat) [][]big.Rat {
	rows := len(A)
	if rows == 0 {
		return nil
	}
	cols := len(A[0])
	M := make([][]big.Rat, rows)
	for i := range A {
		M[i] = make([]big.Rat, cols)
		for j := range A[i] {
			M[i][j].Set(&A[i][j])
		}
	}
	r := 0
	for c := 0; c < cols && r < rows; c++ {
		pivot := -1
		for i := r; i < rows; i++ {
			if M[i][c].RatString() != "0" {
				pivot = i
				break
			}
		}
		if pivot == -1 {
			continue
		}
		M[r], M[pivot] = M[pivot], M[r]
		inv := new(big.Rat).Inv(&M[r][c])
		for j := c; j < cols; j++ {
			M[r][j].Mul(&M[r][j], inv)
		}
		for i := 0; i < rows; i++ {
			if i == r {
				continue
			}
			if M[i][c].RatString() == "0" {
				continue
			}
			f := new(big.Rat).Set(&M[i][c])
			for j := c; j < cols; j++ {
				t := new(big.Rat).Mul(&M[r][j], f)
				M[i][j].Sub(&M[i][j], t)
			}
		}
		r++
	}
	pivots := make([]int, r)
	for i := 0; i < r; i++ {
		for j := 0; j < cols; j++ {
			if M[i][j].RatString() != "0" {
				pivots[i] = j
				break
			}
		}
	}
	pivotSet := make(map[int]bool, r)
	for i := 0; i < r; i++ {
		pivotSet[pivots[i]] = true
	}
	freeCols := []int{}
	for j := 0; j < cols; j++ {
		if !pivotSet[j] {
			freeCols = append(freeCols, j)
		}
	}
	basis := make([][]big.Rat, 0, len(freeCols))
	for _, f := range freeCols {
		vec := make([]big.Rat, cols)
		vec[f].SetInt64(1)
		for i := 0; i < r; i++ {
			vec[pivots[i]].Set(&M[i][f])
		}
		basis = append(basis, vec)
	}
	return basis
}

func ratFromFloat(f float64) *big.Rat {
	num := int64(f * 1e9)
	return new(big.Rat).SetFrac(big.NewInt(num), big.NewInt(1e9))
}

// ConservedPools computes the conserved subspace (left nullspace of the
// stoichiometric matrix) over non-open species. Element and charge pools are
// reported by their physical labels when they are conserved; additional linearly
// independent vectors complete the basis as linear-N pools.
func ConservedPools(species []model.SpeciesRef, reactions []model.ReactionRef, exempt map[int64]bool) ([]model.ConservedPool, error) {
	var idx []int
	var spIDs []int64
	for i, sp := range species {
		if !sp.Open {
			idx = append(idx, i)
			spIDs = append(spIDs, sp.ID)
		}
	}
	m := len(idx)
	n := 0
	for _, rx := range reactions {
		if !exempt[rx.ID] {
			n++
		}
	}
	if m == 0 || n == 0 {
		return nil, nil
	}
	S, _, ids := BuildStoich(species, reactions)
	byID := make(map[int64]int, len(ids))
	for j, rid := range ids {
		byID[rid] = j
	}
	A := make([][]big.Rat, 0, n)
	for _, rx := range reactions {
		if exempt[rx.ID] {
			continue
		}
		j := byID[rx.ID]
		row := make([]big.Rat, m)
		for c, si := range idx {
			row[c] = *ratFromFloat(S[si][j])
		}
		A = append(A, row)
	}

	pools := []poolVec{}

	// Reference vectors: charge, then each element.
	refs := []struct {
		label string
		vec   []big.Rat
	}{}
	chargeVec := make([]big.Rat, m)
	for c, si := range idx {
		chargeVec[c].SetInt64(int64(species[si].Charge))
	}
	refs = append(refs, struct {
		label string
		vec   []big.Rat
	}{"charge", chargeVec})
	for _, e := range elementList(species) {
		e := e
		ev := make([]big.Rat, m)
		for c, si := range idx {
			ev[c].SetInt64(int64(species[si].Composition[e]))
		}
		refs = append(refs, struct {
			label string
			vec   []big.Rat
		}{"element:" + e, ev})
	}
	for _, r := range refs {
		if isZeroVec(r.vec) {
			continue
		}
		if !inNullspace(A, r.vec) {
			continue
		}
		pools = append(pools, poolVec{vec: r.vec, label: r.label})
	}

	// Complete the basis with any linearly independent nullspace vectors.
	for _, b := range NullspaceBasis(A) {
		if !independent(pools, b) {
			continue
		}
		label := ""
		for _, r := range refs {
			if proportional(b, r.vec) {
				label = r.label
				break
			}
		}
		if label == "" {
			label = fmt.Sprintf("linear-%d", len(pools)+1)
		}
		pools = append(pools, poolVec{vec: b, label: label})
	}

	out := make([]model.ConservedPool, 0, len(pools))
	for _, p := range pools {
		members := make([]model.PoolMember, 0, m)
		for c, si := range idx {
			if p.vec[c].RatString() == "0" {
				continue
			}
			members = append(members, model.PoolMember{
				SpeciesID: spIDs[c],
				Symbol:    species[si].Symbol,
				Coeff:     p.vec[c].RatString(),
			})
		}
		if len(members) == 0 {
			continue
		}
		out = append(out, model.ConservedPool{Label: p.label, Members: members})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })
	return out, nil
}

func isZeroVec(v []big.Rat) bool {
	for _, x := range v {
		if x.RatString() != "0" {
			return false
		}
	}
	return true
}

func inNullspace(A [][]big.Rat, v []big.Rat) bool {
	for _, row := range A {
		s := new(big.Rat)
		for c := range v {
			t := new(big.Rat).Mul(&row[c], &v[c])
			s.Add(s, t)
		}
		if s.RatString() != "0" {
			return false
		}
	}
	return true
}

func independent(pools []poolVec, b []big.Rat) bool {
	rows := make([][]big.Rat, 0, len(pools)+1)
	for _, p := range pools {
		rows = append(rows, cloneVec(p.vec))
	}
	rank0 := rankOf(rows)
	rows = append(rows, cloneVec(b))
	return rankOf(rows) > rank0
}

func cloneVec(v []big.Rat) []big.Rat {
	out := make([]big.Rat, len(v))
	for i := range v {
		out[i].Set(&v[i])
	}
	return out
}

func rankOf(rows [][]big.Rat) int {
	if len(rows) == 0 {
		return 0
	}
	cols := len(rows[0])
	M := make([][]big.Rat, len(rows))
	for i := range rows {
		M[i] = make([]big.Rat, cols)
		for j := range rows[i] {
			M[i][j].Set(&rows[i][j])
		}
	}
	r := 0
	for c := 0; c < cols && r < len(M); c++ {
		pivot := -1
		for i := r; i < len(M); i++ {
			if M[i][c].RatString() != "0" {
				pivot = i
				break
			}
		}
		if pivot == -1 {
			continue
		}
		M[r], M[pivot] = M[pivot], M[r]
		inv := new(big.Rat).Inv(&M[r][c])
		for j := c; j < cols; j++ {
			M[r][j].Mul(&M[r][j], inv)
		}
		for i := 0; i < len(M); i++ {
			if i == r {
				continue
			}
			if M[i][c].RatString() == "0" {
				continue
			}
			f := new(big.Rat).Set(&M[i][c])
			for j := c; j < cols; j++ {
				t := new(big.Rat).Mul(&M[r][j], f)
				M[i][j].Sub(&M[i][j], t)
			}
		}
		r++
	}
	return r
}

func proportional(a, b []big.Rat) bool {
	var scale *big.Rat
	for c := range a {
		az, bz := a[c].RatString() == "0", b[c].RatString() == "0"
		if az && bz {
			continue
		}
		if az || bz {
			return false
		}
		k := new(big.Rat).Quo(&a[c], &b[c])
		if scale == nil {
			scale = k
		} else if scale.Cmp(k) != 0 {
			return false
		}
	}
	return scale != nil
}
