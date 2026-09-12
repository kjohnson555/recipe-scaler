package main

import (
	"reflect"
	"testing"

	recipescale "github.com/kjohnson555/recipe-scaler"
)

func TestParseConvertFlag(t *testing.T) {
	cases := []struct {
		in   string
		want map[string]string
	}{
		{"", nil},
		{"flour=g", map[string]string{"flour": "g"}},
		{"Flour=g, Milk=tbsp", map[string]string{"flour": "g", "milk": "tbsp"}},
	}
	for _, c := range cases {
		got, err := parseConvertFlag(c.in)
		if err != nil {
			t.Errorf("parseConvertFlag(%q) returned error: %v", c.in, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("parseConvertFlag(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseConvertFlagInvalid(t *testing.T) {
	cases := []string{"flour", "flour=", "=g"}
	for _, c := range cases {
		if _, err := parseConvertFlag(c); err == nil {
			t.Errorf("parseConvertFlag(%q) expected an error, got nil", c)
		}
	}
}

func TestParseConvertFlagTrailingComma(t *testing.T) {
	got, err := parseConvertFlag("flour=g,")
	if err != nil {
		t.Fatalf("parseConvertFlag returned error: %v", err)
	}
	want := map[string]string{"flour": "g"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseConvertFlag(%q) = %v, want %v", "flour=g,", got, want)
	}
}

func TestConvertIngredients(t *testing.T) {
	r := recipescale.Recipe{
		Name:     "Pancakes",
		Servings: 6,
		Ingredients: []recipescale.Ingredient{
			{Name: "flour", Quantity: 3, Unit: "cup"},
			{Name: "Milk", Quantity: 2.25, Unit: "cup"},
			{Name: "eggs", Quantity: 3},
		},
	}

	if err := convertIngredients(&r, map[string]string{"milk": "tbsp"}); err != nil {
		t.Fatalf("convertIngredients returned error: %v", err)
	}

	if r.Ingredients[0].Unit != "cup" || r.Ingredients[0].Quantity != 3 {
		t.Errorf("flour changed unexpectedly: %+v", r.Ingredients[0])
	}
	if r.Ingredients[1].Unit != "tbsp" || r.Ingredients[1].Quantity != 36 {
		t.Errorf("milk = %+v, want 36 tbsp", r.Ingredients[1])
	}
	if r.Ingredients[2].Unit != "" || r.Ingredients[2].Quantity != 3 {
		t.Errorf("eggs changed unexpectedly: %+v", r.Ingredients[2])
	}
}

func TestConvertIngredientsUnknownName(t *testing.T) {
	r := recipescale.Recipe{
		Ingredients: []recipescale.Ingredient{{Name: "flour", Quantity: 1, Unit: "cup"}},
	}
	if err := convertIngredients(&r, map[string]string{"sugar": "g"}); err == nil {
		t.Error("convertIngredients with an unknown ingredient name expected an error, got nil")
	}
}

func TestConvertIngredientsIncompatibleUnit(t *testing.T) {
	r := recipescale.Recipe{
		Ingredients: []recipescale.Ingredient{{Name: "flour", Quantity: 1, Unit: "cup"}},
	}
	if err := convertIngredients(&r, map[string]string{"flour": "g"}); err == nil {
		t.Error("convertIngredients with an incompatible unit expected an error, got nil")
	}
}
