/** @type {import('eslint').ESLint.ConfigData} */
const eslintConfig = [
  {
    files: ['**/*.js'],
    ignores: ['prettier.config.js', 'eslint.config.js'],
    rules: {
      'prefer-const': 'off',
      'no-use-before-define': 'error',
    },
  },
];

export default eslintConfig;
