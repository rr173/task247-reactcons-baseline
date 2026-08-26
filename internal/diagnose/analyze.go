package diagnose

import (
	"encoding/json"
	"math"
	"sort"
	"strings"

	"task247-reactcons/internal/matrix"
	"task247-reactcons/internal/model"
)

const conservationTol = 1e-9

// Analyze evaluates conservation for all non-exempt reactions, derives minimal
// conflict reaction sets, computes conserved pools and checks declared
// constraints. The result status is conflict when any conservation constraint
// is violated, otherwise publishable.
func Analyze(species []model.SpeciesRef, reactions []model.ReactionRef, exempt map[int64]bool, constraints []model.Constraint) (*model.SolveResult, error) {
	checks := matrix.ReactionsConservation(species, reactions, map[int64]bool{})
	pools, err := matrix.ConservedPools(species, reactions, map[int64]bool{})
	if err != nil {
		return nil, err
	}

	grouped := map[string][]int64{}
	for _, rc := range checks {
		for _, v := range rc.Violations {
			key := v.Kind + ":" + v.Target
			grouped[key] = append(grouped[key], rc.ReactionID)
		}
	}
	conflictSets := make([]model.ConflictSet, 0, len(grouped))
	for key, rids := range grouped {
		parts := strings.SplitN(key, ":", 2)
		conflictSets = append(conflictSets, model.ConflictSet{
			Kind:        parts[0],
			Target:      parts[1],
			ReactionIDs: dedupe(rids),
		})
	}
	sort.Slice(conflictSets, func(i, j int) bool {
		if conflictSets[i].Kind != conflictSets[j].Kind {
			return conflictSets[i].Kind < conflictSets[j].Kind
		}
		return conflictSets[i].Target < conflictSets[j].Target
	})

	constraintChecks, err := checkConstraints(species, reactions, exempt, constraints)
	if err != nil {
		return nil, err
	}

	out := &model.SolveResult{
		ReactionChecks:   checks,
		ConservedPools:   pools,
		ConflictSets:     conflictSets,
		ConstraintChecks: constraintChecks,
	}
	if len(conflictSets) > 0 {
		out.Status = model.NetworkConflict
	}
	for _, c := range constraintChecks {
		if c.Status == "violated" {
			out.Status = model.NetworkConflict
		}
	}
	if out.Status != model.NetworkConflict {
		out.Status = model.NetworkPublishable
	}
	for _, cs := range conflictSets {
		out.ConflictReactionIDs = append(out.ConflictReactionIDs, cs.ReactionIDs...)
	}
	out.ConflictReactionIDs = dedupe(out.ConflictReactionIDs)

	return out, nil
}

func dedupe(ids []int64) []int64 {
	seen := map[int64]struct{}{}
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// checkConstraints evaluates each declared moiety-closure constraint: a declared
// linear combination of species amounts must be conserved by the network.
func checkConstraints(species []model.SpeciesRef, reactions []model.ReactionRef, exempt map[int64]bool, constraints []model.Constraint) ([]model.Constraint, error) {
	symbolToIndex := map[string]int{}
	for i, sp := range species {
		symbolToIndex[sp.Symbol] = i
	}
	v := make([]float64, len(species))
	out := make([]model.Constraint, 0, len(constraints))
	for _, c := range constraints {
		if c.Kind != "moietyclosure" {
			continue
		}
		var mc model.MoietyClosure
		if err := json.Unmarshal(c.Payload, &mc); err != nil {
			c.Status = "violated"
			out = append(out, c)
			continue
		}
		for i := range v {
			v[i] = 0
		}
		ok := true
		for _, mem := range mc.Members {
			i, found := symbolToIndex[mem.Symbol]
			if !found {
				ok = false
				break
			}
			v[i] += mem.Coeff
		}
		if !ok {
			c.Status = "violated"
			out = append(out, c)
			continue
		}
		nets := matrix.ColumnNet(species, reactions, exempt, v)
		maxAbs := 0.0
		for _, net := range nets {
			if math.Abs(net) > maxAbs {
				maxAbs = math.Abs(net)
			}
		}
		if maxAbs > conservationTol {
			c.Status = "violated"
		} else {
			c.Status = "satisfied"
		}
		out = append(out, c)
	}
	return out, nil
}
