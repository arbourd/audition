const HtmlWebpackPlugin = require('html-webpack-plugin');
var CopyWebpackPlugin = require('copy-webpack-plugin');
const webpack = require('webpack');
const path = require('path');

module.exports = {
	entry: './index.js',
	output: {
		path: path.resolve(__dirname, 'dist'),
		filename: 'bundle.js'
	},
	module: {
		loaders: [{
			test: /\.js$/,
			exclude: /node_modules/,
			loader: 'babel-loader'
		}]
	},
	plugins: [
        new HtmlWebpackPlugin({
			template: './index.html',
			assets: {
				style: 'style.[hash].css',
			}
		}),
        new CopyWebpackPlugin([{
			from: 'styles.css',
			to: 'css/styles.css'
		}, {
			from: 'node_modules/bulma/css/bulma.css',
			to: 'css/bulma.css'
		}])
    ]
}
