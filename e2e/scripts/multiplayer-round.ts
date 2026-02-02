/**
 * Multiplayer Round Flow
 *
 * Simulates a complete trivia round with 2 players to validate:
 * - "Host will start the round when ready" messaging
 * - "Waiting for host to continue..." after results
 * - "Overall Leader" on scoreboard
 *
 * Usage: npm run multiplayer
 */

import { chromium, Page, BrowserContext } from '@playwright/test';
import { mkdir } from 'fs/promises';
import { join } from 'path';

const BASE_URL = process.env.BASE_URL || 'http://localhost:3000';
const SCREENSHOTS_DIR = join(import.meta.dirname, '..', 'screenshots', 'multiplayer');

let screenshotIndex = 0;

async function capture(page: Page, role: string, label: string) {
  const filename = `${String(screenshotIndex++).padStart(2, '0')}-${role}-${label}.png`;
  const filepath = join(SCREENSHOTS_DIR, filename);
  await page.screenshot({ path: filepath, fullPage: true });
  console.log(`  [${role}] ${label}: ${filename}`);
  return filepath;
}

async function waitForSelector(page: Page, selector: string, timeout = 5000) {
  try {
    await page.waitForSelector(selector, { timeout });
    return true;
  } catch {
    return false;
  }
}

async function main() {
  await mkdir(SCREENSHOTS_DIR, { recursive: true });

  console.log('\n=== Multiplayer Round Flow ===\n');
  console.log(`Base URL: ${BASE_URL}\n`);

  const browser = await chromium.launch({ headless: true });

  // Create two separate browser contexts (separate sessions/cookies)
  const hostContext = await browser.newContext({ viewport: { width: 1280, height: 800 } });
  const playerContext = await browser.newContext({ viewport: { width: 1280, height: 800 } });

  const hostPage = await hostContext.newPage();
  const playerPage = await playerContext.newPage();

  try {
    // ========== PHASE 1: Create Lobby ==========
    console.log('--- Phase 1: Create Lobby ---');

    await hostPage.goto(`${BASE_URL}/trivia`, { waitUntil: 'networkidle' });
    await capture(hostPage, 'host', 'trivia-home');

    // Host creates lobby
    await hostPage.fill('input[name="name"]', 'Test Game');
    await hostPage.fill('input[name="nickname"]', 'HostPlayer');
    await hostPage.click('button:has-text("INITIALIZE")');
    await hostPage.waitForURL(/\/lobbies\/[A-Z0-9]+/, { timeout: 5000 });
    await hostPage.waitForLoadState('networkidle');

    const lobbyUrl = hostPage.url();
    const lobbyCode = lobbyUrl.split('/lobbies/')[1]?.split('/')[0];
    console.log(`  Lobby created: ${lobbyCode}`);

    await capture(hostPage, 'host', 'lobby-created');

    // ========== PHASE 2: Player Joins ==========
    console.log('\n--- Phase 2: Player Joins ---');

    await playerPage.goto(`${BASE_URL}/lobbies/${lobbyCode}`, { waitUntil: 'networkidle' });
    await capture(playerPage, 'player', 'join-page');

    await playerPage.fill('input[name="nickname"]', 'Player2');
    await playerPage.click('button:has-text("CONNECT")');
    await playerPage.waitForLoadState('networkidle');
    await capture(playerPage, 'player', 'joined-lobby');

    // Wait for host page to update via WebSocket
    await hostPage.waitForTimeout(1000);
    await capture(hostPage, 'host', 'player-joined');

    // ========== PHASE 3: Host Starts Game ==========
    console.log('\n--- Phase 3: Host Starts Game ---');

    // Host clicks start
    const startButton = hostPage.locator('button:has-text("START GAME")');
    if (await startButton.isVisible()) {
      await startButton.click();
      await hostPage.waitForLoadState('networkidle');
    }

    await hostPage.waitForTimeout(500);
    await playerPage.waitForTimeout(500);

    await capture(hostPage, 'host', 'question-submit-form');
    await capture(playerPage, 'player', 'question-submit-form');

    // ========== PHASE 3.5: Template Modal ==========
    console.log('\n--- Phase 3.5: Template Modal ---');

    const templateCta = hostPage.locator('button:has-text("Need help? Use a template")');
    if (await templateCta.isVisible().catch(() => false)) {
      await templateCta.click();
      await hostPage.waitForSelector('#templates-modal', { timeout: 5000 });
      await capture(hostPage, 'host', 'templates-modal-open');

      const firstTemplate = hostPage.locator('#templates-content button').first();
      await firstTemplate.waitFor({ state: 'visible', timeout: 5000 });
      await firstTemplate.click();
      await hostPage.waitForTimeout(300);
      await capture(hostPage, 'host', 'template-selected');

      const templateIdValue = await hostPage.locator('#template_id').inputValue();
      const questionValue = await hostPage.locator('#question_text').inputValue();
      console.log(`  ✓ Template selected: ${templateIdValue !== ''} | question filled: ${questionValue.trim() !== ''}`);
    } else {
      console.log('  ⚠ Template CTA not visible');
    }

    // ========== PHASE 4: Submit Questions ==========
    console.log('\n--- Phase 4: Submit Questions ---');

    // Host submits question
    const hostQuestionValue = (await hostPage.locator('textarea[name="question_text"]').inputValue()).trim();
    const hostCorrectValue = (await hostPage.locator('input[name="correct_answer"]').inputValue()).trim();
    if (!hostQuestionValue || !hostCorrectValue) {
      await hostPage.fill('textarea[name="question_text"]', 'What is the capital of France?');
      await hostPage.fill('input[name="correct_answer"]', 'Paris');
      await hostPage.fill('input[name="wrong_answer_1"]', 'London');
      await hostPage.fill('input[name="wrong_answer_2"]', 'Berlin');
      await hostPage.fill('input[name="wrong_answer_3"]', 'Madrid');
    }
    await hostPage.click('button:has-text("SUBMIT QUESTION")');
    await hostPage.waitForLoadState('networkidle');

    await capture(hostPage, 'host', 'question-submitted');

    // Player submits question
    await playerPage.fill('textarea[name="question_text"]', 'What is 2 + 2?');
    await playerPage.fill('input[name="correct_answer"]', '4');
    await playerPage.fill('input[name="wrong_answer_1"]', '3');
    await playerPage.fill('input[name="wrong_answer_2"]', '5');
    await playerPage.fill('input[name="wrong_answer_3"]', '22');
    await playerPage.click('button:has-text("SUBMIT QUESTION")');
    await playerPage.waitForLoadState('networkidle');

    // KEY VALIDATION: Player should see "Host will start the round when ready"
    await capture(playerPage, 'player', 'question-submitted-waiting-for-host');

    // Check for the text
    const waitingForHostText = await playerPage.locator('text=Host will start the round when ready').isVisible();
    console.log(`  ✓ "Host will start the round when ready" visible: ${waitingForHostText}`);

    // Wait for WebSocket updates
    await hostPage.waitForTimeout(1000);
    await capture(hostPage, 'host', 'all-submitted-can-start');

    // ========== PHASE 5: Host Starts Round ==========
    console.log('\n--- Phase 5: Host Starts Round ---');

    const startRoundButton = hostPage.locator('button:has-text("START ROUND")');
    if (await startRoundButton.isVisible()) {
      await startRoundButton.click();
      await hostPage.waitForLoadState('networkidle');
    }

    await hostPage.waitForTimeout(1000);
    await playerPage.waitForTimeout(1000);

    await capture(hostPage, 'host', 'answering-question');
    await capture(playerPage, 'player', 'answering-question');

    // ========== PHASE 6: Answer Questions ==========
    console.log('\n--- Phase 6: Answer Questions ---');

    // Find and click an answer (not the author's own question)
    // The answer buttons should be visible
    for (let questionNum = 1; questionNum <= 2; questionNum++) {
      console.log(`  Question ${questionNum}:`);

      // Take before-answer screenshots
      await capture(hostPage, 'host', `q${questionNum}-before-answer`);
      await capture(playerPage, 'player', `q${questionNum}-before-answer`);

      // Host answers (if not author - check if form is present)
      const hostAnswerBtn = hostPage.locator('button[type="submit"]:has-text("A."), button[type="submit"]:has-text("B."), button[type="submit"]:has-text("C."), button[type="submit"]:has-text("D.")').first();
      if (await hostAnswerBtn.isVisible({ timeout: 1000 }).catch(() => false)) {
        await hostAnswerBtn.click();
        await hostPage.waitForTimeout(500);
        console.log('    Host answered');
      } else {
        console.log('    Host is author (skipped)');
      }

      // Player answers (if not author)
      const playerAnswerBtn = playerPage.locator('button[type="submit"]:has-text("A."), button[type="submit"]:has-text("B."), button[type="submit"]:has-text("C."), button[type="submit"]:has-text("D.")').first();
      if (await playerAnswerBtn.isVisible({ timeout: 1000 }).catch(() => false)) {
        await playerAnswerBtn.click();
        await playerPage.waitForTimeout(500);
        console.log('    Player answered');
      } else {
        console.log('    Player is author (skipped)');
      }

      // Wait for results to appear
      await hostPage.waitForTimeout(1500);
      await playerPage.waitForTimeout(500);

      // Capture results screen
      await capture(hostPage, 'host', `q${questionNum}-results`);
      await capture(playerPage, 'player', `q${questionNum}-results`);

      // KEY VALIDATION: Player should see "Waiting for host to continue..."
      const waitingForHostContinue = await playerPage.locator('text=Waiting for host to continue').isVisible();
      console.log(`    ✓ "Waiting for host to continue..." visible: ${waitingForHostContinue}`);

      // Host advances to next question (if not last)
      const nextButton = hostPage.locator('button:has-text("NEXT QUESTION")');
      if (await nextButton.isVisible({ timeout: 2000 }).catch(() => false)) {
        await nextButton.click();
        await hostPage.waitForTimeout(1000);
        await playerPage.waitForTimeout(500);
      }
    }

    // ========== PHASE 7: Scoreboard ==========
    console.log('\n--- Phase 7: Scoreboard ---');

    await hostPage.waitForTimeout(1000);
    await playerPage.waitForTimeout(500);

    await capture(hostPage, 'host', 'scoreboard');
    await capture(playerPage, 'player', 'scoreboard');

    // KEY VALIDATION: Check for "Overall Leader" text
    const overallLeaderHost = await hostPage.locator('text=Overall Leader').isVisible();
    const overallLeaderPlayer = await playerPage.locator('text=Overall Leader').isVisible();
    console.log(`  ✓ "Overall Leader" visible (host): ${overallLeaderHost}`);
    console.log(`  ✓ "Overall Leader" visible (player): ${overallLeaderPlayer}`);

    // Check for trophy emoji
    const trophyCount = await hostPage.locator('text=🏆').count();
    console.log(`  ✓ Trophy emoji visible: ${trophyCount > 0} (count: ${trophyCount})`);

    console.log('\n=== Flow Complete ===');
    console.log(`Screenshots saved to: ${SCREENSHOTS_DIR}`);
    console.log(`Total screenshots: ${screenshotIndex}`);

  } catch (error) {
    console.error('\nError during flow:', error);
    await capture(hostPage, 'host', 'error-state');
    await capture(playerPage, 'player', 'error-state');
    throw error;
  } finally {
    await hostContext.close();
    await playerContext.close();
    await browser.close();
  }
}

main();
