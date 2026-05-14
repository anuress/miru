# miru 見る

Live OkHttp network traffic in your terminal. Replaces the Android Studio OkHttp Profiler plugin for IDE-free debugging.

## Install

```bash
go install github.com/anuress/miru@latest
```

## Android setup

Add the interceptor to your app's `build.gradle`:

```groovy
debugImplementation 'io.nerdythings:okhttpprofiler:1.0.8'
```

Register it on your `OkHttpClient`:

```kotlin
OkHttpClient.Builder()
    .addInterceptor(OkHttpProfilerInterceptor())
    .build()
```

## Usage

```bash
miru                              # pick device + process interactively
miru --device emulator-5554       # skip device picker
miru --process com.myapp.debug    # skip process picker
miru --theme github-dark          # use GitHub Dark theme (default: catppuccin-mocha)
```

## Configuration

Create `~/.config/miru/config.json` to persist preferences:

```json
{
  "theme": "catppuccin-mocha"
}
```

Available themes: `catppuccin-mocha` (default), `github-dark`

## Keys

### Request list (left pane)

| Key | Action |
|-----|--------|
| `↑↓` | Navigate requests (newest at top) |
| `f` | Filter by URL · `m:POST` · `s:4xx` |
| `y` | Copy curl command for selected request |
| `c` | Clear request list |
| `Tab` | Switch focus to detail pane |
| `q` | Quit |

### Detail pane (right pane)

| Key | Action |
|-----|--------|
| `↑↓` / `j` `k` | Move cursor line by line |
| `Ctrl+U` / `Ctrl+D` | Jump half page up/down |
| `gg` / `G` | Jump to first / last line |
| `←→` | Switch tabs |
| `y` | Copy value at cursor (header value · JSON value · block) |
| `Y` | Copy full raw line at cursor |
| `/` | Search — type query, `↑↓` navigate matches, `esc` clear |
| `Tab` | Switch focus back to list |

### Tabs

| Tab | Content |
|-----|---------|
| RAW REQUEST | Request headers + body |
| REQ HEADERS | Request headers only |
| RESP HEADERS | Response headers only |
| RESP BODY | Response body (pretty-printed JSON) |

## How it works

miru reads from `adb logcat`, filtering for `OKPRFL_*` log lines emitted by the OkHttpProfilerInterceptor. It assembles request/response pairs by their unique ID and streams them into the TUI in real time — no TCP sockets, no port forwarding.
