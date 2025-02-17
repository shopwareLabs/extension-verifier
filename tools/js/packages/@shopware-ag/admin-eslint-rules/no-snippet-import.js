export default {
    meta: {
      type: "problem",
      docs: {
        description:
          "Forbid passing `snippets` to Shopware.Module.register or Module.register",
        category: "Best Practices",
        recommended: true,
      }
    },
    create(context) {
      return {
        CallExpression(node) {
          // Ensure we are dealing with a MemberExpression call.
          const callee = node.callee;
          if (callee.type !== "MemberExpression") {
            return;
          }
  
          // Check if the method name is "register" (either as an Identifier or Literal)
          const property = callee.property;
          const isRegister =
            (property.type === "Identifier" && property.name === "register") ||
            (property.type === "Literal" && property.value === "register");
  
          if (!isRegister) {
            return;
          }
  
          // Determine if the call is from Shopware.Module.register or (destructured) Module.register.
          let isShopwareModuleRegister = false;
          if (callee.object.type === "MemberExpression") {
            // e.g. Shopware.Module.register(...)
            const objectPart = callee.object;
            if (
              objectPart.object.type === "Identifier" &&
              objectPart.object.name === "Shopware" &&
              objectPart.property.type === "Identifier" &&
              objectPart.property.name === "Module"
            ) {
              isShopwareModuleRegister = true;
            }
          } else if (callee.object.type === "Identifier") {
            // e.g. Module.register(...)
            // For simplicity we'll assume that if `Module` is used directly, it is the one from Shopware.
            if (callee.object.name === "Module") {
              isShopwareModuleRegister = true;
            }
          }
  
          if (!isShopwareModuleRegister) {
            return;
          }
  
          // Check the second argument (options object)
          const args = node.arguments;
          if (!args || args.length < 2) {
            return;
          }
          const options = args[1];
          if (options.type !== "ObjectExpression") {
            return;
          }
  
          // Look for a property key named 'snippets'
          options.properties.forEach((prop) => {
            // Skip spread elements or non-standard notations
            if (prop.type !== "Property" || prop.computed) {
              return;
            }
            const key = prop.key;
            if (
              (key.type === "Identifier" && key.name === "snippets") ||
              (key.type === "Literal" && key.value === "snippets")
            ) {
              context.report({
                node: prop,
                message:
                  "Passing 'snippets' to Shopware.Module.register is forbidden as it increases the bundle size. Snippets are automatically loaded when they are placed in a folder named snippet.",
              });
            }
          });
        },
      };
    },
  };