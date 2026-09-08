package recipescale

import "testing"

func TestConvertQuantity(t *testing.T) {
	cases := []struct {
		qty      float64
		from, to string
		want     float64
	}{
		{3, "tsp", "tbsp", 1},
		{1, "tbsp", "tsp", 3},
		{1, "cup", "tbsp", 16},
		{48, "tsp", "cup", 1},
		{1, "kg", "g", 1000},
		{500, "g", "kg", 0.5},
		{2, "cups", "cup", 2},
		{1, "tablespoon", "teaspoons", 3},
		{1, "cup", "cup", 1},
	}
	for _, c := range cases {
		got, err := ConvertQuantity(c.qty, c.from, c.to)
		if err != nil {
			t.Errorf("ConvertQuantity(%v, %q, %q) returned error: %v", c.qty, c.from, c.to, err)
			continue
		}
		if got != c.want {
			t.Errorf("ConvertQuantity(%v, %q, %q) = %v, want %v", c.qty, c.from, c.to, got, c.want)
		}
	}
}

func TestConvertQuantityIncompatible(t *testing.T) {
	if _, err := ConvertQuantity(1, "cup", "g"); err == nil {
		t.Error("ConvertQuantity(cup -> g) expected an error, got nil")
	}
}

func TestConvertQuantityUnknownUnit(t *testing.T) {
	cases := []struct{ from, to string }{
		{"cup", "smidgen"},
		{"smidgen", "cup"},
	}
	for _, c := range cases {
		if _, err := ConvertQuantity(1, c.from, c.to); err == nil {
			t.Errorf("ConvertQuantity(1, %q, %q) expected an error, got nil", c.from, c.to)
		}
	}
}

func TestConvertIngredient(t *testing.T) {
	ing := Ingredient{Name: "milk", Quantity: 1, Unit: "cup"}
	got, err := ConvertIngredient(ing, "tbsp")
	if err != nil {
		t.Fatalf("ConvertIngredient returned error: %v", err)
	}
	if got.Quantity != 16 {
		t.Errorf("Quantity = %v, want 16", got.Quantity)
	}
	if got.Unit != "tbsp" {
		t.Errorf("Unit = %q, want %q", got.Unit, "tbsp")
	}
	if got.Name != ing.Name {
		t.Errorf("Name = %q, want %q", got.Name, ing.Name)
	}
	if ing.Unit != "cup" {
		t.Errorf("ConvertIngredient modified the input: Unit = %q, want %q", ing.Unit, "cup")
	}
}

func TestConvertIngredientNoUnit(t *testing.T) {
	ing := Ingredient{Name: "eggs", Quantity: 2}
	if _, err := ConvertIngredient(ing, "g"); err == nil {
		t.Error("ConvertIngredient with no unit expected an error, got nil")
	}
}
