package classifier

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/amurru/lk/internal/domain"
	"github.com/amurru/lk/internal/kinds"
)

type Classifier struct {
	Registry kinds.Registry
	UseMagic bool
}

type Classification struct {
	Kind         string
	MatchedBy    domain.MatchSource
	Extension    string
	DetectedMIME string
}

func (c Classifier) Classify(path string) (Classification, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Classification{}, err
	}
	if info.IsDir() {
		return Classification{}, fmt.Errorf("%s is a directory", path)
	}

	ext := kinds.NormalizeExtension(filepath.Ext(path))
	if kind, ok := c.Registry.KindForExtension(ext); ok {
		return Classification{Kind: kind, MatchedBy: domain.MatchSourceExtension, Extension: ext}, nil
	}
	if !c.UseMagic {
		return Classification{Kind: "unknown", MatchedBy: domain.MatchSourceUnknown, Extension: ext}, nil
	}

	prefix, err := readPrefix(path)
	if err != nil {
		return Classification{}, err
	}
	if len(prefix) == 0 {
		return Classification{Kind: "unknown", MatchedBy: domain.MatchSourceUnknown, Extension: ext}, nil
	}

	if hasShebang(prefix) {
		return Classification{Kind: "code", MatchedBy: domain.MatchSourceMagic, Extension: ext, DetectedMIME: mimetypeDetect(prefix)}, nil
	}

	mime := mimetypeDetect(prefix)
	if kind, ok := c.Registry.KindForMIME(mime); ok {
		return Classification{Kind: kind, MatchedBy: domain.MatchSourceMagic, Extension: ext, DetectedMIME: mime}, nil
	}
	if strings.EqualFold(mime, "text/plain") || strings.HasPrefix(strings.ToLower(mime), "text/") {
		return Classification{Kind: "documents", MatchedBy: domain.MatchSourceMagic, Extension: ext, DetectedMIME: mime}, nil
	}

	// Fallback: check for formats whose signatures are beyond the 512-byte header.
	if detectISO9660(path) {
		return Classification{Kind: "diskimages", MatchedBy: domain.MatchSourceMagic, Extension: ext, DetectedMIME: "application/x-iso9660-image"}, nil
	}

	return Classification{Kind: "unknown", MatchedBy: domain.MatchSourceUnknown, Extension: ext, DetectedMIME: mime}, nil
}

func readPrefix(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(io.LimitReader(f, 512))
}

func hasShebang(data []byte) bool {
	data = trimSpace(data)
	return len(data) >= 2 && data[0] == '#' && data[1] == '!'
}

func trimSpace(data []byte) []byte {
	start := 0
	for start < len(data) {
		switch data[start] {
		case ' ', '\n', '\r', '\t', '\f', '\v':
			start++
		default:
			goto endStart
		}
	}
endStart:
	end := len(data)
	for end > start {
		switch data[end-1] {
		case ' ', '\n', '\r', '\t', '\f', '\v':
			end--
		default:
			return data[start:end]
		}
	}
	return data[start:end]
}
