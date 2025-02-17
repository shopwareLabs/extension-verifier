import stylelint from "stylelint";

export const ruleName = "shopware-administration/no-scss-extension-import";

export const messages = stylelint.utils.ruleMessages(ruleName, {
  rejected:
    'Avoid using the ".scss" extension in imports that start with "~scss". ' +
    'Use "@import \'~scss/variables\'" instead.'
});

export default stylelint.createPlugin(
  ruleName,
  (primaryOption, secondaryOptions, context) => {
    return (root, result) => {
      const validOptions = stylelint.utils.validateOptions(result, ruleName, {
        actual: primaryOption,
      });
      if (!validOptions) {
        return;
      }

      root.walkAtRules("import", (atRule) => {
        // Get import params; remove quotes and semicolon elements.
        let importPath = atRule.params.trim();

        // Remove wrapping quotes (either single or double)
        const matchQuotes = importPath.match(/^(['"])(.*)\1$/);
        if (matchQuotes) {
          importPath = matchQuotes[2];
        }

        // Check if it starts with "~scss" and ends with ".scss"
        if (importPath.startsWith("~scss") && importPath.endsWith(".scss")) {
          // Create the fixed import path by stripping the .scss extension.
          const fixedPath = importPath.replace(/\.scss$/, "");

          if (context.fix) {
            // Use the same quote character or default to double quotes.
            const quote = matchQuotes ? matchQuotes[1] : '"';
            atRule.params = `${quote}${fixedPath}${quote}`;
          } else {
            stylelint.utils.report({
              message: messages.rejected,
              node: atRule,
              result,
              ruleName,
            });
          }
        }
      });
    };
  }
);

