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
        return {
            // Handle HttpClient imports
            ImportDeclaration(node) {
                if (
                    node.source.value === 'src/service/http-client.service'
                ) {
                    context.report({
                        node,
                        message: 'Remove HttpClient import as fetch will be used instead',
                        fix(fixer) {
                            return fixer.remove(node);
                        },
                    });
                }
            },

            AssignmentExpression(node) {
                if (
                    node.left.type === 'MemberExpression' &&
                    node.left.property.name === '_httpClient' &&
                    node.right.type === 'NewExpression' &&
                    node.right.callee.name === 'HttpClient'
                ) {
                    context.report({
                        node,
                        message: 'Remove HttpClient assignment as fetch will be used instead',
                        fix(fixer) {
                            return fixer.remove(node.parent); // Remove the entire statement
                        },
                    });
                }
            },
            CallExpression(node) {
                if (
                    node.callee.type === 'MemberExpression' &&
                    node.callee.object.type === 'MemberExpression' &&
                    node.callee.object.property.name === '_httpClient'
                ) {
                    const sourceCode = context.getSourceCode();
                    const method = node.callee.property.name;

                    if (method === 'get') {
                        const [urlArg, callbackFn] = node.arguments;
                        
                        if (!urlArg || !callbackFn || callbackFn.type !== 'ArrowFunctionExpression') {
                            return;
                        }

                        const callbackBody = sourceCode.getText(callbackFn.body);
                        const callbackParamName = callbackFn.params[0].name;
                        
                        const fetchCode = `fetch(${sourceCode.getText(urlArg)})
    .then(response => response.text())
    .then(${callbackParamName} => {
        ${callbackBody.replace(/^\{|\}$/g, '').trim()}
    })`;

                        context.report({
                            node,
                            message: 'Use fetch API instead of _httpClient.get',
                            fix(fixer) {
                                return fixer.replaceText(node, fetchCode);
                            },
                        });
                    } else if (method === 'post') {
                        const [urlArg, dataArg, callbackFn, contentTypeArg] = node.arguments;
                        
                        if (!urlArg || !dataArg || !callbackFn || callbackFn.type !== 'ArrowFunctionExpression') {
                            return;
                        }

                        const callbackBody = sourceCode.getText(callbackFn.body);
                        const contentType = contentTypeArg ? sourceCode.getText(contentTypeArg) : "'application/json'";
                        const callbackParamName = callbackFn.params[0].name;
                        
                        const fetchCode = `fetch(${sourceCode.getText(urlArg)}, {
    method: 'POST',
    headers: {
        'Content-Type': ${contentType}
    },
    body: ${sourceCode.getText(dataArg)}
})
    .then(response => response.text())
    .then(${callbackParamName} => {
        ${callbackBody.replace(/^\{|\}$/g, '').trim()}
    })`;

                        context.report({
                            node,
                            message: 'Use fetch API instead of _httpClient.post',
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