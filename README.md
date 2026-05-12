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

| Key | Action |
|-----|--------|
| `↑↓` | Navigate requests |
| `←→` | Switch detail tabs |
| `Tab` | Toggle list / detail focus |
| `f` | Filter (URL · `m:POST` · `s:4xx`) |
| `y` | Copy curl command |
| `/` | Search response body |
| `c` | Clear request list |
| `q` | Quit |
