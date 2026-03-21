import { expect, test, type Page } from '@playwright/test';

import {
  clickAndWaitForIdle,
  createLobby,
  createPlayerSession,
  joinLobby,
  requirePage,
  safeCloseSessions,
  type PlayerSession,
} from '../support/multiplayer.js';

type TriviaQuestionSpec = {
  author: string;
  text: string;
  correct: string;
  wrongs: [string, string, string];
};

const duplicateQuestion: TriviaQuestionSpec = {
  author: 'Host',
  text: 'Which slot should own the duplicated answer vote?',
  correct: 'Neon',
  wrongs: ['Neon', 'Solar', 'Lunar'],
};

const fillerQuestion: TriviaQuestionSpec = {
  author: 'Player2',
  text: 'Which filler option is correct?',
  correct: 'Mercury',
  wrongs: ['Venus', 'Earth', 'Mars'],
};

async function submitTriviaQuestion(session: PlayerSession, question: TriviaQuestionSpec) {
  const page = requirePage(session);
  await expect(page.locator('#question_text')).toBeVisible({ timeout: 10_000 });
  await page.fill('#question_text', question.text);
  await page.fill('#correct_answer', question.correct);
  await page.fill('#wrong_answer_1', question.wrongs[0]);
  await page.fill('#wrong_answer_2', question.wrongs[1]);
  await page.fill('#wrong_answer_3', question.wrongs[2]);
  await clickAndWaitForIdle(page, 'button:has-text("SUBMIT QUESTION")');
  await expect(page.getByText('QUESTION SUBMITTED')).toBeVisible({ timeout: 10_000 });
}

async function visibleQuestion(page: Page): Promise<string> {
  await expect
    .poll(async () => {
      if (await page.getByText(duplicateQuestion.text, { exact: true }).isVisible()) {
        return duplicateQuestion.text;
      }
      if (await page.getByText(fillerQuestion.text, { exact: true }).isVisible()) {
        return fillerQuestion.text;
      }
      return '';
    }, { timeout: 10_000 })
    .not.toBe('');

  if (await page.getByText(duplicateQuestion.text, { exact: true }).isVisible()) {
    return duplicateQuestion.text;
  }

  return fillerQuestion.text;
}

async function answerByKey(session: PlayerSession, answerKey: string) {
  const page = requirePage(session);
  const button = page.locator(`button[name="answer"][value="${answerKey}"]`).first();
  await expect(button).toBeVisible({ timeout: 10_000 });
  await Promise.all([page.waitForLoadState('networkidle'), button.click()]);
}

async function waitForReveal(page: Page, questionText: string) {
  await expect(page.getByText(questionText, { exact: true })).toBeVisible({ timeout: 10_000 });
  await expect(page.getByText('Correct', { exact: true })).toBeVisible({ timeout: 10_000 });
}

test.describe('Trivia Duplicate Options', () => {
  test('attributes votes to the chosen duplicate slot instead of every matching label', async ({ browser }) => {
    const sessions = await Promise.all([
      createPlayerSession(browser, 'Host'),
      createPlayerSession(browser, 'Player2'),
    ]);

    const [host, player2] = sessions;
    const hostPage = requirePage(host);

    try {
      const code = await createLobby(host, '/trivia', `Trivia Duplicate ${Date.now()}`);
      await joinLobby(player2, code);

      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/start"] button:has-text("START GAME")');

      await submitTriviaQuestion(host, duplicateQuestion);
      await submitTriviaQuestion(player2, fillerQuestion);

      await expect(hostPage.getByText(/2\s*\/\s*2 submitted/i)).toBeVisible({ timeout: 10_000 });
      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/advance"] button:has-text("START ROUND")');

      const firstQuestion = await visibleQuestion(hostPage);
      if (firstQuestion === fillerQuestion.text) {
        await answerByKey(host, 'wrong_answer_1');
        await waitForReveal(hostPage, fillerQuestion.text);
        await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/next-question"] button:has-text("NEXT QUESTION")');
      }

      await expect(hostPage.getByText(duplicateQuestion.text, { exact: true })).toBeVisible({ timeout: 10_000 });
      await answerByKey(player2, 'wrong_answer_1');
      await waitForReveal(hostPage, duplicateQuestion.text);

      const selectedDuplicateRow = hostPage.locator('[data-answer-key="wrong_answer_1"]').first();
      const correctDuplicateRow = hostPage.locator('[data-answer-key="correct_answer"]').first();

      await expect(selectedDuplicateRow).toContainText('Neon');
      await expect(selectedDuplicateRow).toContainText('100%');
      await expect(correctDuplicateRow).toContainText('Neon');
      await expect(correctDuplicateRow).toContainText('0%');
    } finally {
      await safeCloseSessions(sessions);
    }
  });
});
