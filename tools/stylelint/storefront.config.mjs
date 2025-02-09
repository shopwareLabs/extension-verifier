/** @type {import('stylelint').Config} */
export default {
    extends: [
        "stylelint-config-standard"
    ],
    plugins: [
        "stylelint-scss"
    ],
    rules: {
        "selector-class-pattern": null,
        "import-notation": null,
        "declaration-property-value-no-unknown": null,
        "at-rule-no-unknown": null,
        "no-descending-specificity": null,
    }
};