import { defineConfig } from '@playwright/test';

const testPort = Number(process.env.PWA_TEST_PORT ?? 4179);

export default defineConfig({
  testDir: './tests',
  use: { baseURL: `http://localhost:${testPort}`, channel: 'msedge', viewport: { width: 375, height: 812 } },
  webServer: { command: `node scripts/test-server.mjs ${testPort}`, url: `http://localhost:${testPort}`, reuseExistingServer: false },
});
