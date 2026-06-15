# Bundled libmoss Artifacts

This directory contains prebuilt `libmoss` C SDK artifacts used by the Go
bindings when built with `-tags libmoss`.

Source release:

- <https://github.com/usemoss/moss/releases/tag/c-sdk-v0.9.0>

Bundled targets:

- `darwin_arm64/libmoss.dylib`
- `linux_amd64/libmoss.so`
- `linux_arm64/libmoss.so`
- `windows_amd64/moss.dll`
- `windows_amd64/moss.dll.lib`
- `include/libmoss.h`

Upstream archive checksums:

```text
fd403a9bdce3644ed5eb892cb84f57e58f2d2a6aba8d1a45c74a296e2754a9c8  libmoss-v0.9.0-aarch64-apple-darwin.tar.gz
ad7b860235a313578bdec354591715ce81589533287646c4a9f5034dacd1790a  libmoss-v0.9.0-aarch64-unknown-linux-gnu.tar.gz
d0392d896ec157bc88332ee5e7c2bedd7bb52ea6fabc34919e6b1b40e5e20f1c  libmoss-v0.9.0-x86_64-pc-windows-msvc.tar.gz
d0d35c717dc7567c68831e8866fb690d9d7f44ee7746aa628b9c47fd55c32b4d  libmoss-v0.9.0-x86_64-unknown-linux-gnu.tar.gz
```
