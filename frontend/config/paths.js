const path = require('path');

module.exports = {
  pages: path.resolve(__dirname, '../src/views/pages'),
  assets: path.resolve(__dirname, '../src'),
  build: path.resolve(__dirname, '../build'),
  srcImages: path.resolve(__dirname, '../src/assets/images'),
  srcSvg: path.resolve(__dirname, '../src/assets/svg'),
  distImages: path.resolve(__dirname, '../build/assets/images/[name].[ext]'),
};