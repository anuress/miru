# miru 見る

Live OkHttp network traffic in your terminal.

## Install

```bash
go install github.com/anuress/miru@latest
```

## Android setup

Add to your app's `build.gradle`:

```groovy
debugImplementation 'com.itkacher.okhttpprofiler:okhttpprofiler:1.0.8'
```

Add the interceptor:

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
miru --port 6360                  # custom port (default: 6360)
```

## Keys

### Request list (left pane)

| Key | Action |
|-----|--------|
| `↑↓` | Navigate requests |
| `Tab` | Switch focus to detail pane |
| `f` | Filter (URL · `m:POST` · `s:4xx`) |
| `y` | Copy curl command for selected request |
| `c` | Clear request list |
| `q` | Quit |

### Detail pane (right pane)

| Key | Action |
|-----|--------|
| `↑↓` / `j` `k` | Move cursor line by line |
| `Ctrl+U` / `Ctrl+D` | Move cursor half page up/down |
| `gg` | Jump cursor to first line |
| `G` | Jump cursor to last line |
| `←→` | Switch tabs (resets cursor) |
| `y` | Copy value under cursor (header value / JSON value / block) |
| `Y` | Copy full raw line under cursor |
| `/` | Search in response body |
| `n` / `N` | Next / previous search match (also moves cursor) |
| `Tab` | Switch focus back to list |
