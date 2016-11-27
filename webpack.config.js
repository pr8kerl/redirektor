module.exports = {
  resolve: {
    alias: {
      'vue$': 'vue/dist/vue.common.js'
    }
  },
  entry: './public/js/app.js',
  output: {
    filename: 'bundle.js',
    path: './public/js'
  }
}
