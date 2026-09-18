// Package recipescale scales recipe ingredient quantities to a different
// number of servings.
package recipescale

// Ingredient is one line of a recipe: how much of something, in what unit.
// Unit is free text ("cup", "tsp", "g") and may be empty for countable
// items like "2 eggs".
//
// MaxQuantity is nonzero for a range such as "1-2 cloves garlic", where
// Quantity holds the low end and MaxQuantity the high end. It is zero for
// an ordinary single-value quantity.
type Ingredient struct {
	Name        string
	Quantity    float64
	MaxQuantity float64
	Unit        string
}

// IsRange reports whether ing was written as a quantity range.
func (ing Ingredient) IsRange() bool {
	return ing.MaxQuantity != 0
}

// Recipe is a named list of ingredients sized for Servings people.
type Recipe struct {
	Name        string
	Servings    int
	Ingredients []Ingredient
}
