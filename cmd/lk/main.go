package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/amurru/lk/internal/domain"
	"github.com/amurru/lk/internal/kinds"
	"github.com/amurru/lk/internal/output"
	"github.com/amurru/lk/internal/walker"
)

const version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	registry := kinds.DefaultRegistry()
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "-h", "--help":
		printUsage(stdout)
		return 0
	case "-v", "--version":
		fmt.Fprintf(stdout, "lk version %s\n", version)
		return 0
	case "-k", "--kinds", "kinds":
		writeKinds(stdout, registry)
		return 0
	}

	kindName, ok := registry.CanonicalName(args[0])
	if !ok {
		fmt.Fprintf(stderr, "unknown kind %q\n", args[0])
		printUsage(stderr)
		return 2
	}

	fs := flag.NewFlagSet("lk", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { printUsage(stderr) }

	recursive := fs.Bool("recursive", false, "walk subdirectories recursively")
	fs.BoolVar(recursive, "r", false, "walk subdirectories recursively")
	hidden := fs.Bool("hidden", false, "include hidden files")
	fs.BoolVar(hidden, "a", false, "include hidden files")
	format := fs.String("format", "table", "output format")
	fs.StringVar(format, "f", "table", "output format")
	sortBy := fs.String("sort", "name", "sort by name, size, or modified")
	limit := fs.Int("limit", 0, "limit results")
	fs.IntVar(limit, "l", 0, "limit results")
	noMagic := fs.Bool("no-magic", false, "disable magic byte inspection")
	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
	versionFlag := fs.Bool("version", false, "show version")
	fs.BoolVar(versionFlag, "v", false, "show version")
	kindsFlag := fs.Bool("kinds", false, "list kinds")
	fs.BoolVar(kindsFlag, "k", false, "list kinds")

	if err := fs.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintln(stderr, err)
		return 2
	}

	if *help {
		printUsage(stdout)
		return 0
	}
	if *versionFlag {
		fmt.Fprintf(stdout, "lk version %s\n", version)
		return 0
	}
	if *kindsFlag {
		writeKinds(stdout, registry)
		return 0
	}

	outputFormat, err := parseFormat(*format)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	sortOrder, err := parseSort(*sortBy)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	path := "."
	if fs.NArg() > 1 {
		fmt.Fprintln(stderr, "expected at most one path argument")
		printUsage(stderr)
		return 2
	}
	if fs.NArg() == 1 {
		path = fs.Arg(0)
	}

	scanner := walker.New(registry, !*noMagic)
	entries, walkErrs, err := scanner.Walk(path, domain.ScanOptions{
		Recursive:     *recursive,
		IncludeHidden: *hidden,
		UseMagic:      !*noMagic,
		SortBy:        sortOrder,
		Limit:         *limit,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	filtered := filterByKind(entries, kindName)
	filtered = applyLimit(filtered, *limit)
	for _, walkErr := range walkErrs {
		fmt.Fprintln(stderr, walkErr)
	}

	writer := output.Writer{Stdout: stdout, Stderr: stderr}
	if err := writer.Write(filtered, kindName, outputFormat); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if len(walkErrs) > 0 {
		return 1
	}
	return 0
}

func parseFormat(raw string) (domain.OutputFormat, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "table":
		return domain.OutputFormatTable, nil
	case "simple":
		return domain.OutputFormatSimple, nil
	case "json":
		return domain.OutputFormatJSON, nil
	default:
		return "", fmt.Errorf("unknown format %q", raw)
	}
}

func parseSort(raw string) (domain.SortBy, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "name":
		return domain.SortByName, nil
	case "size":
		return domain.SortBySize, nil
	case "modified", "modtime":
		return domain.SortByModified, nil
	default:
		return "", fmt.Errorf("unknown sort order %q", raw)
	}
}

func filterByKind(entries []domain.FileEntry, kind string) []domain.FileEntry {
	filtered := make([]domain.FileEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.Kind == kind {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func applyLimit(entries []domain.FileEntry, limit int) []domain.FileEntry {
	if limit <= 0 || limit >= len(entries) {
		return entries
	}
	return entries[:limit]
}

func writeKinds(w io.Writer, registry kinds.Registry) {
	for _, kind := range registry.Kinds() {
		fmt.Fprintf(w, "%s\taliases=%s\textensions=%s\n", kind.Name, strings.Join(kind.Aliases, ","), strings.Join(kind.Extensions, ","))
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "lk - List Kind")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  lk <kind> [flags] [path]")
	fmt.Fprintln(w, "  lk kinds")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  -r, --recursive   walk subdirectories recursively")
	fmt.Fprintln(w, "  -a, --hidden      include hidden files")
	fmt.Fprintln(w, "  -f, --format      table, simple, or json")
	fmt.Fprintln(w, "  -s, --sort        name, size, or modified")
	fmt.Fprintln(w, "  -l, --limit       limit number of results")
	fmt.Fprintln(w, "      --no-magic    disable magic byte inspection")
	fmt.Fprintln(w, "  -k, --kinds       list kinds and exit")
	fmt.Fprintln(w, "  -h, --help       show help")
	fmt.Fprintln(w, "  -v, --version    show version")
}
