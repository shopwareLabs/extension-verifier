/** @type {import('stylelint').Config} */
export default {
    extends: [
        "stylelint-config-recommended-scss"
    ],
    customSyntax: "postcss-scss",
    plugins: [
        "stylelint-scss"
    ],
    rules: {
        "selector-class-pattern": null,
        "import-notation": null,
        "declaration-property-value-no-unknown": null,
        "at-rule-no-unknown": null,
        "no-descending-specificity": null,
        "max-nesting-depth": [3, {
            "ignore": ["blockless-at-rules", "pseudo-classes"],
            "severity": "warning"
        }],
        "selector-max-type": [0, {
            "message": "Selectors containing elements like \"%s\" should be avoided because the element type might change. Prefer to use .classes, #ids and [data-attributes] instead."
        }],
    }
};
