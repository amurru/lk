package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/amurru/lk/internal/domain"
)

type Writer struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (w Writer) Write(entries []domain.FileEntry, kind string, format domain.OutputFormat) error {
	switch format {
	case "", domain.OutputFormatTable:
		return writeTable(w.Stdout, entries, kind)
	case domain.OutputFormatSimple:
		return writeSimple(w.Stdout, entries)
	case domain.OutputFormatJSON:
		return writeJSON(w.Stdout, entries)
	default:
		return fmt.Errorf("unknown format %q", format)
	}
}

func writeJSON(w io.Writer, entries []domain.FileEntry) error {
	if entries == nil {
		entries = []domain.FileEntry{}
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	value := float64(bytes)
	unit := ""
	for _, next := range units {
		value /= 1024
		unit = next
		if value < 1024 {
			break
		}
	}
	return fmt.Sprintf("%.1f %s", value, unit)
}

func formatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04")
}

func describeEmpty(kind string) string {
	if kind == "" {
		return "No files found."
	}
	return fmt.Sprintf("No %s files found.", kind)
}
