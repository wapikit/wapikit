module.exports = {
	arrowParens: 'avoid',
	trailingComma: 'none',
	tabWidth: 4,
	printWidth: 100,
	semi: false,
	useTabs: true,
	bracketSpacing: true,
	bracketSameLine: false,
	singleQuote: true,
	plugins: [
		'prettier-plugin-tailwindcss'
	],
	tailwindConfig: './tailwind.config.ts'
}
