package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		start, err := vedastro.NewTime("00:00 01/01/2020 +05:30", r.Birth.Location)
		if err != nil {
			return err
		}
		end, err := vedastro.NewTime("23:59 31/12/2020 +05:30", r.Birth.Location)
		if err != nil {
			return err
		}
		return r.Show("Vimshottari dasa", func() (any, error) {
			return r.Client.DasaAtRange(r.Context, r.Birth, start, end, vedastro.DasaAtRangeOptions{
				Levels: vedastro.Ptr(3), PrecisionHours: vedastro.Ptr(100),
			})
		})
	})
}
