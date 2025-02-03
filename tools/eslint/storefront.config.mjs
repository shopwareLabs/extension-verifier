import globals from "globals";
import pluginJs from "@eslint/js";
import tseslint from "typescript-eslint";

import storefrontRules from './rules/storefront/index.js';


/** @type {import('eslint').Linter.Config[]} */
export default [
  {languageOptions: { globals: globals.browser }},
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  {
    plugins: {
      "shopware-storefront": storefrontRules,
    },
    rules: {
      "shopware-storefront/migrate-plugin-manager": 'error',
      'shopware-storefront/no-dom-access-helper': 'error',
      'shopware-storefront/no-http-client': 'error',
      'shopware-storefront/no-query-string': 'error',
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