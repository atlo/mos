const fs = require('fs');
const paths = require('./paths');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');
const SVGSpritePlugin = require('./plugins/svg-sprite-plugin');

const pages = fs.readdirSync(paths.pages).map(file => file.replace(/\.pug/, ''));

module.exports = {
    entry: ['babel-polyfill', './src/js/main.js'],
    output: {
        filename: 'main.bundle.js',
        path: paths.build
    },
    stats: {
        modules: false,
    },
    module: {
        rules: [
            {
                test: /\.js$/g,
                loader: 'babel-loader',
                options: {
                    presets: ['env'],
                    plugins: [
                        'transform-async-to-generator',
                        'syntax-async-functions',
                        ['transform-runtime', {
                            helpers: false,
                            polyfill: false,
                            regenerator: true,
                            moduleName: 'babel-runtime'
                        }]
                    ]
                }
            },
            {
                test: /\.scss$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    'css-loader',
                    'sass-loader',
                ],
            },
            {
                test: /\.pug$/,
                loader: 'pug-loader'
            },
            {
                test: /\.(woff|woff2|eot|ttf|otf)$/,
                loader: 'file-loader',
                options: {
                    name: '[name].[ext]',
                    useRelativePath: true,
                }
            }
        ]
    },
    plugins: [
        new CleanWebpackPlugin('**/*.*', {
            root: paths.build,
        }),
        new MiniCssExtractPlugin({
            filename: '[name].css',
        }),
        ...pages.map(page =>
            new HtmlWebpackPlugin({
              template: `src/views/pages/${page}.pug`,
              filename: `${page}.html`,
              inject: 'body',
            }),
        ),
        new CopyWebpackPlugin(
          [{
              from: paths.srcImages,
              to: paths.distImages,
          }]
        ),
        new SVGSpritePlugin({
            path: paths.srcSvg,
        }),
    ]
};
