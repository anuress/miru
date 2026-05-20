# miru 見る

Live OkHttp network traffic in your terminal. An alternative to the Android Studio OkHttp Profiler plugin — useful when you want to monitor network calls without opening the IDE.

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
| `r` | Copy response body (pretty-printed) |
| `R` | Copy raw request (headers + body) |
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

## Credits

miru is built on top of the [OkHttp Profiler](https://github.com/itkacher/OkHttpProfiler) library by [nerdythings](https://github.com/itkacher). The interceptor handles data capture on the Android side — miru provides the terminal UI to view that data without requiring Android Studio.
