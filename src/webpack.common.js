const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const genIcons = require('./generate.icons.js')

const PROJECT_NAME = 'wtf'
const PROJECT_USER = 'blunt1337'
const PROJECT_REPO = 'wtfcmd'
const PROJECT_DESCRIPTION = 'super alias your commands'
const PROJECT_COLOR = '#508cb6'

// Build icons
genIcons(PROJECT_COLOR)

const htmlPluginOptions = {
	title: PROJECT_NAME,
	description: PROJECT_DESCRIPTION,
	color: PROJECT_COLOR,
	template: 'index.ejs',
	hash: true,
	minify: {
		collapseWhitespace: true,
	},
	excludeChunks: ['critical'],
}

module.exports = {
	entry: {
		index: path.resolve(__dirname, 'js/index.js'),
		critical: path.resolve(__dirname, 'sass/critical.scss'),
	},
	output: {
		path: path.resolve(__dirname, '..'),
		filename: '[name].js',
		crossOriginLoading: 'anonymous',
	},
	module: {
		rules: [
			{
				test: /\.eval\.js?$/,
				use: [{
					loader: 'babel-loader',
					options: {
						presets: ['es2015'],
						plugins: ['syntax-dynamic-import', 'syntax-object-rest-spread'],
					}
				}, 'val-loader'],
			},
			{
				test: /\.js?$/,
				use: {
					loader: 'babel-loader',
					options: {
						presets: ['es2015'],
						plugins: ['syntax-dynamic-import', 'syntax-object-rest-spread'],
					}
				},
			},
			{
				test: /critical\.scss$/,
				use: [
					MiniCssExtractPlugin.loader,
					'css-loader',
					{
						loader: 'sass-loader',
						options: {
							data: '$primary: ' + PROJECT_COLOR + ';'
						},
					}
				],
			},
			{
				test: /\.vue$/,
				loader: 'vue-loader',
				options: {
					loaders: {
						html: 'markup-inline-loader',
						sass: ['vue-style-loader', 'css-loader', {
							loader: 'sass-loader',
							options: {
								data: '$primary: ' + PROJECT_COLOR + ';'
							},
						}],
					},
				}
			},
			{
				test: /\.(png|jpg|gif|svg)$/,
				loader: 'file-loader',
				options: {
					name: '[name].[ext]?[hash]',
				}
			},
			{
				test: /\.md?$/,
				use: {
					loader: 'vue-markdown-loader',
					options: require('./markdown.config.js'),
				},
			},
		]
	},
	stats: { colors: true },
	resolve: {
		extensions: ['.js', '.vue', '.scss', '.md', '.html', '.eval.js'],
		alias: {
			js: path.resolve(__dirname, 'js'),
			pages: path.resolve(__dirname, 'pages'),
			sass: path.resolve(__dirname, 'sass'),
		},
	},
	plugins: [
		new MiniCssExtractPlugin({
			filename: '[name].css',
		}),
		new HtmlWebpackPlugin(htmlPluginOptions),
		new HtmlWebpackPlugin({
			...htmlPluginOptions,
			filename: '404.html',
		}),
		new webpack.DefinePlugin({
			PROJECT_NAME: JSON.stringify(PROJECT_NAME),
			PROJECT_USER: JSON.stringify(PROJECT_USER),
			PROJECT_REPO: JSON.stringify(PROJECT_REPO),
			PROJECT_DESCRIPTION: JSON.stringify(PROJECT_DESCRIPTION),
			PROJECT_COLOR: JSON.stringify(PROJECT_COLOR),
		}),
	],
}