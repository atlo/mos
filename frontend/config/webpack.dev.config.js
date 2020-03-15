const fs = require('fs');
const paths = require('./paths');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const StyleLintPlugin = require('stylelint-webpack-plugin');
const SVGSpritePlugin = require('./plugins/svg-sprite-plugin');

const pages = fs.readdirSync(paths.pages).map(file => file.replace(/\.pug/, ''));

module.exports = {
    devServer: {
        contentBase: paths.assets,
        compress: true,
        port: 3001,
        open: true,
        overlay: {
          errors: true
        },
        stats: {
          modules: false,
        }
    },
    entry: './src/js/main.js',
    output: {
        filename: 'main.bundle.js',
    },
    module: {
        rules: [
            {
                test: /\.js$/,
                loader: 'babel-loader',
                exclude: /node_modules/,
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
                test: /\.js$/,
                exclude: /node_modules/,
                loader: "eslint-loader",
            },
            {
                test: /\.scss$/,
                use: [
                    { loader: "style-loader" }, 
                    { 
                        loader: "css-loader",
                        options: {
                            sourceMap: true
                        }
                    }, 
                    { 
                        loader: "sass-loader",
                        options: {
                            sourceMap: true
                        }
                    },
                ]
            },
            {
                test: /\.pug$/, 
                loader: 'pug-loader'
            },
            { 
                test: /\.(woff|woff2|eot|ttf|otf)$/,
                loader: 'url-loader?limit=100000'
            },
            {
                test: /\.(jpe?g|png|gif)$/,
                loader: 'file-loader',
            },
        ]
    },
    plugins: [
      ...pages.map(page =>
        new HtmlWebpackPlugin({
            template: `src/views/pages/${page}.pug`,
            filename: `${page}.html`,
            inject: 'body',
        })
      ),
      new StyleLintPlugin(),
      new SVGSpritePlugin({
        path: paths.srcSvg,
      }),
    ]
};
