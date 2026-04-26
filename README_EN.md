# Kugou Music Go SDK

[中文](README.md) | English

> **Note**: This project is a Go version migrated and updated using Codex based on the source project [MakcRe/KuGouMusicApi](https://github.com/MakcRe/KuGouMusicApi) (JavaScript version). Due to potential discrepancies between the migration implementation and API behavior, the current version may contain issues. Please use with caution and evaluate risks independently.

## Quick Start

Install:

```bash
go get github.com/lfhy/kugou-music-api
```

Recommended root import:

```go
import kg "github.com/lfhy/kugou-music-api"
```

Minimal example:

```go
package main

import (
	"context"
	"fmt"
	"log"

	kg "github.com/lfhy/kugou-music-api"
)

func main() {
	client, err := kg.New()
	if err != nil {
		log.Fatalf("init client failed: %v", err)
	}

	resp, err := client.Search(context.Background(), kg.SearchRequest{
		Keywords: "Jay Chou",
		Page:     1,
		Pagesize: 10,
	})
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	fmt.Printf("status: %d\n", resp.Status)
	fmt.Printf("body: %s\n", string(resp.RawBody))
}
```

Compatibility:

- New projects should import the root package `github.com/lfhy/kugou-music-api`
- Existing code can continue using `github.com/lfhy/kugou-music-api/sdk`

Session and auto refresh:

- The SDK creates clients with the `lite` platform behavior by default
- When an API call returns an expired-login response, currently matched mainly by `error_code=20018`, the SDK will try `LoginByToken(...)` with the current `token + userid` and replay the original request once
- If `token/userid` are missing, or the refresh attempt fails, the SDK returns the error to the caller
- If you want to manage session refresh yourself, disable it explicitly:

```go
client, err := kg.New(
    kg.WithCookie(cookie),
    kg.WithAutoRefresh(false),
)
```

## API Documentation

- Technical catalog (routes/methods/models): `sdk/API_CATALOG.md`
- Chinese descriptions (auto-extracted from comments): `sdk/API_CATALOG_ZH.md`
- Compatibility audit (per-interface risk assessment): `sdk/API_COMPAT_AUDIT.md`
- Fix status (per-interface verification progress): `sdk/API_FIX_STATUS.md`
- Signature/encryption checklist: `sdk/API_SIGN_CHECK.md`

## License

MIT License - see [LICENSE](LICENSE) file for details.
