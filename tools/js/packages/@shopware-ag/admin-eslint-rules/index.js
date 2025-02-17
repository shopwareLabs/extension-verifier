import noSnippetImport from "./no-snippet-import.js";
import noSrcImport from "./no-src-import.js";

export default {
    plugins: {
        "shopware-admin": {
            rules: {
                "no-src-import": noSrcImport,
                "no-snippet-import": noSnippetImport
            },
        }
    },
    rules: {
        'shopware-admin/no-src-import': 'error',
        'shopware-admin/no-snippet-import': 'error',
    }
}