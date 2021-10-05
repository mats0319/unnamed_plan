module.exports = {
  root: true,
  env: {
    node: true
  },
  'extends': [
    'plugin:vue/essential',
    'eslint:recommended',
    '@vue/typescript/recommended'
  ],
  parserOptions: {
    ecmaVersion: 2020
  },
  rules: {
    'no-console': "off",
    'no-debugger': "off",
    "@typescript-eslint/no-var-requires": "off", // mainly for vue.config.js
    "@typescript-eslint/ban-ts-comment": "off"
  }
}
