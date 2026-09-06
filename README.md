# recipe-scaler

Scale a recipe from one serving count to another without ending up with
"0.3333333 cups of flour". Recipe quantities are almost always written as
whole numbers or simple fractions (1/2, 1 3/4, ...), but multiplying them by
a scaling factor rarely lands on a clean fraction. This library parses those
quantities, scales them, and rounds the result back to the nearest eighth so
the output is something you can actually measure.

## Library

```go
package main

import (
	"fmt"

	recipescale "github.com/kjohnson555/recipe-scaler"
)

func main() {
	pancakes := recipescale.Recipe{
		Name:     "Pancakes",
		Servings: 4,
		Ingredients: []recipescale.Ingredient{
			{Name: "flour", Quantity: 2, Unit: "cup"},
			{Name: "milk", Quantity: 1.5, Unit: "cup"},
			{Name: "salt", Quantity: 0.25, Unit: "tsp"},
			{Name: "eggs", Quantity: 2},
		},
	}

	scaled, err := recipescale.Scale(pancakes, 6)
	if err != nil {
		panic(err)
	}

	for _, ing := range scaled.Ingredients {
		fmt.Println(recipescale.FormatQuantity(ing.Quantity), ing.Unit, ing.Name)
	}
	// 3 cup flour
	// 2 1/4 cup milk
	// 3/8 tsp salt
	// 3 eggs
}
```

`Scale` is a pure function: it takes a `Recipe` and a target serving count
and returns a new `Recipe`, leaving the input untouched. `ParseQuantity` and
`FormatQuantity` are pure as well, so all three are trivial to unit test with
plain input/output pairs.

## CLI

The `recipescale` command reads a recipe in a small pipe-delimited text
format and prints it scaled to `-servings`.

Recipe file format:

```
<name>
servings: <n>
<quantity> | <unit> | <ingredient name>
```

Example `pancakes.recipe`:

```
Pancakes
servings: 4
2 | cup | flour
1 1/2 | cup | milk
1/4 | tsp | salt
2 | | eggs
```

Run it:

```sh
go run ./cmd/recipescale -servings 6 -file pancakes.recipe
```

```
Pancakes (serves 6)
  3 cup flour
  2 1/4 cup milk
  3/8 tsp salt
  3 eggs
```

Omit `-file` to read the recipe from stdin instead.

## Status

Early skeleton: quantity parsing, scaling, and formatting work and are
covered by tests; unit conversion (e.g. tsp -> tbsp -> cup) is not built yet.
