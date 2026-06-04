package output

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/amurru/lk/internal/domain"
)

type Writer struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (w Writer) Write(entries []domain.FileEntry, kind string, format domain.OutputFormat, null bool) error {
	switch format {
	case "", domain.OutputFormatTable:
		if null {
			return writeNull(w.Stdout, entries)
		}
		return writeTable(w.Stdout, entries, kind)
	case domain.OutputFormatSimple:
		if null {
			return writeNull(w.Stdout, entries)
		}
		return writeSimple(w.Stdout, entries)
	case domain.OutputFormatJSON:
		return writeJSON(w.Stdout, entries)
	case domain.OutputFormatXML:
		return writeXML(w.Stdout, entries)
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

func writeXML(w io.Writer, entries []domain.FileEntry) error {
	if entries == nil {
		entries = []domain.FileEntry{}
	}
	type xmlEntry struct {
		XMLName   xml.Name `xml:"file"`
		Path      string   `xml:"path"`
		Name      string   `xml:"name"`
		Extension string   `xml:"extension,omitempty"`
		SizeBytes int64    `xml:"sizeBytes"`
		Modified  string   `xml:"modifiedAt"`
		Kind      string   `xml:"kind"`
		MatchedBy string   `xml:"matchedBy"`
	}
	type xmlRoot struct {
		XMLName xml.Name   `xml:"files"`
		Entries []xmlEntry `xml:"file"`
	}
	xentries := make([]xmlEntry, len(entries))
	for i, e := range entries {
		xentries[i] = xmlEntry{
			Path:      e.Path,
			Name:      e.Name,
			Extension: e.Extension,
			SizeBytes: e.SizeBytes,
			Modified:  e.ModifiedAt.UTC().Format(time.RFC3339),
			Kind:      e.Kind,
			MatchedBy: string(e.MatchedBy),
		}
	}
	data, err := xml.MarshalIndent(xmlRoot{Entries: xentries}, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, xml.Header+string(data))
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
