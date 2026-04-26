# Kugou Music Go SDK

[中文](README.md) | English

> **Note**: This project is a Go version migrated and updated using Codex based on the source project [MakcRe/KuGouMusicApi](https://github.com/MakcRe/KuGouMusicApi) (JavaScript version). Due to potential discrepancies between the migration implementation and API behavior, the current version may contain issues. Please use with caution and evaluate risks independently.

## API Documentation

- Technical catalog (routes/methods/models): `sdk/API_CATALOG.md`
- Chinese descriptions (auto-extracted from comments): `sdk/API_CATALOG_ZH.md`
- Compatibility audit (per-interface risk assessment): `sdk/API_COMPAT_AUDIT.md`
- Fix status (per-interface verification progress): `sdk/API_FIX_STATUS.md`
- Signature/encryption checklist: `sdk/API_SIGN_CHECK.md`

## Package

```go
import "github.com/lfhy/kugou-music-api/sdk"
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/lfhy/kugou-music-api/sdk"
)

func main() {
    client := sdk.NewClient()
    
    // Example: Get daily recommendations
    resp, err := client.DailyRecommend(context.Background(), &sdk.DailyRecommendRequest{})
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Recommendations: %+v\n", resp)
}
```

## License

MIT License - see [LICENSE](LICENSE) file for details.
