module.exports = {
  devServer: {
    port: 8999,
    disableHostCheck: true
  },
  pages: {
    index: {
      entry: 'src/main.ts',
      template: 'public/index.html',
      filename: 'index.html'
    },
    signin: {
      entry: 'src/main.ts',
      template: 'public/p.html',
      filename: 'p.html'
    }
  },
  css: {
    loaderOptions: {
      scss: {
        prependData: '@import "@/style.scss";'
      }
    }
  }
};
