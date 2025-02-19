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

    function getBaseIndent(node) {
      const text = sourceCode.getText(node);
      const match = text.match(/^\s*/);
      return match ? match[0] : "";
    }

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

      VariableDeclaration(node) {
        if (
          node.kind === "const" &&
          node.declarations.length === 1 &&
          node.declarations[0].id.type === "ObjectPattern" &&
          node.declarations[0].init &&
          node.declarations[0].init.type === "CallExpression" &&
          node.declarations[0].init.callee.type === "MemberExpression"
        ) {
          const callee = node.declarations[0].init.callee;
          if (
            callee.object.name === "Component" &&
            callee.property.name === "getComponentHelper"
          ) {
            context.report({
              node,
              message:
                "Remove the unused Component.getComponentHelper() import.",
              fix(fixer) {
                // Remove the entire variable declaration.
                return fixer.remove(node);
              },
            });
          }
        }
      },

      // 2. Transform spread elements of mapState.
      SpreadElement(node) {
        const arg = node.argument;
        if (
          arg &&
          arg.type === "CallExpression" &&
          arg.callee &&
          (arg.callee.name === "mapState" || arg.callee.name === "mapGetters") &&
          arg.arguments.length === 2 &&
          arg.arguments[0].type === "Literal" &&
          (arg.arguments[1].type === "ArrayExpression" ||
            arg.arguments[1].type === "ObjectExpression")
        ) {
          const storeName = arg.arguments[0].value;
          const baseIndent = getBaseIndent(node);
          let computedText = "";

          // a. Handle array syntax:
          if (arg.arguments[1].type === "ArrayExpression") {
            const props = arg.arguments[1].elements
              .filter(
                (el) =>
                  el &&
                  el.type === "Literal" &&
                  typeof el.value === "string"
              )
              .map((el) => el.value);

            computedText = props
              .map((prop) => {
                return (
                  "\n" +
                  baseIndent +
                  `${prop}() {\n` +
                  baseIndent +
                  "    " +
                  `return Shopware.Store.get('${storeName}').${prop};\n` +
                  baseIndent +
                  "}" + (prop !== props[props.length - 1] ? "," : "")
                );
              })
              .join("");
          }
          // b. Handle object syntax.
          else if (arg.arguments[1].type === "ObjectExpression") {
            computedText = arg.arguments[1].properties
              .filter((prop) => prop.type === "Property")
              .map((prop) => {
                // Determine the computed property name.
                let propName = "";
                if (prop.key.type === "Identifier") {
                  propName = prop.key.name;
                } else if (prop.key.type === "Literal") {
                  propName = prop.key.value;
                } else {
                  return "";
                }

                // Case 1: property value is a literal string (mapping).
                if (
                  prop.value.type === "Literal" &&
                  typeof prop.value.value === "string"
                ) {
                  const mappedProp = prop.value.value;
                  return (
                    "\n" +
                    baseIndent +
                    `${propName}() {\n` +
                    baseIndent +
                    "    " +
                    `return Shopware.Store.get('${storeName}').${mappedProp};\n` +
                    baseIndent +
                    "},"
                  );
                }
                // Case 2: property value is a function.
                else if (
                  prop.value.type === "ArrowFunctionExpression" ||
                  prop.value.type === "FunctionExpression"
                ) {
                  const fn = prop.value;
                  if (!fn.params || fn.params.length === 0) return "";
                  if (fn.params[0].type !== "Identifier") return "";
                  const paramName = fn.params[0].name;

                  let returnedExpr = null;
                  if (fn.body.type !== "BlockStatement") {
                    // Arrow function with implicit return.
                    returnedExpr = fn.body;
                  } else {
                    // Function with block body — find the return statement.
                    const retStmt = fn.body.body.find(
                      (n) => n.type === "ReturnStatement" && n.argument
                    );
                    if (retStmt) {
                      returnedExpr = retStmt.argument;
                    }
                  }
                  if (!returnedExpr) return "";
                  const returnedText = sourceCode.getText(returnedExpr);
                  // Ensure returned text starts with the parameter.
                  const regex = new RegExp("^" + paramName + "\\b");
                  if (!regex.test(returnedText)) return "";
                  const newReturnedText = returnedText.replace(
                    regex,
                    `Shopware.Store.get('${storeName}')`
                  );

                  return (
                    "\n" +
                    baseIndent +
                    `${propName}() {\n` +
                    baseIndent +
                    "    " +
                    `return ${newReturnedText};\n` +
                    baseIndent +
                    "},"
                  );
                }
                return "";
              })
              .join("");

              if (computedText.endsWith(",")) {
                computedText = computedText.slice(0, -1);
              }
          }

          context.report({
            node,
            message:
              "Replace spread mapState call with explicit computed property definitions.",
            fix(fixer) {
              return fixer.replaceText(node, computedText);
            },
          });
        }
      },

      CallExpression(node) {
        if (
          node.callee &&
          node.callee.type === "MemberExpression" &&
          node.callee.property &&
          node.callee.property.name === "commit"
        ) {
          // Allow two cases:
          // a) Shopware.State.commit(...)
          // b) State.commit(...)
          const stateMember = node.callee.object;
          let isValid = false;
          let shortHand = false;

          if (
            stateMember.type === "MemberExpression" &&
            stateMember.object &&
            stateMember.object.type === "Identifier" &&
            stateMember.object.name === "Shopware" &&
            stateMember.property &&
            stateMember.property.name === "State"
          ) {
            isValid = true;
          } else if (
            stateMember.type === "Identifier" &&
            stateMember.name === "State"
          ) {
            isValid = true;
            shortHand = true;
          }

          if (!isValid) return;

          if (node.arguments.length < 1) return;
          const firstArg = node.arguments[0];
          if (
            firstArg.type === "Literal" &&
            typeof firstArg.value === "string"
          ) {
            const parts = firstArg.value.split("/");
            if (parts.length === 2) {
              const [storeName, methodName] = parts;
              const args = node.arguments.slice(1);
              const argsText = args
                .map((arg) => sourceCode.getText(arg))
                .join(", ");
              const newCode = `${shortHand ? '' : 'Shopware.'}Store.get('${storeName}').${methodName}(${argsText ? argsText : ""
                })`;
              context.report({
                node,
                message:
                  "Replace State.commit/Shopware.State.commit call with Shopware.Store.get call.",
                fix(fixer) {
                  return fixer.replaceText(node, newCode);
                },
              });
            }
          }
        }
      },
    };
  }
};
