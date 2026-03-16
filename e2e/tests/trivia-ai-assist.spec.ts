import { expect, test, type Page } from '@playwright/test';

import {
  installTriviaAiAssistMock,
  triviaAiAssistMockResponse,
  triviaAiAssistMockTopic,
} from '../support/trivia-ai-assist-mock.js';

async function createTriviaLobby(page: Page) {
  await page.goto('/trivia', { waitUntil: 'networkidle' });
  await expect(page.locator('form[action="/lobbies"]')).toBeVisible({ timeout: 10000 });

  await page.fill('form[action="/lobbies"] input[name="name"]', `Trivia AI E2E ${Date.now()}`);
  await page.fill('form[action="/lobbies"] input[name="nickname"]', 'Host');
  await Promise.all([
    page.waitForURL(/\/lobbies\/[A-Z0-9]{6}$/, { timeout: 10000 }),
    page.click('form[action="/lobbies"] button:has-text("INITIALIZE")'),
  ]);
}

async function startTriviaGame(page: Page) {
  const startButton = page.locator('form[action$="/trivia/start"] button:has-text("START GAME")');
  await expect(startButton).toBeVisible({ timeout: 10000 });
  await Promise.all([
    page.waitForLoadState('networkidle'),
    startButton.click(),
  ]);

  await expect(page.locator('form[action$="/questions"]')).toBeVisible({ timeout: 10000 });
}

test.describe('Trivia AI Assist', () => {
  test('shows loading feedback, keeps keyboard focus usable, and populates the submit form', async ({ page }) => {
    const aiAssistMock = await installTriviaAiAssistMock(page, { delayMs: 400 });

    await createTriviaLobby(page);
    await startTriviaGame(page);

    const assistSummary = page.locator('#assist-panel summary');
    await expect(assistSummary).toBeVisible({ timeout: 10000 });
    await assistSummary.click();

    await page.fill('#assist_topic', triviaAiAssistMockTopic);
    await page.locator('#assist_topic').focus();
    await page.keyboard.press('Tab');
    await expect(page.locator('#assist_generate_button')).toBeFocused();

    const responsePromise = page.waitForResponse((response) => response.url().includes('/trivia/generate-question') && response.request().method() === 'POST');
    await page.keyboard.press('Enter');

    await expect(page.locator('#assist_generate_button')).toBeDisabled();
    await expect(page.locator('#assist_generate_button_busy')).toBeVisible();
    await expect(page.locator('#assist_status_spinner')).toBeVisible();
    await expect(page.locator('#assist_status')).toContainText('Generating draft...');

    await responsePromise;

    await expect(page.locator('#question_text')).toHaveValue(triviaAiAssistMockResponse.question_text);
    await expect(page.locator('#correct_answer')).toHaveValue(triviaAiAssistMockResponse.correct_answer);
    await expect(page.locator('#wrong_answer_1')).toHaveValue(triviaAiAssistMockResponse.wrong_answer_1);
    await expect(page.locator('#wrong_answer_2')).toHaveValue(triviaAiAssistMockResponse.wrong_answer_2);
    await expect(page.locator('#wrong_answer_3')).toHaveValue(triviaAiAssistMockResponse.wrong_answer_3);
    await expect(page.locator('#assist_generate_button')).toBeEnabled();
    await expect(page.locator('#assist_generate_button_busy')).toBeHidden();
    await expect(page.locator('#assist_status')).toContainText('Draft generated from AI model.');

    expect(aiAssistMock.callCount).toBe(1);
    expect(aiAssistMock.lastTopic).toBe(triviaAiAssistMockTopic);
  });
});
