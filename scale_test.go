package recipescale

import "testing"

func testRecipe() Recipe {
	return Recipe{
		Name:     "Pancakes",
		Servings: 4,
		Ingredients: []Ingredient{
			{Name: "flour", Quantity: 2, Unit: "cup"},
			{Name: "milk", Quantity: 1.5, Unit: "cup"},
			{Name: "salt", Quantity: 0.25, Unit: "tsp"},
			{Name: "eggs", Quantity: 2},
		},
	}
}

func TestScaleDoublesQuantities(t *testing.T) {
	r := testRecipe()
	scaled, err := Scale(r, 8)
	if err != nil {
		t.Fatalf("Scale returned error: %v", err)
	}

	if scaled.Servings != 8 {
		t.Errorf("Servings = %d, want 8", scaled.Servings)
	}
	if scaled.Name != r.Name {
		t.Errorf("Name = %q, want %q", scaled.Name, r.Name)
	}

	want := []float64{4, 3, 0.5, 4}
	for i, ing := range scaled.Ingredients {
		if ing.Quantity != want[i] {
			t.Errorf("Ingredients[%d].Quantity = %v, want %v", i, ing.Quantity, want[i])
		}
		if ing.Name != r.Ingredients[i].Name {
			t.Errorf("Ingredients[%d].Name = %q, want %q", i, ing.Name, r.Ingredients[i].Name)
		}
		if ing.Unit != r.Ingredients[i].Unit {
			t.Errorf("Ingredients[%d].Unit = %q, want %q", i, ing.Unit, r.Ingredients[i].Unit)
		}
	}
}

func TestScaleDoesNotModifyInput(t *testing.T) {
	r := testRecipe()
	original := r.Ingredients[0].Quantity

	if _, err := Scale(r, 8); err != nil {
		t.Fatalf("Scale returned error: %v", err)
	}

	if r.Ingredients[0].Quantity != original {
		t.Errorf("Scale modified the input recipe: got %v, want %v", r.Ingredients[0].Quantity, original)
	}
}

func TestScaleSameServings(t *testing.T) {
	r := testRecipe()
	scaled, err := Scale(r, r.Servings)
	if err != nil {
		t.Fatalf("Scale returned error: %v", err)
	}
	for i, ing := range scaled.Ingredients {
		if ing.Quantity != r.Ingredients[i].Quantity {
			t.Errorf("Ingredients[%d].Quantity = %v, want %v", i, ing.Quantity, r.Ingredients[i].Quantity)
		}
	}
}

func TestScaleInvalidTargetServings(t *testing.T) {
	r := testRecipe()
	for _, target := range []int{0, -1} {
		if _, err := Scale(r, target); err == nil {
			t.Errorf("Scale(r, %d) expected an error, got nil", target)
		}
	}
}

func TestScaleInvalidRecipeServings(t *testing.T) {
	r := testRecipe()
	r.Servings = 0
	if _, err := Scale(r, 4); err == nil {
		t.Error("Scale with zero recipe servings expected an error, got nil")
	}
}

func TestScaleEmptyIngredients(t *testing.T) {
	r := Recipe{Name: "Water", Servings: 1}
	scaled, err := Scale(r, 2)
	if err != nil {
		t.Fatalf("Scale returned error: %v", err)
	}
	if len(scaled.Ingredients) != 0 {
		t.Errorf("Ingredients = %v, want empty", scaled.Ingredients)
	}
}
