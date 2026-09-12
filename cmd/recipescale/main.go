// Command recipescale reads a recipe from a file (or stdin) and prints it
// scaled to a target number of servings.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	recipescale "github.com/kjohnson555/recipe-scaler"
)

func main() {
	servings := flag.Int("servings", 0, "target number of servings")
	path := flag.String("file", "", "path to recipe file (reads stdin if omitted)")
	convert := flag.String("convert", "", `convert ingredients to different units after scaling, as "name=unit,name=unit" (e.g. "flour=g,milk=tbsp")`)
	flag.Parse()

	if *servings <= 0 {
		fmt.Fprintln(os.Stderr, "recipescale: -servings must be a positive integer")
		os.Exit(1)
	}

	in := os.Stdin
	if *path != "" {
		f, err := os.Open(*path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "recipescale: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		in = f
	}

	recipe, err := parseRecipe(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recipescale: %v\n", err)
		os.Exit(1)
	}

	scaled, err := recipescale.Scale(recipe, *servings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recipescale: %v\n", err)
		os.Exit(1)
	}

	targetUnits, err := parseConvertFlag(*convert)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recipescale: %v\n", err)
		os.Exit(1)
	}
	if len(targetUnits) > 0 {
		if err := convertIngredients(&scaled, targetUnits); err != nil {
			fmt.Fprintf(os.Stderr, "recipescale: %v\n", err)
			os.Exit(1)
		}
	}

	printRecipe(os.Stdout, scaled)
}

// parseConvertFlag parses the -convert flag's "name=unit,name=unit" syntax
// into a lookup keyed by lowercased ingredient name.
func parseConvertFlag(s string) (map[string]string, error) {
	if s == "" {
		return nil, nil
	}
	targets := make(map[string]string)
	for _, pair := range strings.Split(s, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		name, unit, ok := strings.Cut(pair, "=")
		if !ok {
			return nil, fmt.Errorf("invalid -convert entry %q, want \"name=unit\"", pair)
		}
		name = strings.TrimSpace(name)
		unit = strings.TrimSpace(unit)
		if name == "" || unit == "" {
			return nil, fmt.Errorf("invalid -convert entry %q, want \"name=unit\"", pair)
		}
		targets[strings.ToLower(name)] = unit
	}
	return targets, nil
}

// convertIngredients converts, in place, every ingredient in r whose name
// (case-insensitive) appears in targetUnits.
func convertIngredients(r *recipescale.Recipe, targetUnits map[string]string) error {
	matched := make(map[string]bool, len(targetUnits))
	for i, ing := range r.Ingredients {
		unit, ok := targetUnits[strings.ToLower(ing.Name)]
		if !ok {
			continue
		}
		matched[strings.ToLower(ing.Name)] = true
		converted, err := recipescale.ConvertIngredient(ing, unit)
		if err != nil {
			return err
		}
		r.Ingredients[i] = converted
	}
	for name := range targetUnits {
		if !matched[name] {
			return fmt.Errorf("-convert: no ingredient named %q in recipe", name)
		}
	}
	return nil
}

// parseRecipe reads the plain-text recipe format documented in the README:
//
//	<name>
//	servings: <n>
//	<quantity> | <unit> | <ingredient name>
//	...
func parseRecipe(r io.Reader) (recipescale.Recipe, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return recipescale.Recipe{}, err
	}
	if len(lines) < 2 {
		return recipescale.Recipe{}, fmt.Errorf("recipe must have a name and a servings line")
	}

	name := lines[0]

	const prefix = "servings:"
	servingsLine := lines[1]
	if !strings.HasPrefix(strings.ToLower(servingsLine), prefix) {
		return recipescale.Recipe{}, fmt.Errorf("expected a %q line, got %q", "servings: N", servingsLine)
	}
	servingsStr := strings.TrimSpace(servingsLine[len(prefix):])
	servings, err := strconv.Atoi(servingsStr)
	if err != nil {
		return recipescale.Recipe{}, fmt.Errorf("invalid servings count %q: %w", servingsStr, err)
	}

	var ingredients []recipescale.Ingredient
	for _, line := range lines[2:] {
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			return recipescale.Recipe{}, fmt.Errorf("invalid ingredient line %q, want \"quantity | unit | name\"", line)
		}
		qty, err := recipescale.ParseQuantity(strings.TrimSpace(parts[0]))
		if err != nil {
			return recipescale.Recipe{}, fmt.Errorf("ingredient %q: %w", line, err)
		}
		ingredients = append(ingredients, recipescale.Ingredient{
			Quantity: qty,
			Unit:     strings.TrimSpace(parts[1]),
			Name:     strings.TrimSpace(parts[2]),
		})
	}

	return recipescale.Recipe{Name: name, Servings: servings, Ingredients: ingredients}, nil
}

func printRecipe(w io.Writer, r recipescale.Recipe) {
	fmt.Fprintf(w, "%s (serves %d)\n", r.Name, r.Servings)
	for _, ing := range r.Ingredients {
		qty := recipescale.FormatQuantity(ing.Quantity)
		if ing.Unit == "" {
			fmt.Fprintf(w, "  %s %s\n", qty, ing.Name)
		} else {
			fmt.Fprintf(w, "  %s %s %s\n", qty, ing.Unit, ing.Name)
		}
	}
}
