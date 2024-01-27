const path = require('path');
const webpack = require("webpack");
const Visualizer = require('webpack-visualizer-plugin2');
const NodePolyfillPlugin = require("node-polyfill-webpack-plugin");

const { statSync: fsStat, readdirSync: fsReadDir } = require("fs");
const { config: dotEnvConfig } = require("dotenv");


setupEnv({ rootPath: __dirname });

module.exports = {
  entry: './src/index.ts',
  output: {
    filename: 'dist.build.js',
    path: path.resolve(__dirname, 'public'),
    publicPath: '/',
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      }
    ],
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
  },
  plugins: [
    new webpack.EnvironmentPlugin(
      Object.keys(process.env)
    ),
    new Visualizer({
      filename: './hidden.stats.html'
    }),
    new NodePolyfillPlugin({
      includeAliases: ['buffer', 'util', 'events']
    })
  ],
  mode: process.env.NODE_ENV || "development",
  devtool: 'inline-source-map',
  optimization: {
    usedExports: true,
  },
  devServer: {
    hot: true,
    port: process.env.HTTP_PORT || 8080,
    host: '0.0.0.0',
    historyApiFallback: true,
    static: {
      directory: path.join(__dirname, 'public'),
    },
    allowedHosts: "all",
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, PATCH, OPTIONS",
      "Access-Control-Allow-Headers": "X-Requested-With, content-type, Authorization"
    },
  },
};

function setupEnv({ rootPath }) {
  const envpath = path.resolve(rootPath, './env');
  const stat = fsStat(envpath, { throwIfNoEntry: false }); // just checking if folder exists

  if (typeof stat === "undefined") {
    return console.log("no env directory");
  }

  if (!stat.isDirectory()) {
    return console.log("no env directory");
  }

  const isClient = /.*-client-.*/;
  const isEnv = /.*\.env/;
  const files = fsReadDir(envpath);

  files.forEach((file) => {
    if (!isClient.test(file)) return;
    if (!isEnv.test(file)) return;
    const filePath = path.resolve(envpath, file);
    const vars = dotEnvConfig({ path: filePath });
  });
}