module.exports = {
	root: true,
	env: {
		browser: true,
		node: true,
		es6: true,
	},
	parserOptions: {
		parser: 'babel-eslint'
	},
	extends: [
		// https://github.com/vuejs/eslint-plugin-vue#priority-a-essential-error-prevention
		// consider switching to `plugin:vue/strongly-recommended` or `plugin:vue/recommended` for stricter rules.
		'plugin:vue/essential',
		'eslint:recommended'
	],
	// required to lint *.vue files
	plugins: [
		'vue'
	],
	// add your custom rules here
	rules: {
		'no-console': 'off',
		'comma-dangle': 'off',
		'indent': ['error', 'tab', { SwitchCase: 1 }],
		'no-tabs': 'off',
		'camelcase': 'off',
		'space-before-function-paren': ['error', { anonymous: 'always', named: 'never', asyncArrow: 'always' }],
		'no-trailing-spaces': ['error', { skipBlankLines: true }],
	},
}