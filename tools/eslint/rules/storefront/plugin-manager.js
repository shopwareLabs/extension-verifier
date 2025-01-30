export default {
  meta: {
    type: "suggestion",
    docs: {
      description: "Migrate PluginManager import to window.PluginManager assignment",
      category: "Migration",
      recommended: false
    },
    fixable: "code"
  },

  create(context) {
    return {
      ImportDeclaration(node) {
        // Check if the import is from 'src/plugin-system/plugin.manager'
        if (node.source.value === "src/plugin-system/plugin.manager") {
          // Get the imported variable name
          const importedName = node.specifiers[0]?.local?.name;

          if (importedName) {
            context.report({
              node,
              message: `Import from plugin.manager should use window.PluginManager`,
              fix(fixer) {
                return fixer.replaceText(
                  node,
                  `const ${importedName} = window.PluginManager;`
                );
              },
            });
          }
        }

        if (node.source.value === 'src/plugin-system/plugin.class') {
          // Get the imported variable name
          const importedName = node.specifiers[0]?.local?.name;

          if (importedName) {
            context.report({
              node,
              message: `Import from src/plugin-system/plugin.class should use window.PluginBaseClass`,
              fix(fixer) {
                return fixer.replaceText(
                  node,
                  `const ${importedName} = window.PluginBaseClass;`
                );
              },
            });
          }
        }
      },
    };
  },
};