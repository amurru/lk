package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/amurru/lk/internal/domain"
)

func writeTable(w io.Writer, entries []domain.FileEntry, kind string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if len(entries) == 0 {
		_, err := fmt.Fprintln(tw, describeEmpty(kind))
		if err != nil {
			return err
		}
		return tw.Flush()
	}
	_, err := fmt.Fprintln(tw, "kind\tmatched_by\tsize\tmodified\tname")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		_, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", entry.Kind, entry.MatchedBy, formatSize(entry.SizeBytes), formatTime(entry.ModifiedAt), entry.Path)
		if err != nil {
			return err
		}
	}
	return tw.Flush()
}
