# png

`png` (short for *ping*) is a fast Linux CLI health-checker that pings grouped HTTP endpoints, measures latency (in µs/ms), and prints clean, colorized status lines in your terminal.

---

## Command Syntax & Usage

| Command | Action |
| :--- | :--- |
| `png` | List all saved groups |
| `png <group>` | Ping all endpoints in `<group>` |
| `png <group> ls` | List endpoints (name, URL, index) in `<group>` |
| `png <group> new` | Create a new `<group>` |
| `png <group> add <name> <url>` | Add an endpoint to `<group>` |
| `png <group> rm` | Delete `<group>` (with confirmation) |
| `png <group> rm <name>` | Delete endpoint `<name>` from `<group>` |

---

## Configuration

Config file is stored at `~/.config/png/config.json` (via `os.UserConfigDir()`):

```json
{
  "groups": [
    {
      "name": "local",
      "endpoints": [
        {
          "name": "products",
          "url": "http://localhost:8080/products"
        },
        {
          "name": "docs",
          "url": "http://localhost:8000/docs"
        }
      ]
    }
  ]
}
```

---

## Roadmap & TODO

### 1. Config & Data Structures
- [x] Create config directory and default config file (`~/.config/png/config.json`)
- [x] Load and unmarshal config from disk into Go structs
- [ ] Add `Endpoint` struct (`name`, `url`) and update `Group` to use `[]Endpoint`
- [ ] Implement `saveConfig(path, cfg)` to write indented JSON back to disk

### 2. Command Dispatching (CLI Router)
- [ ] Migrate from `flag` package to positional `os.Args` parsing
- [ ] Implement `png` (list all groups)
- [ ] Implement `png <group> new`
- [ ] Implement `png <group> ls`
- [ ] Implement `png <group> add <name> <url>`
- [ ] Implement `png <group> rm` (group deletion with confirmation)
- [ ] Implement `png <group> rm <name>` (endpoint deletion)

### 3. HTTP Pinging Engine
- [ ] Ping URLs concurrently or sequentially using `net/http`
- [ ] Measure request latency using `time.Since`
- [ ] Formatted, colorized terminal output (Status Code, Latency, Endpoint Name)

### 4. Distribution
- [ ] Build standalone static binary
- [ ] Create `PKGBUILD` for Arch Linux (AUR) release