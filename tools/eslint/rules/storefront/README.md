# @shopware-ag/storefront-eslint-rules

This package provides custom ESLint rules for Shopware Storefront projects. These rules help to ensure code quality, promote best practices, and assist in migrating from older code patterns to more modern approaches.

## Installation

```
npm install @shopware-ag/storefront-eslint-rules --save-dev
```

## Usage

Add the following to your ESLint configuration file (e.g., .eslintrc.js):

```diff
+import storefrontRules from '@shopware-ag/storefront-eslint-rules';

/** @type {import('eslint').Linter.Config[]} */
export default [
  {languageOptions: { globals: globals.browser }},
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
+  storefrontRules,
  {
    rules: {
      '@typescript-eslint/no-unused-vars': 'warn',
      '@typescript-eslint/no-unused-expressions': 'warn',
      '@typescript-eslint/no-this-alias': 'warn',
      '@typescript-eslint/no-require-imports': 'off',
      'no-undef': 'off',
      'no-alert': 'error',
      'no-console': ['error', { allow: ['warn', 'error'] }],
    }
  }
];
```

## Rules

### migrate-plugin-manager

This rule flags imports from src/plugin-system/plugin.manager and suggests using window.PluginManager instead. It also flags imports from src/plugin-system/plugin.class and suggests using window.PluginBaseClass.

### no-dom-access-helper

This rule identifies usages of the DomAccessHelper and suggests using native DOM methods instead.

### no-http-client

This rule identifies usages of the HttpClient and suggests using the fetch API instead.

### no-query-string

This rule identifies usages of the query-string library and suggests using URLSearchParams instead.

## Contributing

Contributions are welcome! Please read the contributing guidelines before submitting a pull request.

## License

MIT
