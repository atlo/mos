const glob = require('glob');
const path = require('path');
const SVGSpriter = require('svg-sprite');
const fs = require('fs');

let spriter;

const compileSVG = function (srcPath, callback) {
  glob(srcPath, null, (err, files) => {
    if (err) {
      throw err;
    } else {
      spriter = new SVGSpriter({
        dest: '../public',
        shape: {
          id: {
            separator: '--',
            generator: (err, file) => (file.basename.replace('.svg', '')),
            pseudo: '~'
          },
          dimension: {
            maxWidth: 2000,
            maxHeight: 2000,
            precision: 2,
            attributes: false,
          },
          spacing: {
            padding : 0,
            box: 'content'
          },
          transform: [{
            svgo: {
              plugins: [
                { convertStyleToAttrs: true },
                { removeUnknownsAndDefaults: true },
                { removeEditorsNSData: true },
                { removeTitle: true },
                { convertColors: true },
                { removeUnusedNS: true },
                { removeComments: true },
                { removeEmptyContainers: true },
                { cleanupEnableBackground: true },
                { minifyStyles: true },
                { removeUselessStrokeAndFill: true },
                { sortAttrs: true },
                { removeStyleElement: false },
                { cleanupIDs: true },
                { removeDoctype: true },
                { removeViewBox: true }
              ]
            }
          }],
          meta: null,
          align: null,
          dest: null
        },
        svg: {
          xmlDeclaration: false,
          doctypeDeclaration: false,
          namespaceIDs: false,
          namespaceClassnames: false,
          dimensionAttributes: false
        },
        mode: {
          symbol: true
        }
      });

      files.forEach(function(fileName) {
        const fileContent = fs.readFileSync(path.resolve(fileName));
        if (fileContent) {
          spriter.add(fileName, null, fileContent);
        }
      });

      spriter.compile(function(err, result) {
        if (err) {
          throw err;
        } else {
          let svgBuffer;
          for(let mode in result) {
            if (((result[mode] || {}).sprite || {}).contents) {
              svgBuffer = result[mode].sprite.contents;
            }
          }

          let svgFileContent = svgBuffer.toString();
          callback(svgFileContent);
        }
      });
    }
  });
};

function SVGSpritePlugin(options) {
  this.options = options;
}

SVGSpritePlugin.prototype.apply = function(compiler) {
  const self = this;
  if (compiler.hooks) {
    compiler.hooks.compilation.tap('SVGSpritePlugin', (compilation) => {
      compilation.hooks.htmlWebpackPluginAfterHtmlProcessing.tapAsync(
        'SVGSpritePlugin',
        (data, callback) => {
          compileSVG(`${self.options.path}/**/*.svg`, (svgSpriteContent) => {
            data.html += svgSpriteContent;
            callback(null, data);
          });
        }
      )
    });
  }
};

module.exports = SVGSpritePlugin;
