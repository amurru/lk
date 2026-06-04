package kinds

import (
	"fmt"
	"sort"
	"strings"
)

type Kind struct {
	Name       string
	Aliases    []string
	Extensions []string
	MIMEs      []string
}

type Registry struct {
	kinds        []Kind
	nameIndex    map[string]Kind
	aliasIndex   map[string]string
	extensionMap map[string]string
	mimeMap      map[string]string
}

func DefaultRegistry() Registry {
	kinds := []Kind{
		{Name: "documents", Aliases: []string{"doc", "docs"}, Extensions: []string{"csv", "doc", "docm", "docx", "dot", "dotm", "dotx", "epub", "log", "md", "odp", "ods", "odt", "pdf", "pot", "potm", "potx", "pps", "ppsm", "ppsx", "ppt", "pptm", "pptx", "rtf", "txt", "xlam", "xls", "xlsb", "xlsm", "xlsx", "xlt", "xltm", "xltx", "xps"}, MIMEs: []string{"application/pdf", "application/msword", "application/epub+zip", "text/plain", "text/markdown", "text/csv", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.openxmlformats-officedocument.spreadsheetml.template", "application/vnd.ms-excel.sheet.macroEnabled.12", "application/vnd.ms-excel.sheet.binary.macroEnabled.12", "application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation", "application/vnd.openxmlformats-officedocument.presentationml.template", "application/vnd.openxmlformats-officedocument.presentationml.slideshow", "application/vnd.ms-powerpoint.presentation.macroEnabled.12", "application/vnd.ms-powerpoint.slideshow.macroEnabled.12", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.openxmlformats-officedocument.wordprocessingml.template", "application/vnd.ms-word.document.macroEnabled.12", "application/vnd.oasis.opendocument.text", "application/vnd.oasis.opendocument.spreadsheet", "application/vnd.oasis.opendocument.presentation", "application/vnd.ms-xpsdocument"}},
		{Name: "images", Aliases: []string{"image", "pics", "photo", "photos"}, Extensions: []string{"jpg", "jpeg", "png", "gif", "svg", "webp", "bmp", "tif", "tiff"}, MIMEs: []string{"image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/webp", "image/bmp", "image/tiff"}},
		{Name: "videos", Aliases: []string{"vid", "vids", "movie", "movies"}, Extensions: []string{"mp4", "avi", "mkv", "mov", "webm", "m4v", "mpeg", "mpg"}, MIMEs: []string{"video/mp4", "video/x-matroska", "video/webm", "video/quicktime"}},
		{Name: "archives", Aliases: []string{"archive", "zips"}, Extensions: []string{"zip", "tar", "gz", "rar", "7z", "tgz", "bz2"}, MIMEs: []string{"application/zip", "application/x-tar", "application/x-rar-compressed", "application/gzip", "application/x-7z-compressed"}},
		{Name: "audio", Aliases: []string{"sound", "music"}, Extensions: []string{"aac", "ac3", "aif", "aiff", "amr", "au", "flac", "m4a", "mid", "midi", "mka", "mp3", "mpc", "ogg", "opus", "wav", "wma", "wv"}, MIMEs: []string{"audio/mpeg", "audio/wav", "audio/x-wav", "audio/ogg", "audio/flac", "audio/mp4", "audio/aac", "audio/x-ms-wma", "audio/opus", "audio/aiff", "audio/x-aiff", "audio/basic", "audio/webm", "audio/ac3", "audio/amr", "audio/midi"}},
		{Name: "code", Aliases: []string{"src", "source", "script"}, Extensions: []string{"go", "py", "js", "sh", "html", "css", "ts", "tsx", "json", "yaml", "yml", "rb", "java", "c", "h", "cpp", "hpp"}, MIMEs: []string{"text/x-go", "text/x-python", "text/javascript", "application/javascript", "text/html", "text/css", "application/json", "application/x-yaml", "text/x-shellscript", "application/xml"}},
	}
	return NewRegistry(kinds)
}

func NewRegistry(entries []Kind) Registry {
	r := Registry{
		kinds:        make([]Kind, 0, len(entries)),
		nameIndex:    map[string]Kind{},
		aliasIndex:   map[string]string{},
		extensionMap: map[string]string{},
		mimeMap:      map[string]string{},
	}
	for _, kind := range entries {
		normalized := Kind{
			Name:       strings.ToLower(strings.TrimSpace(kind.Name)),
			Aliases:    normalizeStrings(kind.Aliases),
			Extensions: normalizeExtensions(kind.Extensions),
			MIMEs:      normalizeStrings(kind.MIMEs),
		}
		r.kinds = append(r.kinds, normalized)
		r.nameIndex[normalized.Name] = normalized
		for _, alias := range normalized.Aliases {
			r.aliasIndex[alias] = normalized.Name
		}
		for _, ext := range normalized.Extensions {
			r.extensionMap[ext] = normalized.Name
		}
		for _, mime := range normalized.MIMEs {
			r.mimeMap[mime] = normalized.Name
		}
	}
	sort.Slice(r.kinds, func(i, j int) bool { return r.kinds[i].Name < r.kinds[j].Name })
	return r
}

func (r Registry) Names() []string {
	out := make([]string, 0, len(r.kinds))
	for _, kind := range r.kinds {
		out = append(out, kind.Name)
	}
	return out
}

func (r Registry) Kinds() []Kind {
	out := make([]Kind, len(r.kinds))
	copy(out, r.kinds)
	return out
}

func (r Registry) CanonicalName(name string) (string, bool) {
	n := strings.ToLower(strings.TrimSpace(name))
	if n == "" {
		return "", false
	}
	if _, ok := r.nameIndex[n]; ok {
		return n, true
	}
	if canonical, ok := r.aliasIndex[n]; ok {
		return canonical, true
	}
	return "", false
}

func (r Registry) KindForExtension(ext string) (string, bool) {
	ext = NormalizeExtension(ext)
	kind, ok := r.extensionMap[ext]
	return kind, ok
}

func (r Registry) KindForMIME(mime string) (string, bool) {
	mime = strings.ToLower(strings.TrimSpace(mime))
	if mime == "" {
		return "", false
	}
	if kind, ok := r.mimeMap[mime]; ok {
		return kind, true
	}
	if strings.HasPrefix(mime, "text/") {
		return "documents", true
	}
	if strings.HasPrefix(mime, "audio/") {
		return "audio", true
	}
	if strings.HasPrefix(mime, "video/") {
		return "videos", true
	}
	if strings.HasPrefix(mime, "image/") {
		return "images", true
	}
	return "", false
}

func (r Registry) KindByName(name string) (Kind, bool) {
	n, ok := r.CanonicalName(name)
	if !ok {
		return Kind{}, false
	}
	kind, ok := r.nameIndex[n]
	return kind, ok
}

func (r Registry) Describe() string {
	lines := make([]string, 0, len(r.kinds))
	for _, kind := range r.kinds {
		parts := []string{kind.Name}
		if len(kind.Aliases) > 0 {
			parts = append(parts, fmt.Sprintf("aliases=%s", strings.Join(kind.Aliases, ",")))
		}
		if len(kind.Extensions) > 0 {
			parts = append(parts, fmt.Sprintf("extensions=%s", strings.Join(kind.Extensions, ",")))
		}
		lines = append(lines, strings.Join(parts, " "))
	}
	return strings.Join(lines, "\n")
}

func normalizeStrings(in []string) []string {
	out := make([]string, 0, len(in))
	seen := map[string]struct{}{}
	for _, s := range in {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func normalizeExtensions(in []string) []string {
	out := make([]string, 0, len(in))
	seen := map[string]struct{}{}
	for _, s := range in {
		s = NormalizeExtension(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func NormalizeExtension(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	ext = strings.TrimPrefix(ext, ".")
	return ext
}
