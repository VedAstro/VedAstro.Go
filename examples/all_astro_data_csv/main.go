package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
	"os"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		rows := [][]string{{"planet", "result_json"}}
		for _, planet := range []vedastro.PlanetName{vedastro.PlanetSun, vedastro.PlanetMoon, vedastro.PlanetMars,
			vedastro.PlanetMercury, vedastro.PlanetJupiter, vedastro.PlanetVenus, vedastro.PlanetSaturn, vedastro.PlanetRahu, vedastro.PlanetKetu} {
			value, err := r.Fetch(func() (any, error) { return r.Client.AllPlanetData(r.Context, planet, r.Birth) })
			if err != nil {
				return err
			}
			encoded, err := json.Marshal(value)
			if err != nil {
				return err
			}
			rows = append(rows, []string{string(planet), string(encoded)})
		}
		file, err := os.OpenFile("planet_data.csv", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		if err := writer.WriteAll(rows); err != nil {
			return err
		}
		fmt.Println("Saved planet_data.csv")
		return nil
	})
}
