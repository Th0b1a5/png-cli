# png

`png` (short for *ping*) is a fast Linux CLI health-checker that pings grouped HTTP endpoints, measures latency (in µs/ms), and prints clean, colorized status lines in your terminal.

---

## Command Syntax & Usage

| Command | Action |
| :--- | :--- |
| `png` | Show help and usage |
| `png <group>` | Ping all URLs in `<group>` (fast shortcut) |
| `png ping <group>` | Ping all URLs in `<group>` |
| `png ls` | List all saved groups |
| `png ls <group>` | List URLs in `<group>` |
| `png new <group>` | Create a new `<group>` |
| `png add <group> <url>` | Add a URL to `<group>` |
| `png rm <group>` | Delete `<group>` (with confirmation) |
| `png rm <group> <url>` | Delete `<url>` from `<group>` |

> **Reserved Words**: `ls`, `new`, `add`, `rm`, `ping`, `version`, and `help` cannot be used as group names.

---

## Configuration

Config file is stored at `~/.config/png/config.json` (via `os.UserConfigDir()`):

```json
{
  "groups": [
    {
      "name": "local",
      "urls": [
        "http://localhost:8080/products",
        "http://localhost:8000/docs"
      ]
    }
  ]
}
```

---

## Roadmap & TODO

### 1. Config & Data Structures
- [x] Create config directory and default config file (`~/.config/png/config.json`)
- [x] Load and unmarshal config from disk into Go structs (`Config`, `Group` with `URLs []string`)
- [ ] Implement `saveConfig(path, cfg)` to write indented JSON back to disk

### 2. Command Dispatching (CLI Router)
- [x] Migrate from `flag` package to positional `os.Args` parsing
- [ ] Update router to Action-First syntax (`png <action> [group] [url]`)
- [ ] Implement `png ls` and `png ls <group>`
- [ ] Implement `png new <group>`
- [ ] Implement `png add <group> <url>`
- [ ] Implement `png rm <group>` and `png rm <group> <url>`
- [ ] Implement `png <group>` fallback for instant pinging

### 3. HTTP Pinging Engine & Terminal UI
- [ ] Ping URLs concurrently or sequentially using `net/http`
- [ ] Measure request latency using `time.Since`
- [ ] Fixed-width formatting (`%-30s`) so latency columns line up
- [ ] Colorized terminal output (Green for OK, Yellow for in-flight, Red for unreachable)

### 4. Distribution
- [ ] Build standalone static binary
- [ ] Create `PKGBUILD` for Arch Linux (AUR) release