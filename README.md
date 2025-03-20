# Extension Verifier

The idea of this tool is to provide a Tool which can find automated issues in your extensions, fix them automatically and also format them.

### Linting

```shell
docker run --rm -v $(pwd):/ext:ro ghcr.io/shopwarelabs/extension-verifier:latest check /ext
```

### Automatic Fixes

```shell
docker run --rm -v $(pwd):/ext ghcr.io/shopwarelabs/extension-verifier:latest fix /ext
```

For using experimental Twig AI diffing, you will need a Google Gemini API key, which needs to be set as environment variable `GEMINI_API_KEY`. The Key can be found in [Google AI Studio Dashboard](https://aistudio.google.com/).

### Formatting

```shell
docker run --rm -v $(pwd):/ext ghcr.io/shopwarelabs/extension-verifier:latest format /ext
```

## CI Usage

```yaml
jobs:
    check:
        runs-on: ubuntu-24.04
        strategy:
            fail-fast: false
            matrix:
                version-selection: [ 'lowest', 'highest']
        steps:
            - name: Checkout
              uses: actions/checkout@v4

            - name: Check extension
              uses: shopware/github-actions/extension-verifier@main
              with:
                   action: check
                   check-against: ${{ matrix.version-selection }}
```

## Reporter

The check command supports various reporters:

- summary
- markdown
- json
- github (GitHub Code annotations + Summary)
- junit

The JUnit output can be used by various CI/CD providers to provide better summary.

# Contribution

To run this tool locally you need PHP, Node, NPM and Go.

```shell
composer install -d tools/php
npm install --prefix tools/js

go run . <the command you wanna run>
```

# FAQ

## Missing classes in Storefront/Elasticsearch bundle

You're plugin typically require only `shopware/core`, but when you use classes from Storefront or Elasticsearch Bundle and they are required, you have to add `shopware/storefront` or `shopware/elasticsearch` also to the `require` in the composer.json. If those features are optional with `class_exists` checks, you want to add them into `require-dev`, so the dependencies are installed only for development and PHPStan can reconize the files.
