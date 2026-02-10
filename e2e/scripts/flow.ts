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
 *   npm run flow -- coordinates CODE # Run Cluster game flow
 *   npm run flow -- templates        # Open templates modal flow
 */

import { chromium, Page, BrowserContext } from '@playwright/test';
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

  await page.goto(`${BASE_URL}/trivia`, { waitUntil: 'networkidle' });
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

  await page.goto(`${BASE_URL}/lobbies/${lobbyCode}`, { waitUntil: 'networkidle' });
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

  await page.goto(`${BASE_URL}/lobbies/${lobbyCode}`, { waitUntil: 'networkidle' });
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

async function coordinatesFlow(page: Page, lobbyCode: string) {
  flowName = 'coordinates';
  let code = lobbyCode;
  const extraPlayers: Array<{ context: BrowserContext; page: Page }> = [];
  console.log(`\n=== Coordinates Flow (${lobbyCode}) ===\n`);

  if (lobbyCode.toUpperCase() === 'NEW') {
    await page.goto(`${BASE_URL}/cluster`, { waitUntil: 'networkidle' });
    await capture(page, 'cluster-home');

    const createForm = page.locator('form[action="/lobbies"]').first();
    await createForm.waitFor({ state: 'visible', timeout: 10000 });
    await createForm.locator('input[name="name"]').fill(`Coordinates Flow ${Date.now()}`);
    await createForm.locator('input[name="nickname"]').fill('FlowHost');
    await Promise.all([
      page.waitForURL(/\/lobbies\/[A-Z0-9]{6}/, { timeout: 10000 }),
      createForm.locator('button[type="submit"]').click(),
    ]);

    const match = page.url().match(/\/lobbies\/([A-Z0-9]{6})/);
    code = match ? match[1] : lobbyCode;
    console.log(`  Created lobby code: ${code}`);

    // Add two extra players so Cluster can start (minimum 3 players).
    const browser = page.context().browser();
    if (browser) {
      const joinAs = async (nickname: string) => {
        const ctx = await browser.newContext({ viewport: { width: 1280, height: 720 } });
        const p = await ctx.newPage();
        await p.goto(`${BASE_URL}/lobbies/${code}`, { waitUntil: 'networkidle' });
        const joinInput = p.locator('input[name="nickname"]').first();
        if (await joinInput.isVisible().catch(() => false)) {
          await joinInput.fill(nickname);
          await Promise.all([
            p.waitForURL(new RegExp(`/lobbies/${code}$`), { timeout: 10000 }),
            p.locator('button:has-text("CONNECT")').first().click(),
          ]);
        }
        extraPlayers.push({ context: ctx, page: p });
      };

      await joinAs('FlowPlayer2');
      await joinAs('FlowPlayer3');
      await page.waitForTimeout(1000);
    }
  } else {
    await page.goto(`${BASE_URL}/lobbies/${code}`, { waitUntil: 'networkidle' });

    const joinInput = page.locator('input[name="nickname"]').first();
    if (await joinInput.isVisible().catch(() => false)) {
      await joinInput.fill('FlowPlayer');
      const connectBtn = page.locator('button:has-text("CONNECT")').first();
      if (await connectBtn.isVisible().catch(() => false)) {
        await Promise.all([
          page.waitForURL(new RegExp(`/lobbies/${code}$`), { timeout: 10000 }),
          connectBtn.click(),
        ]);
      }
    }
  }

  await capture(page, 'cluster-page');

  const startBtn = page.locator('form[action$="/cluster/start"] button:has-text("START CLUSTER")').first();
  if (await startBtn.isVisible().catch(() => false)) {
    await page.waitForSelector('#player-list li:nth-child(3)', { timeout: 10000 }).catch(() => {});
    await Promise.all([
      page.waitForLoadState('networkidle'),
      startBtn.click(),
    ]);
    await capture(page, 'round-started');
  }

  const plane = page.locator('#cluster-plane-input').first();
  if (await plane.isVisible().catch(() => false)) {
    const box = await plane.boundingBox();
    if (box) {
      const clickX = box.x + (box.width * 0.33);
      const clickY = box.y + (box.height * (1 - 0.77));
      await page.mouse.click(clickX, clickY);
      await capture(page, 'coordinates-selected');
    }

    const submitBtn = page.locator('button[type="submit"]', { hasText: 'SUBMIT POINT' }).first();
    if (await submitBtn.isVisible().catch(() => false)) {
      await submitBtn.click();
      await page.waitForLoadState('networkidle');
      await capture(page, 'coordinate-submitted');
    }
  }

  if (extraPlayers.length > 0) {
    const extraPoints = [
      { x: 0.82, y: 0.22 },
      { x: 0.61, y: 0.73 },
    ];

    for (let i = 0; i < extraPlayers.length; i++) {
      const p = extraPlayers[i].page;
      const pt = extraPoints[i] || { x: 0.5, y: 0.5 };
      const pPlane = p.locator('#cluster-plane-input').first();
      if (await pPlane.isVisible().catch(() => false)) {
        const box = await pPlane.boundingBox();
        if (box) {
          const clickX = box.x + (box.width * pt.x);
          const clickY = box.y + (box.height * (1 - pt.y));
          await p.mouse.click(clickX, clickY);
        }

        const pSubmit = p.locator('button[type="submit"]', { hasText: 'SUBMIT POINT' }).first();
        if (await pSubmit.isVisible().catch(() => false)) {
          await pSubmit.click();
          await p.waitForLoadState('networkidle');
        }
      }
    }

    await page.waitForTimeout(1200);
    await capture(page, 'reveal-state');

    const standingsHeader = page.locator('text=Round pts').first();
    if (await standingsHeader.isVisible().catch(() => false)) {
      await standingsHeader.scrollIntoViewIfNeeded();
      await page.waitForTimeout(200);
      await capture(page, 'reveal-standings');
    }
  }

  await capture(page, 'final-state');
  console.log(`\nFinal URL: ${page.url()}`);

  for (const extra of extraPlayers) {
    await extra.context.close();
  }
}

async function templatesFlow(page: Page) {
  flowName = 'templates';
  console.log('\n=== Templates Modal Flow ===\n');

  await page.goto(`${BASE_URL}/trivia`, { waitUntil: 'networkidle' });
  await capture(page, 'homepage');

  const createForm = page.locator('form[action="/lobbies"]').first();
  if (!(await createForm.isVisible())) {
    const triviaLink = page.locator('a[href="/trivia"], a:has-text("TRIVIA")').first();
    if (await triviaLink.isVisible()) {
      await triviaLink.click();
      await page.waitForLoadState('networkidle');
    }
  }
  await createForm.waitFor({ state: 'visible', timeout: 5000 });

  await createForm.locator('input[name="name"]').fill('Template Session');
  await createForm.locator('input[name="nickname"]').fill('TemplateHost');
  await capture(page, 'create-form-filled');

  const createSubmit = createForm.locator('button[type="submit"]').first();
  await Promise.all([
    page.waitForURL(/\/lobbies\//, { timeout: 10000 }),
    createSubmit.click(),
  ]);
  await capture(page, 'lobby-created');

  const startButton = page.locator('button', { hasText: 'START GAME' }).first();
  if (await startButton.isVisible()) {
    await Promise.all([
      page.waitForLoadState('networkidle'),
      startButton.click(),
    ]);
  }

  const submitForm = page.locator('form[action$="/questions"]').first();
  await submitForm.waitFor({ state: 'visible', timeout: 10000 });

  const templatesButton = page.locator('button', { hasText: 'Need help? Use a template' }).first();
  await templatesButton.click();

  const templatesModal = page.locator('#templates-modal');
  await templatesModal.waitFor({ state: 'visible', timeout: 5000 });
  await page.locator('#templates-content').waitFor({ state: 'visible', timeout: 10000 });

  const categoryHeaders = await page.locator('#templates-content h3').allTextContents();
  if (categoryHeaders.length > 0) {
    console.log(`  Categories: ${categoryHeaders.join(' | ')}`);
  }

  const namePlaceholderCount = await page.locator('#templates-content button p', { hasText: '[my name]' }).count();
  console.log(`  "[my name]" placeholders: ${namePlaceholderCount}`);

  await capture(page, 'templates-modal');
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
      case 'coordinates':
      case 'cluster':
        if (!param) throw new Error('coordinates requires a lobby code');
        await coordinatesFlow(page, param);
        break;
      case 'templates':
        await templatesFlow(page);
        break;
      default:
        console.error(`Unknown flow: ${command}`);
        console.error('Available: create, join <code>, trivia <code>, coordinates <code>, templates');
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
