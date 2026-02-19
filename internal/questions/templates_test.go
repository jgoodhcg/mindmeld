package questions

import "testing"

func TestGetAvailableTemplatesFiltersByRatingAndUsedIDs(t *testing.T) {
	available := GetAvailableTemplates([]string{"quick-001", "work-001"}, 10)
	if len(available) == 0 {
		t.Fatalf("expected kids-safe templates, got none")
	}

	for _, tpl := range available {
		if tpl.ID == "quick-001" || tpl.ID == "work-001" {
			t.Fatalf("expected used template %s to be excluded", tpl.ID)
		}
		if tpl.MinRating > 10 {
			t.Fatalf("template %s exceeds kids rating", tpl.ID)
		}
	}
}

func TestBuildPackSectionsHonorsPackRating(t *testing.T) {
	kidsSections := BuildPackSections(nil, 10)
	if len(kidsSections) != 2 {
		t.Fatalf("expected 2 kids-safe packs, got %d", len(kidsSections))
	}
	if kidsSections[0].Pack.ID != PackQuickBrain {
		t.Fatalf("expected first kids pack %s, got %s", PackQuickBrain, kidsSections[0].Pack.ID)
	}
	if kidsSections[1].Pack.ID != PackWorldSnapshot {
		t.Fatalf("expected second kids pack %s, got %s", PackWorldSnapshot, kidsSections[1].Pack.ID)
	}

	workSections := BuildPackSections(nil, 20)
	if len(workSections) != len(AllPacks) {
		t.Fatalf("expected all packs for work rating, got %d", len(workSections))
	}
	if workSections[0].Pack.ID != PackWorkEssentials {
		t.Fatalf("expected first work pack %s, got %s", PackWorkEssentials, workSections[0].Pack.ID)
	}
}
