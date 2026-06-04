package walker

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/amurru/lk/internal/classifier"
	"github.com/amurru/lk/internal/domain"
	"github.com/amurru/lk/internal/kinds"
)

type Walker struct {
	Classifier classifier.Classifier
}

func New(registry kinds.Registry, useMagic bool) Walker {
	return Walker{Classifier: classifier.Classifier{Registry: registry, UseMagic: useMagic}}
}

func (w Walker) Walk(root string, opts domain.ScanOptions) ([]domain.FileEntry, []error, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, nil, err
	}

	var entries []domain.FileEntry
	var errs []error

	addFile := func(path string, name string, info os.FileInfo) {
		if !opts.IncludeHidden && isHidden(name) {
			return
		}
		classification, err := w.Classifier.Classify(path)
		if err != nil {
			errs = append(errs, err)
			return
		}
		entries = append(entries, domain.FileEntry{
			Path:       path,
			Name:       name,
			Extension:  kinds.NormalizeExtension(filepath.Ext(path)),
			SizeBytes:  info.Size(),
			ModifiedAt: info.ModTime().UTC(),
			Kind:       classification.Kind,
			MatchedBy:  classification.MatchedBy,
		})
	}

	if !info.IsDir() {
		addFile(root, info.Name(), info)
		sortEntries(entries, opts.SortBy)
		return entries, errs, nil
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			errs = append(errs, walkErr)
			return nil
		}
		if d.IsDir() {
			if path == root {
				return nil
			}
			if !opts.Recursive {
				return filepath.SkipDir
			}
			if !opts.IncludeHidden && isHidden(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !opts.IncludeHidden && isHidden(d.Name()) {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			errs = append(errs, err)
			return nil
		}
		addFile(path, d.Name(), info)
		return nil
	})
	if err != nil {
		return nil, errs, err
	}

	sortEntries(entries, opts.SortBy)
	return entries, errs, nil
}

func sortEntries(entries []domain.FileEntry, sortBy domain.SortBy) {
	switch sortBy {
	case "", domain.SortByName:
		sort.SliceStable(entries, func(i, j int) bool {
			if entries[i].Name == entries[j].Name {
				return entries[i].Path < entries[j].Path
			}
			return entries[i].Name < entries[j].Name
		})
	case domain.SortBySize:
		sort.SliceStable(entries, func(i, j int) bool {
			if entries[i].SizeBytes == entries[j].SizeBytes {
				return entries[i].Path < entries[j].Path
			}
			return entries[i].SizeBytes < entries[j].SizeBytes
		})
	case domain.SortByModified:
		sort.SliceStable(entries, func(i, j int) bool {
			if entries[i].ModifiedAt.Equal(entries[j].ModifiedAt) {
				return entries[i].Path < entries[j].Path
			}
			return entries[i].ModifiedAt.Before(entries[j].ModifiedAt)
		})
	}
}

func applyLimit(entries []domain.FileEntry, limit int) []domain.FileEntry {
	if limit <= 0 || limit >= len(entries) {
		return entries
	}
	return entries[:limit]
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}
