import noSnippetImport from "./no-snippet-import.js";
import noSrcImport from "./no-src-import.js";

export default {
    rules: {
        "no-src-import": noSrcImport,
        "no-snippet-import": noSnippetImport
    },
    configs: {
        recommended: {
            plugins: ['shopware-admin'],
            rules: {
                'shopware-admin/no-src-import': 'error',
                'shopware-admin/no-snippet-import': 'error',
            }
        }
    }
}