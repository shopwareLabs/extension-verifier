export default {
    meta: {
        type: 'suggestion',
        docs: {
            description: 'Transform HttpClient to fetch API',
            category: 'Best Practices',
            recommended: false,
        },
        fixable: 'code',
    },
    create(context) {
        let httpClientPropertyName = null;
        let httpClientVariableName = null;

        return {
            // Handle HttpClient imports
            ImportDeclaration(node) {
                if (node.source.value === 'src/service/http-client.service') {
                    context.report({
                        node,
                        message:
                            'Remove HttpClient import as fetch will be used instead',
                        fix(fixer) {
                            return fixer.remove(node);
                        },
                    });
                }
            },

            // Handle local variable declarations like: let client = new HttpClient(...)
            VariableDeclarator(node) {
                if (
                    node.init &&
                    node.init.type === 'NewExpression' &&
                    node.init.callee.name === 'HttpClient'
                ) {
                    httpClientVariableName = node.id.name;
                    context.report({
                        node: node.parent,
                        message: `Remove HttpClient assignment for '${httpClientVariableName}' as fetch will be used instead.`,
                        fix(fixer) {
                            // Remove the entire variable declaration
                            return fixer.remove(node.parent);
                        },
                    });
                }
            },

            AssignmentExpression(node) {
                if (
                    node.left.type === 'MemberExpression' &&
                    node.left.object.type === 'ThisExpression' && // Ensure it's like this.propertyName
                    node.right.type === 'NewExpression' &&
                    node.right.callee.name === 'HttpClient'
                ) {
                    httpClientPropertyName = node.left.property.name;
                    context.report({
                        node,
                        message: `Remove HttpClient assignment for '${httpClientPropertyName}' as fetch will be used instead.`,
                        fix(fixer) {
                            // node.parent should be the ExpressionStatement if the assignment is a standalone statement.
                            // This is usually correct for removing the whole line.
                            return fixer.remove(node.parent);
                        },
                    });
                }
            },
            CallExpression(node) {
                const isClassPropertyCall = 
                    httpClientPropertyName && // Ensure the property name was captured
                    node.callee.type === 'MemberExpression' &&
                    node.callee.object.type === 'MemberExpression' &&
                    node.callee.object.object.type === 'ThisExpression' && // Ensures it's this.httpClientPropertyName
                    node.callee.object.property.name === httpClientPropertyName;
                
                const isLocalVariableCall = 
                    httpClientVariableName && // Ensure the variable name was captured
                    node.callee.type === 'MemberExpression' &&
                    node.callee.object.type === 'Identifier' &&
                    node.callee.object.name === httpClientVariableName;
                
                if (isClassPropertyCall || isLocalVariableCall) {
                    const sourceCode = context.getSourceCode();
                    const method = node.callee.property.name;

                    let urlArg, dataArg, callbackFnArg, contentTypeArgText;
                    let callbackFnIndex = -1;

                    if (method === 'get' || method === 'post') {
                        urlArg = node.arguments[0];
                        if (!urlArg) return;

                        if (method === 'post') {
                            dataArg = node.arguments[1];
                            if (!dataArg) return;
                        }

                        // Find the callback function argument
                        const startIndexForCallbackSearch =
                            method === 'post' ? 2 : 1;
                        for (
                            let i = startIndexForCallbackSearch;
                            i < node.arguments.length;
                            i++
                        ) {
                            const arg = node.arguments[i];
                            if (
                                arg.type === 'ArrowFunctionExpression' ||
                                arg.type === 'FunctionExpression' ||
                                (arg.type === 'CallExpression' &&
                                    arg.callee.type === 'MemberExpression' &&
                                    arg.callee.property.name === 'bind')
                            ) {
                                callbackFnArg = arg;
                                callbackFnIndex = i;
                                break;
                            }
                        }

                        if (!callbackFnArg) return; // No suitable callback found

                        // For post, try to find contentTypeArg if it's after the callback
                        if (method === 'post') {
                            contentTypeArgText = "'application/json'"; // Default
                            if (callbackFnIndex + 1 < node.arguments.length) {
                                const potentialContentTypeArg =
                                    node.arguments[callbackFnIndex + 1];
                                if (
                                    potentialContentTypeArg.type ===
                                        'Literal' &&
                                    typeof potentialContentTypeArg.value ===
                                        'string'
                                ) {
                                    contentTypeArgText = sourceCode.getText(
                                        potentialContentTypeArg
                                    );
                                }
                            }
                        }

                        let fetchCode;
                        const urlText = sourceCode.getText(urlArg);
                        const callbackFnText =
                            sourceCode.getText(callbackFnArg);

                        // Try to inline simple arrow functions
                        if (
                            callbackFnArg.type === 'ArrowFunctionExpression' &&
                            callbackFnArg.params.length <= 1
                        ) {
                            const callbackParamName = callbackFnArg.params[0]
                                ? sourceCode.getText(callbackFnArg.params[0])
                                : '_response';
                            let callbackBodyText = sourceCode.getText(
                                callbackFnArg.body
                            );

                            if (callbackFnArg.body.type === 'BlockStatement') {
                                callbackBodyText = callbackBodyText
                                    .replace(/^\\{|\\}$/g, '')
                                    .trim();
                            } else {
                                // If not a block statement, it's an expression, so ensure it's returned in the new block
                                callbackBodyText = `return ${callbackBodyText};`;
                            }

                            const thenClause = `.then((${callbackParamName}) => {
        ${callbackBodyText}
    })`;

                            if (method === 'get') {
                                fetchCode = `fetch(${urlText})\n    .then(response => response.text())\n    ${thenClause}`;
                            } else {
                                // post
                                const dataText = sourceCode.getText(dataArg);
                                fetchCode = `fetch(${urlText}, {\n    method: 'POST',\n    headers: {\n        'Content-Type': ${contentTypeArgText}\n    },\n    body: ${dataText}\n})\n    .then(response => response.text())\n    ${thenClause}`;
                            }
                        } else {
                            // For FunctionExpression, .bind() calls, or complex ArrowFunctions
                            const thenClause = `.then(${callbackFnText})`;
                            if (method === 'get') {
                                fetchCode = `fetch(${urlText})\n    .then(response => response.text())\n    ${thenClause}`;
                            } else {
                                // post
                                const dataText = sourceCode.getText(dataArg);
                                fetchCode = `fetch(${urlText}, {\n    method: 'POST',\n    headers: {\n        'Content-Type': ${contentTypeArgText}\n    },\n    body: ${dataText}\n})\n    .then(response => response.text())\n    ${thenClause}`;
                            }
                        }

                        context.report({
                            node,
                            message: isClassPropertyCall 
                                ? `Use fetch API instead of this.${httpClientPropertyName}.${method}`
                                : `Use fetch API instead of ${httpClientVariableName}.${method}`,
                            fix(fixer) {
                                return fixer.replaceText(node, fetchCode);
                            },
                        });
                    }
                }
            },
        };
    },
};
