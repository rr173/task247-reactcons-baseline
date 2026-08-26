package diagnose

import (
	"encoding/json"
	"testing"

	"task247-reactcons/internal/model"
)

func TestAnalyzeMarksViolatedMoietyClosure(t *testing.T) {
	species := []model.SpeciesRef{
		{ID: 1, Symbol: "H2", Composition: map[string]int{"H": 2}},
		{ID: 2, Symbol: "H2O", Composition: map[string]int{"H": 2, "O": 1}},
	}
	reactions := []model.ReactionRef{{
		ID:       7,
		Equation: "H2 -> H2O",
		Participants: []model.ParticipantRef{
			{SpeciesID: 1, Role: "reactant", Coeff: 1},
			{SpeciesID: 2, Role: "product", Coeff: 1},
		},
	}}
	payload, err := json.Marshal(model.MoietyClosure{Members: []model.MoietyMember{{Symbol: "H2", Coeff: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	result, err := Analyze(species, reactions, nil, []model.Constraint{{ID: 3, Kind: "moietyclosure", Payload: payload}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != model.NetworkConflict || len(result.ConstraintChecks) != 1 || result.ConstraintChecks[0].Status != "violated" {
		t.Fatalf("unexpected constraint result: %+v", result)
	}
}
