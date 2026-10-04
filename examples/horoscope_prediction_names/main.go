package main

import (
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Horoscope predictions", func() (any, error) { return r.Client.HoroscopePredictions(r.Context, r.Birth) }); err != nil {
			return err
		}
		return nil
	})
}
