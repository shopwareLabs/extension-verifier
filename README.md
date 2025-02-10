# Extension Verifier


**This tool is still in development and experimental, maybe not finished**

The idea of this tool is to provide a Tool which can find automated issues in your extensions, fix them automatically and also format them.

### Linting

```shell
docker run --rm -v $(pwd):/ext:ro ghcr.io/shopwarelabs/extension-verifier:latest check /ext
```

### Automatic Fixes

```shell
docker run --rm -v $(pwd):/ext:ro ghcr.io/shopwarelabs/extension-verifier:latest fix /ext
```

### Formatting

```shell
docker run --rm -v $(pwd):/ext:ro ghcr.io/shopwarelabs/extension-verifier:latest format /ext
```
