package domain

import "time"

type MatchSource string

const (
	MatchSourceExtension MatchSource = "extension"
	MatchSourceMagic     MatchSource = "magic"
	MatchSourceUnknown   MatchSource = "unknown"
)

type SortBy string

const (
	SortByName     SortBy = "name"
	SortBySize     SortBy = "size"
	SortByModified SortBy = "modified"
)

type OutputFormat string

const (
	OutputFormatTable  OutputFormat = "table"
	OutputFormatSimple OutputFormat = "simple"
	OutputFormatJSON   OutputFormat = "json"
	OutputFormatXML   OutputFormat = "xml"
)

type FileEntry struct {
	Path       string      `json:"path"`
	Name       string      `json:"name"`
	Extension  string      `json:"extension,omitempty"`
	SizeBytes  int64       `json:"size_bytes"`
	ModifiedAt time.Time   `json:"modified_at"`
	Kind       string      `json:"kind"`
	MatchedBy  MatchSource `json:"matched_by"`
}

type ScanOptions struct {
	Recursive     bool
	IncludeHidden bool
	UseMagic      bool
	SortBy        SortBy
	Limit         int
}
