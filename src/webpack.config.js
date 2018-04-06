const path = require('path')
const webpack = require('webpack')

process.env.PROJECT_NAME = 'wtf'
process.env.PROJECT_USER = 'blunt1337'
process.env.PROJECT_REPO = 'wtfcmd'
process.env.PROJECT_DESCRIPTION = 'super alias your commands'
process.env.PROJECT_COLOR = '#508cb6'

//-----------------------------
//-- Base webpack config
//-----------------------------
const cfg = module.exports = {
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
		rules: [],
	},
	stats: { colors: true },
	resolve: {
		extensions: [],
		alias: {
			'@': path.resolve(__dirname),
			node_modules: path.resolve(__dirname, 'node_modules'),
		},
	},
	plugins: [],
}

//-----------------------------
//-- Loaders
//-----------------------------

// JS
cfg.resolve.extensions.push('.js')
cfg.module.rules.push({
	test: /\.js?$/,
	use: [{
		loader: 'babel-loader',
		options: {
			presets: ['es2015'],
			plugins: ['syntax-dynamic-import', 'syntax-object-rest-spread'],
		},
	}],
})

// Eval js
cfg.resolve.extensions.push('.eval.js')
cfg.module.rules.push({
	test: /\.eval\.js?$/,
	use: [
		...cfg.module.rules[0].use,
		'val-loader'
	],
	enforce: 'post',
})

// Vuejs + svg inline
const { VueLoaderPlugin } = require('vue-loader')
cfg.resolve.extensions.push('.vue')
cfg.module.rules.push({
	test: /\.vue$/,
	use: ['vue-loader', 'markup-inline-loader'],
})
cfg.plugins.push(new VueLoaderPlugin())

// Scss + css
cfg.resolve.extensions.push('.scss', '.css')
cfg.module.rules.push({
	test: /\.s?css$/,
	exclude: /critical\.scss$/,
	use: [
		'vue-style-loader',
		'css-loader',
		{
			loader: 'sass-loader',
			options: {
				data: '$primary: ' + process.env.PROJECT_COLOR + ';'
			},
		}
	],
})

// Critical css
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
cfg.module.rules.push({
	test: /critical\.scss$/,
	use: [
		MiniCssExtractPlugin.loader,
		'css-loader',
		{
			loader: 'sass-loader',
			options: {
				data: '$primary: ' + process.env.PROJECT_COLOR + ';'
			},
		}
	],
})
cfg.plugins.push(new MiniCssExtractPlugin({
	filename: '[name].css',
}))

// Vue markdown + svg inline
cfg.resolve.extensions.push('.md')
cfg.module.rules.push({
	test: /\.md$/,
	use: [
		'vue-loader',
		'markup-inline-loader',
		{
			loader: path.resolve(__dirname, 'node_modules/vue-markdown-loader/lib/markdown-compiler.js'),
			options: require('./markdown.config.js'),
		}
	],
})

// Assets
cfg.module.rules.push({
	test: /\.(png|woff|woff2|eot|ttf|svg)$/,
	loader: 'url-loader',
	options: {
		limit: 8192,
		name: '[name].[ext]?[hash]',
	},
})

// ESLint
cfg.module.rules.push({
	test: /\.js$/,
	enforce: 'pre',
	exclude: /node_modules/,
	use: 'eslint-loader',
})

//-----------------------------
//-- Extra plugins
//-----------------------------

// HTML pages
const HtmlWebpackPlugin = require('html-webpack-plugin')
const html_options = {
	title: process.env.PROJECT_NAME,
	description: process.env.PROJECT_DESCRIPTION,
	color: process.env.PROJECT_COLOR,
	template: 'index.ejs',
	hash: true,
	minify: {
		collapseWhitespace: true,
	},
	excludeChunks: ['critical'],
}
cfg.plugins.push(new HtmlWebpackPlugin(html_options))

// 404 page, same as index
html_options.filename = '404.html'
cfg.plugins.push(new HtmlWebpackPlugin(html_options))

// Define all envs inside js scripts
const defines = {}
for (let [key, value] of Object.entries(process.env)) {
	defines['process.env.' + key] = JSON.stringify(value)
}
cfg.plugins.push(new webpack.DefinePlugin(defines))

//-----------------------------
//-- Icons generation
//-----------------------------
const fs = require('fs')
const logo_path = path.resolve(__dirname, 'src/logo.svg')
if (fs.existsSync(logo_path)) {
	const buildIcons = require('./build-icons')
	const editSvg = svg => (svg + '').replace(/class="fill-primary"/g, 'style="fill:' + process.env.PROJECT_COLOR + '"')
	
	// Build icons now
	buildIcons(logo_path, cfg.output.path, editSvg)

	// Build with watcher
	cfg.plugins.push({
		apply(compiler) {
			compiler.plugin('watch-run', () => {
				fs.watchFile(logo_path, () => buildIcons(logo_path, cfg.output.path, editSvg))
			})
		}
	})
}

if (process.env.NODE_ENV === 'production') {
	//-----------------------------
	//-- Production addons
	//-----------------------------
	cfg.mode = 'production'
	
	// Compress js
	const CompressPlugin = require('uglifyjs-webpack-plugin')
	cfg.plugins.push(new CompressPlugin())
	
	// Improve speed
	const OptimizeJsPlugin = require('optimize-js-plugin')
	cfg.plugins.push(new OptimizeJsPlugin({
		sourceMap: false
	}))
	
	// Compress/improve css
	const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin')
	const cssnano = require('cssnano')
	cfg.plugins.push(new OptimizeCSSAssetsPlugin({
		cssProcessor: cssnano,
		cssProcessorOptions: {
			discardComments: { removeAll: true },
			autoprefixer: { add: true, browsers: ['> 5%'] },
			zindex: false,
		},
		canPrint: false,
	}))
	
	// Security
	const SriPlugin = require('webpack-subresource-integrity')
	cfg.plugins.push(new SriPlugin({
		hashFuncNames: ['sha256', 'sha384'],
		enabled: true,
	}))
} else {
	//-----------------------------
	//-- Development addons
	//-----------------------------
	cfg.mode = 'development'
	cfg.devtool = 'cheap-module-eval-source-map'
}