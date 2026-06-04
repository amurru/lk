package output

import (
	"fmt"
	"io"

	"github.com/amurru/lk/internal/domain"
)

func writeSimple(w io.Writer, entries []domain.FileEntry) error {
	for _, entry := range entries {
		if _, err := fmt.Fprintln(w, entry.Path); err != nil {
			return err
		}
	}
	return nil
}
