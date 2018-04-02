const path = require('path')
const common = require('./webpack.common.js')
const merge = require('webpack-merge')
const webpack = require('webpack')
const CompressPlugin = require('uglifyjs-webpack-plugin')
const OptimizeJsPlugin = require('optimize-js-plugin')
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin')
const cssnano = require('cssnano')
const glob = require('glob')
const SriPlugin = require('webpack-subresource-integrity')

module.exports = merge.smart({
	module: {
		rules: [
			{
				test: /\.vue$/,
				loader: 'vue-loader',
				options: {
					preLoaders: {
						html: 'markup-inline-loader'
					},
				},
			},
		],
	},
	mode: 'production',
	plugins: [
		new webpack.DefinePlugin({
			'process.env': {
				NODE_ENV: '"production"'
			}
		}),
		new CompressPlugin(),
		new OptimizeJsPlugin({
			sourceMap: false
		}),
		new OptimizeCSSAssetsPlugin({
			cssProcessor: cssnano,
			cssProcessorOptions: {
				discardComments: { removeAll: true },
				autoprefixer: { add: true, browsers: ['> 5%'] },
				zindex: false,
			},
			canPrint: false,
		}),
		new SriPlugin({
			hashFuncNames: ['sha256', 'sha384'],
			enabled: true,
		}),
	]
}, common)