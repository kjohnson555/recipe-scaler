// Package recipescale scales recipe ingredient quantities to a different
// number of servings.
package recipescale

// Ingredient is one line of a recipe: how much of something, in what unit.
// Unit is free text ("cup", "tsp", "g") and may be empty for countable
// items like "2 eggs".
type Ingredient struct {
	Name     string
	Quantity float64
	Unit     string
}

// Recipe is a named list of ingredients sized for Servings people.
type Recipe struct {
	Name        string
	Servings    int
	Ingredients []Ingredient
}
