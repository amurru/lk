# lk — List Kind

A fast, zero-config file lister that groups files by kind instead of making you remember extensions. Find all images, documents, videos, audio, code, archives, and more in a directory with a single command. Supports magic-byte detection for more accurate classification, multiple output formats, and configurable sorting.

## Usage

```
lk <kind> [flags] [path]
lk kinds
```

If no path is given, the current directory is used.

### Kinds

| Kind        | Aliases              | Extensions                                     |
|-------------|----------------------|------------------------------------------------|
| archives    | archive, zips        | 7z, bz2, gz, rar, tar, tgz, zip               |
| audio       | sound, music         | aac, ac3, aif, aiff, amr, au, flac, m4a, mid, midi, mka, mp3, mpc, ogg, opus, wav, wma, wv |
| code        | script, source, src  | c, cpp, css, go, h, hpp, html, java, js, json, py, rb, sh, ts, tsx, yaml, yml |
| databases   | data, db             | accdb, db, dbf, mdb, sqlite, sqlite3           |
| diskimages  | disk, img, iso       | dmg, img, iso, vhd, vmdk, wim                  |
| documents   | doc, docs            | csv, doc, docm, docx, dot, dotm, dotx, epub, log, md, odp, ods, odt, pdf, pot, potm, potx, pps, ppsm, ppsx, ppt, pptm, pptx, rtf, txt, xls, xlsb, xlsm, xlsx, xlt, xltm, xltx, xps |
| executables | bin, binary, exe     | bat, cmd, com, dll, dylib, exe, ps1, so        |
| fonts       | font, typefaces      | eot, otf, ttf, woff, woff2                     |
| images      | image, photo, pics   | bmp, gif, jpeg, jpg, png, svg, tif, tiff, webp |
| models      | 3d, cad              | 3ds, blend, dae, fbx, glb, gltf, obj, ply, stl, u3d |
| videos      | movie, vid, vids     | avi, m4v, mkv, mov, mp4, mpeg, mpg, webm       |

### Flags

| Flag          | Short | Description                       |
|---------------|-------|-----------------------------------|
| `--recursive` | `-r`  | walk subdirectories recursively   |
| `--hidden`    | `-a`  | include hidden files              |
| `--format`    | `-f`  | output format: `table` (default), `simple`, `json`, `xml` |
| `--sort`      | `-s`  | sort by: `name` (default), `size`, `modified` |
| `--limit`     | `-l`  | limit number of results           |
| `--null`      | `-0`  | null-terminate output (for xargs) |
| `--exec`      |       | run command per file (`{}` = path)|
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
lk fonts /path/to/project
lk executables -0 | xargs -0 file
lk images --exec md5sum {}
lk kinds
```

## Install

```sh
go build -o lk ./cmd/lk/
```
