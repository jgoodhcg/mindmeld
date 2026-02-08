/**
 * Screenshot Script
 *
 * Takes screenshots of specified pages. Useful for quick visual validation.
 *
 * Usage:
 *   npm run screenshot                    # Screenshot homepage
 *   npm run screenshot -- /lobbies/ABC123 # Screenshot specific path
 *   npm run screenshot -- --full          # Full page screenshot
 */

import { chromium } from '@playwright/test';
import { mkdir } from 'fs/promises';
import { join } from 'path';

const BASE_URL = process.env.BASE_URL || 'http://localhost:3000';
const SCREENSHOTS_DIR = join(import.meta.dirname, '..', 'screenshots');

async function screenshot() {
  const args = process.argv.slice(2);
  const fullPage = args.includes('--full');
  const path = args.find(a => a.startsWith('/')) || '/';

  await mkdir(SCREENSHOTS_DIR, { recursive: true });

  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1280, height: 720 } });

  const url = `${BASE_URL}${path}`;
  console.log(`Navigating to: ${url}`);

  try {
    await page.goto(url, { waitUntil: 'networkidle' });
  } catch (e) {
    console.error(`Failed to load ${url} - is the server running?`);
    await browser.close();
    process.exit(1);
  }

  const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
  const safePath = path.replace(/\//g, '_') || 'home';
  const filename = `${safePath}-${timestamp}.png`;
  const filepath = join(SCREENSHOTS_DIR, filename);

  await page.screenshot({ path: filepath, fullPage });
  console.log(`Screenshot saved: ${filepath}`);

  await browser.close();
}

screenshot();
