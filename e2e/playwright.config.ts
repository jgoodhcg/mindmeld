import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  timeout: 30000,
  retries: 0,

  use: {
    baseURL: process.env.BASE_URL || 'http://localhost:3000',
    headless: true,
    screenshot: 'on',
    trace: 'retain-on-failure',
  },

  outputDir: './results',

  projects: [
    {
      name: 'chromium',
      use: { browserType: 'chromium' },
    },
  ],
});
