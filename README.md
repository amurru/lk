# LK - List Kind

List files by their kind (e.g. documents, images, videos, ...etc).

## Usage

```text
lk <kind> [flags] [path]
lk kinds
```

### Examples

```text
lk documents
lk images /path/to/files
lk videos --recursive --hidden /path/to/files
lk documents --format json /path/to/files
lk documents --no-magic /path/to/files
lk kinds
```

### Flags

- `-r`, `--recursive`: walk subdirectories recursively
- `-a`, `--hidden`: include hidden files
- `-f`, `--format`: `table`, `simple`, or `json`
- `-s`, `--sort`: `name`, `size`, or `modified`
- `-n`, `--limit`: limit the number of results
- `--no-magic`: disable magic-byte inspection
