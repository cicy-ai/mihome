<h1 align="center">
  <img src="Meta.png" alt="Meta Kennel" width="200">
  <br>Meta Kernel<br>
</h1>

<h3 align="center">Another Mihomo Kernel.</h3>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/MetaCubeX/mihomo">
    <img src="https://goreportcard.com/badge/github.com/MetaCubeX/mihomo?style=flat-square">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/MetaCubeX/mihomo/Alpha?style=flat-square">
  <a href="https://github.com/MetaCubeX/mihomo/releases">
    <img src="https://img.shields.io/github/release/MetaCubeX/mihomo/all.svg?style=flat-square">
  </a>
  <a href="https://github.com/MetaCubeX/mihomo">
    <img src="https://img.shields.io/badge/release-Meta-00b4f0?style=flat-square">
  </a>
</p>

## cicy-ai fork additions

This fork ([cicy-ai/cicy-mihomo](https://github.com/cicy-ai/cicy-mihomo)) is a
multi-tenant proxy front-end for the cicy-ai worker fleet. The kernel is
unchanged; only the auth and rule-routing surface has new knobs so a single
mihomo instance can serve many ad-hoc workers without per-worker config churn.

### Global password — shared secret, free-form usernames (v1.10.1)

`globalPassword:` at the top of `mihomo.yaml` accepts **any non-empty user**
that presents this password. Callers don't need to be pre-declared under
`authentication:`. Per-user passwords (when declared) still work independently.

```yaml
globalPassword: "cicy_19879c4a4fd42ec545acce7bad39f584"

# Optional — only needed when you want a per-user secret instead of the global one.
authentication:
  - "service-a:per-user-pass"
```

With this, downstream rule matching can key off the request username
(`metadata.InUser`) without coordinating user registration up front.

### `IN-USER-PREFIX` rule — prefix-based user routing (v1.10.2)

Existing `IN-USER` matches by **exact** username equality, so routing N workers
required N lines. `IN-USER-PREFIX` matches `metadata.InUser` by **`HasPrefix`**,
collapsing a whole namespace of usernames into one rule:

```yaml
rules:
  # specific overrides come first (regular IN-USER still works)
  - IN-USER,w-10001,proxy_us
  # prefix catch-all: any user whose name starts with `w-` goes through default_proxy_group
  - IN-USER-PREFIX,w-,default_proxy_group
  - MATCH,REJECT
```

Multiple prefixes are supported in a single rule using `/` as the separator
(same convention as `IN-USER`): `IN-USER-PREFIX,w-/svc-,default_proxy_group`.

Behavior:

| Username       | Matches               | Routed to            |
|----------------|-----------------------|----------------------|
| `w-10001`      | `IN-USER` (explicit)  | `proxy_us`           |
| `w-99999`      | `IN-USER-PREFIX,w-`   | `default_proxy_group`|
| `w-anything`   | `IN-USER-PREFIX,w-`   | `default_proxy_group`|
| `foo`          | `MATCH`               | `REJECT`             |

Combined with `globalPassword`, new workers come online without editing
`mihomo.yaml` — they just authenticate with their username + the shared
password, and the prefix rule decides where their traffic goes.

### Tag-driven multi-platform release workflow

Pushing a tag matching `v*`, `alpha*`, `Alpha*`, `beta*`, or `Beta*` triggers
`.github/workflows/release.yml`, which builds linux/darwin/windows ×
amd64/arm64 via `make all VERSION=<tag>` and publishes a GitHub Release with
SHA256SUMS attached. Tags containing `alpha` / `beta` / `rc` are marked as
prerelease automatically.

---

## Features

- Local HTTP/HTTPS/SOCKS server with authentication support
- VMess, VLESS, Shadowsocks, Trojan, Snell, TUIC, Hysteria protocol support
- Built-in DNS server that aims to minimize DNS pollution attack impact, supports DoH/DoT upstream and fake IP.
- Rules based off domains, GEOIP, IPCIDR or Process to forward packets to different nodes
- Remote groups allow users to implement powerful rules. Supports automatic fallback, load balancing or auto select node
  based off latency
- Remote providers, allowing users to get node lists remotely instead of hard-coding in config
- Netfilter TCP redirecting. Deploy Mihomo on your Internet gateway with `iptables`.
- Comprehensive HTTP RESTful API controller

## Dashboard

A web dashboard with first-class support for this project has been created; it can be checked out at [metacubexd](https://github.com/MetaCubeX/metacubexd).

## Configration example

Configuration example is located at [/docs/config.yaml](https://github.com/MetaCubeX/mihomo/blob/Alpha/docs/config.yaml).

## Docs

Documentation can be found in [mihomo Docs](https://wiki.metacubex.one/).

## For development

Requirements:
[Go 1.20 or newer](https://go.dev/dl/)

Build mihomo:

```shell
git clone https://github.com/MetaCubeX/mihomo.git
cd mihomo && go mod download
go build
```

Set go proxy if a connection to GitHub is not possible:

```shell
go env -w GOPROXY=https://goproxy.io,direct
```

Build with gvisor tun stack:

```shell
go build -tags with_gvisor
```

### IPTABLES configuration

Work on Linux OS which supported `iptables`

```yaml
# Enable the TPROXY listener
tproxy-port: 9898

iptables:
  enable: true # default is false
  inbound-interface: eth0 # detect the inbound interface, default is 'lo'
```

## Debugging

Check [wiki](https://wiki.metacubex.one/api/#debug) to get an instruction on using debug
API.

## Credits

- [Dreamacro/clash](https://github.com/Dreamacro/clash)
- [SagerNet/sing-box](https://github.com/SagerNet/sing-box)
- [riobard/go-shadowsocks2](https://github.com/riobard/go-shadowsocks2)
- [v2ray/v2ray-core](https://github.com/v2ray/v2ray-core)
- [WireGuard/wireguard-go](https://github.com/WireGuard/wireguard-go)
- [yaling888/clash-plus-pro](https://github.com/yaling888/clash)

## License

This software is released under the GPL-3.0 license.

**In addition, any downstream projects not affiliated with `MetaCubeX` shall not contain the word `mihomo` in their names.**