/** @type {import('stylelint').Config} */
export default {
    extends: [
        "stylelint-config-standard"
    ],
    plugins: [
        "stylelint-scss",
        "./rules/administration/wrong-scss-import.js"
    ],
    rules: {
        "selector-class-pattern": null,
        "import-notation": null,
        "declaration-property-value-no-unknown": null,
        "at-rule-no-unknown": null,
        "shopware-administration/no-scss-extension-import": true
    }
};