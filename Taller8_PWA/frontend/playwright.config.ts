import { defineConfig } from '@playwright/test';
export default defineConfig({
  testDir: './tests',
  use: { baseURL: 'http://localhost:4178', channel: 'msedge', viewport: { width: 375, height: 812 } },
  webServer: { command: 'node scripts/test-server.mjs', url: 'http://localhost:4178', reuseExistingServer: false },
});
