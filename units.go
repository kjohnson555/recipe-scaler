package recipescale

import (
	"fmt"
	"strings"
)

// unitKind groups units that can be converted into one another. Volume and
// mass are never compatible with each other since the conversion depends on
// the ingredient's density.
type unitKind int

const (
	unitVolume unitKind = iota
	unitMass
)

// unitDef describes a unit as a multiplier against its kind's base unit:
// teaspoon for volume, gram for mass. Keeping everything in base units
// means converting any unit to any other of the same kind is one
// multiply and one divide, with no per-pair table to maintain.
type unitDef struct {
	kind   unitKind
	toBase float64
}

var unitDefs = map[string]unitDef{
	"tsp":  {unitVolume, 1},
	"tbsp": {unitVolume, 3},
	"cup":  {unitVolume, 48},
	"g":    {unitMass, 1},
	"kg":   {unitMass, 1000},
}

// unitAliases maps the spellings a recipe file is likely to use onto the
// canonical keys in unitDefs.
var unitAliases = map[string]string{
	"tsp": "tsp", "tsps": "tsp", "teaspoon": "tsp", "teaspoons": "tsp",
	"tbsp": "tbsp", "tbsps": "tbsp", "tablespoon": "tbsp", "tablespoons": "tbsp",
	"cup": "cup", "cups": "cup",
	"g": "g", "gram": "g", "grams": "g",
	"kg": "kg", "kilogram": "kg", "kilograms": "kg",
}

func canonicalUnit(u string) (string, bool) {
	key := strings.ToLower(strings.TrimSpace(u))
	canon, ok := unitAliases[key]
	return canon, ok
}

// ConvertQuantity converts qty from one unit to another. The two units must
// be of the same kind (both volume or both mass); converting between volume
// and mass would require the ingredient's density, which this library does
// not track.
func ConvertQuantity(qty float64, from, to string) (float64, error) {
	fromCanon, ok := canonicalUnit(from)
	if !ok {
		return 0, fmt.Errorf("unknown unit %q", from)
	}
	toCanon, ok := canonicalUnit(to)
	if !ok {
		return 0, fmt.Errorf("unknown unit %q", to)
	}

	fromDef := unitDefs[fromCanon]
	toDef := unitDefs[toCanon]
	if fromDef.kind != toDef.kind {
		return 0, fmt.Errorf("cannot convert %q to %q: incompatible units", from, to)
	}

	return qty * fromDef.toBase / toDef.toBase, nil
}

// ConvertIngredient returns a copy of ing with its quantity converted to
// targetUnit. ing is not modified.
func ConvertIngredient(ing Ingredient, targetUnit string) (Ingredient, error) {
	if ing.Unit == "" {
		return Ingredient{}, fmt.Errorf("ingredient %q has no unit to convert", ing.Name)
	}
	qty, err := ConvertQuantity(ing.Quantity, ing.Unit, targetUnit)
	if err != nil {
		return Ingredient{}, fmt.Errorf("ingredient %q: %w", ing.Name, err)
	}
	return Ingredient{Name: ing.Name, Quantity: qty, Unit: targetUnit}, nil
}
