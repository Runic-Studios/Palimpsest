# Palimpsest

Palimpsest is a simple CLI tool for merging config files across multiple overlay directories.

It supports merging the following file formats:
- YAML (`.yaml`, `.yml`)
- JSON (`.json`)
- TOML (`.toml`)
- Java Properties (`.properties`)

Overlays are applied in order: later overlays override earlier ones.

Palimpsest completely respects and preserves directory structure between overlays and for the final output.

---

## Usage

```bash
palimpsest --overlay overlay1/ --overlay overlay2/ --overlay overlay3/ --output output/
```

or shorthand:
```bash
palimpsest -o overlay1/ -o overlay2/ -o overlay3/ -t output/
```

### Arguments
- `--overlay (-o)` — Specify an overlay directory (can be repeated).
- `--output (-t)` — Specify the output directory (required).

## Example

Given the following overlays:
```
overlay1/config.yaml
overlay2/config.yaml
overlay3/config.yaml
```

If all contain the same key `foo`, the value from `overlay3` will win.

If only `overlay3` contains `foo`, that key will be copied into the output.

## Installation

### Build from source:
```bash
go build -o palimpsest ./cmd
```