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
  test('uses a deterministic mocked AI response to populate the submit form', async ({ page }) => {
    const aiAssistMock = await installTriviaAiAssistMock(page);

    await createTriviaLobby(page);
    await startTriviaGame(page);

    const assistSummary = page.locator('#assist-panel summary');
    await expect(assistSummary).toBeVisible({ timeout: 10000 });
    await assistSummary.click();

    await page.fill('#assist_topic', triviaAiAssistMockTopic);

    await Promise.all([
      page.waitForResponse((response) => response.url().includes('/trivia/generate-question') && response.request().method() === 'POST'),
      page.click('#assist_generate_button'),
    ]);

    await expect(page.locator('#question_text')).toHaveValue(triviaAiAssistMockResponse.question_text);
    await expect(page.locator('#correct_answer')).toHaveValue(triviaAiAssistMockResponse.correct_answer);
    await expect(page.locator('#wrong_answer_1')).toHaveValue(triviaAiAssistMockResponse.wrong_answer_1);
    await expect(page.locator('#wrong_answer_2')).toHaveValue(triviaAiAssistMockResponse.wrong_answer_2);
    await expect(page.locator('#wrong_answer_3')).toHaveValue(triviaAiAssistMockResponse.wrong_answer_3);
    await expect(page.locator('#assist_status')).toContainText('Draft generated from AI model.');

    expect(aiAssistMock.callCount).toBe(1);
    expect(aiAssistMock.lastTopic).toBe(triviaAiAssistMockTopic);
  });
});
