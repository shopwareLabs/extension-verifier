export default {
    meta: {
        type: 'suggestion',
        docs: {
            description: 'Transform querystring.parse to URLSearchParams',
            category: 'Modernization',
            recommended: true,
        },
        fixable: 'code',
        schema: [],
    },
    create(context) {
        return {
            ImportDeclaration(node) {
                if (node.source.value === 'query-string') {
                    context.report({
                        node,
                        message: 'Remove querystring import as URLSearchParams is used instead',
                        fix(fixer) {
                            return fixer.remove(node);
                        },
                    });
                }
            },

            CallExpression(node) {
                if (
                    node.callee.type === 'MemberExpression' &&
                    node.callee.object.name === 'querystring' &&
                    node.callee.property.name &&
                    (node.callee.property.name === 'parse' || node.callee.property.name === 'stringify') &&
                    node.arguments.length > 0
                ) {
                    const sourceCode = context.getSourceCode();
                    const argumentSource = sourceCode.getText(node.arguments[0]);

                    if (node.callee.property.name === 'parse') {
                        context.report({
                            node,
                            message: 'Use URLSearchParams instead of querystring.parse',
                            fix(fixer) {
                                return fixer.replaceText(
                                    node,
                                    `Object.fromEntries(new URLSearchParams(${argumentSource}).entries())`
                                );
                            },
                        });
                    } else {
                        context.report({
                            node,
                            message: 'Use URLSearchParams instead of querystring.stringify',
                            fix(fixer) {
                                return fixer.replaceText(
                                    node,
                                    `new URLSearchParams(${argumentSource}).toString()`
                                );
                            },
                        });
                    }
                }
            },
        };
    },
};