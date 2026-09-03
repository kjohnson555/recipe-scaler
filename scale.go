package recipescale

import "fmt"

// Scale returns a new Recipe with every ingredient quantity multiplied by
// targetServings/r.Servings. It does not modify r.
func Scale(r Recipe, targetServings int) (Recipe, error) {
	if targetServings <= 0 {
		return Recipe{}, fmt.Errorf("target servings must be positive, got %d", targetServings)
	}
	if r.Servings <= 0 {
		return Recipe{}, fmt.Errorf("recipe servings must be positive, got %d", r.Servings)
	}

	factor := float64(targetServings) / float64(r.Servings)

	scaled := Recipe{
		Name:        r.Name,
		Servings:    targetServings,
		Ingredients: make([]Ingredient, len(r.Ingredients)),
	}
	for i, ing := range r.Ingredients {
		scaled.Ingredients[i] = Ingredient{
			Name:     ing.Name,
			Quantity: ing.Quantity * factor,
			Unit:     ing.Unit,
		}
	}
	return scaled, nil
}
