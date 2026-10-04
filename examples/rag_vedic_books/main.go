package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Available books", func() (any, error) { return r.Client.GetAvailableSourceTexts(r.Context) }); err != nil {
			return err
		}
		if err := r.Show("Saturn in the seventh house", func() (any, error) {
			return r.Client.SearchSourceText(r.Context, "effects of Saturn in the 7th house", vedastro.SearchSourceTextOptions{TopK: vedastro.Ptr(3)})
		}); err != nil {
			return err
		}
		return nil
	})
}
