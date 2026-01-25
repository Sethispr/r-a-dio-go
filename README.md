# r-a-dio-go (in development, new qol soon)

[![Go Report Card](https://goreportcard.com/badge/github.com/sethispr/r-a-dio-go)](https://goreportcard.com/report/github.com/sethispr/r-a-dio-go) [![Go Version](https://img.shields.io/github/go-mod/go-version/sethispr/r-a-dio-go)](https://golang.org/doc/devel/release.html)

Barebones implementation for bypassing [r-a-d.io](https://r-a-d.io/search)'s 30 minute request limit. ([see note](https://github.com/Sethispr/r-a-dio-go/blob/main/README.md#note))

This cli [gets proxies](https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt) then reads the site source code and scrape the one time [gorilla/csrf](https://github.com/gorilla/csrf) token and tricks the server to accept your song request.

- go's concurrency can verify proxies and go through all song reqs in less than a second
- stateless networking net/http which disposes each request making it fresh
- regex to get csrf token and stores it in a cookiejar
- inject POST request with stolen token, song id and spoofed Referer headers

---

## Installation

```bash
go install github.com/sethispr/r-a-dio-go@latest
```

## Build from source

1. **Clone repo**

```bash
git clone https://github.com/sethispr/r-a-dio-go.git
cd r-a-dio-go
```

2. **Build executable**

```bash
go build .
```

### Usage

Run the compiled binary from your terminal:

```bash
r-a-dio-go
```

> [!NOTE]
> For precompiled binaries for different OS architectures, check the [Releases](https://www.google.com/search?q=https://github.com/sethispr/r-a-dio-go/releases) repo page.

---

### Disclaimer
This tool is for educational and research purposes only. By using this software, you acknowledge that:

The developer does not condone or support the malicious spamming of community run services.

You assume all risks associated with bypassing rate limits or automated scraping, the developer is not responsible for any IP bans, blacklisting, or legal consequences.

### Note

> [!WARNING]
> The staff or anyone with admin can just remove all your spammed songs from queue also wasting the song's cd.
> Use this cli in a resepectable or non suspicious way so your songs don't just get removed.
