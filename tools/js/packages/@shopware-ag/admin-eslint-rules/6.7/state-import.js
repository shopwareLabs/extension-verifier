export default {
  meta: {
    type: "suggestion",
    docs: {
      description:
        "Replace Shopware.State with Shopware.Store (and destructured State accordingly).",
      category: "Best Practices",
      recommended: false
    },
    fixable: "code",
    schema: [],
    minShopwareVersion: '6.7.0.0'
  },

  create(context) {
    const sourceCode = context.getSourceCode();
    const stateVariableNames = new Set();

    function isShopwareState(node) {
      return (
        node &&
        node.type === "MemberExpression" &&
        !node.computed &&
        node.object &&
        node.object.type === "Identifier" &&
        node.object.name === "Shopware" &&
        node.property &&
        node.property.type === "Identifier" &&
        node.property.name === "State"
      );
    }

    return {
      // Check for direct usage like: Shopware.State.get('context')...
      MemberExpression(node) {
        if (isShopwareState(node)) {
          context.report({
            node,
            message:
              "Do not use 'Shopware.State', use 'Shopware.Store' instead.",
            fix(fixer) {
              // Replace "...State" to "...Store"
              const stateText = sourceCode.getText(node.property);
              if (stateText === "State") {
                return fixer.replaceText(node.property, "Store");
              }
              return null;
            }
          });
        }
      },

      // Check for destructuring: const { State } = Shopware;
      VariableDeclarator(node) {
        if (
          node.init &&
          node.init.type === "Identifier" &&
          node.init.name === "Shopware" &&
          node.id.type === "ObjectPattern"
        ) {
          node.id.properties.forEach((prop) => {
            if (
              prop.type === "Property" &&
              prop.key &&
              prop.key.type === "Identifier" &&
              prop.key.name === "State"
            ) {
              // Mark this local variable name (it could be renamed)
              const localName = prop.value.name;
              stateVariableNames.add(localName);
              context.report({
                node: prop,
                message:
                  "Do not use destructured 'State', use destructured 'Store' instead.",
                fix(fixer) {
                  // Fix the property key to change "State" to "Store"
                  // Preserve possible aliasing e.g., { State: MyState }
                  const fixedKey = fixer.replaceText(
                    prop.key,
                    "Store"
                  );
                  return fixedKey;
                }
              });
            }
          });
        }
      },

      // Replace anywhere in the code where the destructured State variable is used.
      Identifier(node) {
        if (
          stateVariableNames.has(node.name) &&
          // make sure it's not a key in an object property, etc.
          node.parent &&
          // avoid double fixing if it is in a declaration already (we already fixed it)
          node.parent.type !== "Property" &&
          node.parent.type !== "VariableDeclarator"
        ) {
          context.report({
            node,
            message:
              "Do not use destructured 'State', use 'Store' instead.",
            fix(fixer) {
              return fixer.replaceText(node, "Store");
            }
          });
        }
      }
    };
  }
};
