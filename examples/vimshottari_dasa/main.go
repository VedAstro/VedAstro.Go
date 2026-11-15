package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		return r.Show("Vimshottari dasa", func() (any, error) {
			return r.Client.DasaForLife(r.Context, r.Birth, vedastro.DasaForLifeOptions{
				System: vedastro.Ptr(vedastro.DasaSystemVimshottariJanma),
				Depth:  vedastro.Ptr(3),
			})
		})
	})
}
