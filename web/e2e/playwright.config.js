// @ts-check
const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
  testDir: '.',
  timeout: 30_000,
  use: {
    baseURL: 'http://localhost:8931',
  },
  webServer: {
    command: 'python3 -m http.server 8931 --directory ..',
    url: 'http://localhost:8931',
    reuseExistingServer: true,
  },
});
