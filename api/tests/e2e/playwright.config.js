// @ts-check
const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
  testDir: '.',
  timeout: 10000,
  use: {
    baseURL: 'file:///Users/kuimovmihail/Downloads/yarik-varit/yarik-varit/index.html',
    headless: true,
  },
});