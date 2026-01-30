/**
 * Flow Script
 *
 * Drives UI interactions and captures screenshots at key checkpoints.
 * Designed for agentic validation - run a flow, get visual feedback.
 *
 * Usage:
 *   npm run flow                     # Run default flow (create lobby)
 *   npm run flow -- join ABC123      # Run join lobby flow
 *   npm run flow -- trivia ABC123    # Run trivia game flow
 */

import { chromium, Page } from '@playwright/test';
import { mkdir } from 'fs/promises';
import { join } from 'path';

const BASE_URL = process.env.BASE_URL || 'http://localhost:3000';
const SCREENSHOTS_DIR = join(import.meta.dirname, '..', 'screenshots');

let screenshotIndex = 0;
let flowName = 'flow';

async function capture(page: Page, label: string) {
  const filename = `${flowName}-${String(screenshotIndex++).padStart(2, '0')}-${label}.png`;
  const filepath = join(SCREENSHOTS_DIR, filename);
  await page.screenshot({ path: filepath });
  console.log(`  [${label}] ${filepath}`);
  return filepath;
}

async function createLobbyFlow(page: Page) {
  flowName = 'create-lobby';
  console.log('\n=== Create Lobby Flow ===\n');

  await page.goto(BASE_URL, { waitUntil: 'networkidle' });
  await capture(page, 'homepage');

  // Look for create lobby button/link
  const createBtn = page.locator('text=Create').first();
  if (await createBtn.isVisible()) {
    await createBtn.click();
    await page.waitForLoadState('networkidle');
    await capture(page, 'after-create-click');
  }

  // If there's a form, try to fill it
  const nameInput = page.locator('input[name="name"], input[placeholder*="name" i]').first();
  if (await nameInput.isVisible()) {
    await nameInput.fill('TestPlayer');
    await capture(page, 'name-filled');

    const submitBtn = page.locator('button[type="submit"], input[type="submit"]').first();
    if (await submitBtn.isVisible()) {
      await submitBtn.click();
      await page.waitForLoadState('networkidle');
      await capture(page, 'form-submitted');
    }
  }

  await capture(page, 'final-state');
  console.log(`\nFinal URL: ${page.url()}`);
}

async function joinLobbyFlow(page: Page, lobbyCode: string) {
  flowName = 'join-lobby';
  console.log(`\n=== Join Lobby Flow (${lobbyCode}) ===\n`);

  await page.goto(`${BASE_URL}/lobby/${lobbyCode}`, { waitUntil: 'networkidle' });
  await capture(page, 'lobby-page');

  const nameInput = page.locator('input[name="name"], input[placeholder*="name" i]').first();
  if (await nameInput.isVisible()) {
    await nameInput.fill('TestPlayer2');
    await capture(page, 'name-filled');

    const submitBtn = page.locator('button[type="submit"], input[type="submit"]').first();
    if (await submitBtn.isVisible()) {
      await submitBtn.click();
      await page.waitForLoadState('networkidle');
      await capture(page, 'joined');
    }
  }

  await capture(page, 'final-state');
  console.log(`\nFinal URL: ${page.url()}`);
}

async function triviaFlow(page: Page, lobbyCode: string) {
  flowName = 'trivia';
  console.log(`\n=== Trivia Flow (${lobbyCode}) ===\n`);

  await page.goto(`${BASE_URL}/lobby/${lobbyCode}/trivia`, { waitUntil: 'networkidle' });
  await capture(page, 'trivia-page');

  // Look for answer buttons or inputs
  const answerBtns = page.locator('button').filter({ hasText: /^[A-D]|answer/i });
  const count = await answerBtns.count();

  if (count > 0) {
    console.log(`  Found ${count} answer buttons`);
    await answerBtns.first().click();
    await page.waitForTimeout(500);
    await capture(page, 'answer-selected');
  }

  await capture(page, 'final-state');
  console.log(`\nFinal URL: ${page.url()}`);
}

async function main() {
  await mkdir(SCREENSHOTS_DIR, { recursive: true });

  const args = process.argv.slice(2);
  const command = args[0] || 'create';
  const param = args[1];

  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1280, height: 720 } });

  try {
    switch (command) {
      case 'create':
        await createLobbyFlow(page);
        break;
      case 'join':
        if (!param) throw new Error('join requires a lobby code');
        await joinLobbyFlow(page, param);
        break;
      case 'trivia':
        if (!param) throw new Error('trivia requires a lobby code');
        await triviaFlow(page, param);
        break;
      default:
        console.error(`Unknown flow: ${command}`);
        console.error('Available: create, join <code>, trivia <code>');
        process.exit(1);
    }
  } catch (e) {
    await capture(page, 'error');
    throw e;
  } finally {
    await browser.close();
  }

  console.log('\n✓ Flow complete. Screenshots in e2e/screenshots/');
}

main();
