import MigratePluginManager from './plugin-manager.js';
import DomAccessHelper from "./dom-access-helper.js";
import HttpClient from "./http-client.js";
import QueryString from "./query-string.js";

export default {
    rules: {
        "migrate-plugin-manager": MigratePluginManager,
        "no-dom-access-helper": DomAccessHelper,
        "no-http-client": HttpClient,
        'no-query-string': QueryString,
    },
    configs: {
        recommended: {
            plugins: ['shopware-storefront'],
            rules: {
                'shopware-storefront/migrate-plugin-manager': 'error',
                'shopware-storefront/no-dom-access-helper': 'error',
                'shopware-storefront/no-http-client': 'error',
                'shopware-storefront/no-query-string': 'error',
            }
        }
    }
}