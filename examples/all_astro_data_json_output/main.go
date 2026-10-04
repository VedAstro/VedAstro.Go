package main

import (
	"encoding/json"
	"fmt"
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
	"os"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		data := map[string]any{}
		calls := []struct {
			name string
			call func() (any, error)
		}{
			{"planet", func() (any, error) { return r.Client.AllPlanetData(r.Context, vedastro.PlanetSun, r.Birth) }},
			{"house", func() (any, error) { return r.Client.AllHouseData(r.Context, vedastro.House1, r.Birth) }},
			{"zodiac", func() (any, error) { return r.Client.AllZodiacSignData(r.Context, vedastro.ZodiacGemini, r.Birth) }},
		}
		for _, item := range calls {
			value, err := r.Fetch(item.call)
			if err != nil {
				return err
			}
			data[item.name] = value
		}
		file, err := os.OpenFile("astro_data.json", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(data); err != nil {
			return err
		}
		fmt.Println("Saved astro_data.json")
		return nil
	})
}
