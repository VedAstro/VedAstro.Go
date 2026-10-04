package main

import (
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Birth number", func() (any, error) { return r.Client.BirthNumber(r.Context, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("Destiny number", func() (any, error) { return r.Client.DestinyNumber(r.Context, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("Name prediction", func() (any, error) { return r.Client.NameNumberPrediction(r.Context, "Sengiv") }); err != nil {
			return err
		}
		return nil
	})
}
