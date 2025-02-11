import globals from "globals";
import pluginJs from "@eslint/js";
import tseslint from "typescript-eslint";

import adminRules from './rules/admin/index.js';

/** @type {import('eslint').Linter.Config[]} */
export default [
  { languageOptions: { globals: globals.browser } },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  {
    plugins: {
      "shopware-admin": adminRules,
    },
    rules: {
      ...adminRules.configs.recommended.rules,
      'no-undef': 'off',
      'no-alert': 'error',
      'no-console': ['error', { allow: ['warn', 'error'] }],
    }
  }
];