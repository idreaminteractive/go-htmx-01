/** @type {import('postcss-load-config').Config} */
const config = {
    plugins: [
        require('postcss-nested'),
      require('autoprefixer'),
      require('cssnano')({
        preset: 'default',
    }),
    require('postcss-fail-on-warn')
    ]
  }
  
  module.exports = config