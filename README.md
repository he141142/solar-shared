# solar-shared

The contract shared between the infra repo and the app repos of the solar trading system.

| Package | Contents | Used by |
|---|---|---|
| `domain` | Event names (source / detail-type), payloads (`Forecast`, `TradeIntent`, `OrderResult`), Parameter Store paths | solar-infra, every app |
| `bus` | Publish events to EventBridge | apps |
| `zonecfg` | Read and cache config from Parameter Store | apps |

## Versioning (semver)
Each consumer pins a version in its own `go.mod`, so upgrading is each repo's own decision:
- **patch** `v0.1.x`: bug fix, no contract change
- **minor** `v0.x.0`: add a field/event/param (backward compatible)
- **major**: rename or remove a field, change an event name. Rules and consumers must be updated, so **avoid this**; prefer adding a new event alongside the old one.

```bash
git tag v0.1.0 && git push origin v0.1.0
# in a consumer repo:
go get github.com/he141142/solar-shared@v0.1.0
```
