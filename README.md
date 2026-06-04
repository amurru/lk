# lk — List Kind

A fast, zero-config file lister that groups files by kind instead of making you remember extensions. Find all images, documents, videos, code, or archives in a directory with a single command. Supports magic-byte detection for more accurate classification, multiple output formats, and configurable sorting.

## Usage

```
lk <kind> [flags] [path]
lk kinds
```

If no path is given, the current directory is used.

### Kinds

| Kind      | Aliases              | Extensions                                     |
|-----------|----------------------|------------------------------------------------|
| archives  | archive, zips        | 7z, bz2, gz, rar, tar, tgz, zip               |
| code      | script, source, src  | c, cpp, css, go, h, hpp, html, java, js, json, py, rb, sh, ts, tsx, yaml, yml |
| documents | doc, docs            | csv, doc, docm, docx, dot, dotm, dotx, epub, log, md, odp, ods, odt, pdf, pot, potm, potx, pps, ppsm, ppsx, ppt, pptm, pptx, rtf, txt, xls, xlsb, xlsm, xlsx, xlt, xltm, xltx, xps |
| images    | image, photo, pics   | bmp, gif, jpeg, jpg, png, svg, tif, tiff, webp |
| videos    | movie, vid, vids     | avi, m4v, mkv, mov, mp4, mpeg, mpg, webm       |

### Flags

| Flag          | Short | Description                       |
|---------------|-------|-----------------------------------|
| `--recursive` | `-r`  | walk subdirectories recursively   |
| `--hidden`    | `-a`  | include hidden files              |
| `--format`    | `-f`  | output format: `table` (default), `simple`, `json`, `xml` |
| `--sort`      | `-s`  | sort by: `name` (default), `size`, `modified` |
| `--limit`     | `-l`  | limit number of results           |
| `--no-magic`  |       | disable magic-byte inspection     |
| `--kinds`     | `-k`  | list kinds and exit               |

### Examples

```sh
lk documents
lk images /path/to/files
lk videos -ra /path/to/files
lk documents -f json /path/to/files
lk code -f xml -r /path/to/files
lk documents -s size -l 10
lk documents --no-magic /path/to/files
lk kinds
```

## Install

```sh
go build -o lk ./cmd/lk/
```
