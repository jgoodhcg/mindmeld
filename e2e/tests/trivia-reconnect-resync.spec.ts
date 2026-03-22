import { expect, test } from '@playwright/test';

import {
  clickAndWaitForIdle,
  createLobby,
  createPlayerSession,
  disconnectPlayer,
  joinLobby,
  reconnectPlayer,
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

const hostQuestion: TriviaQuestionSpec = {
  author: 'Host',
  text: 'Host reconnect resync question?',
  correct: 'Alpha',
  wrongs: ['Beta', 'Gamma', 'Delta'],
};

const player2Question: TriviaQuestionSpec = {
  author: 'Player2',
  text: 'Player2 reconnect resync question?',
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

test.describe('Trivia Reconnect Resync', () => {
  test('returns a reconnecting player to the active question instead of rewinding the room state', async ({ browser }) => {
    const sessions = await Promise.all([
      createPlayerSession(browser, 'Host'),
      createPlayerSession(browser, 'Player2'),
    ]);

    const [host, player2] = sessions;
    const hostPage = requirePage(host);

    try {
      const code = await createLobby(host, '/trivia', `Trivia Resync ${Date.now()}`);
      await joinLobby(player2, code);

      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/start"] button:has-text("START GAME")');

      await submitTriviaQuestion(host, hostQuestion);
      await submitTriviaQuestion(player2, player2Question);

      await expect(hostPage.getByText(/2\s*\/\s*2 submitted/i)).toBeVisible({ timeout: 10_000 });
      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/advance"] button:has-text("START ROUND")');

      const currentQuestion = (await hostPage.getByText(hostQuestion.text, { exact: true }).isVisible())
        ? hostQuestion
        : player2Question;

      const reconnectingSession = currentQuestion.author === host.name ? player2 : host;
      await disconnectPlayer(reconnectingSession);

      const reconnectedPage = await reconnectPlayer(reconnectingSession, code);
      await expect(reconnectedPage.getByText(currentQuestion.text, { exact: true })).toBeVisible({
        timeout: 10_000,
      });
      await expect(reconnectedPage.getByText('QUESTION SUBMITTED')).not.toBeVisible();
      await expect(reconnectedPage.locator('button[name="answer"]').first()).toBeVisible({ timeout: 10_000 });

      await Promise.all([
        reconnectedPage.waitForLoadState('networkidle'),
        reconnectedPage.locator('button[name="answer"]').first().click(),
      ]);

      await expect(hostPage.getByText(currentQuestion.text, { exact: true })).toBeVisible({ timeout: 10_000 });
      await expect(hostPage.getByText('Correct', { exact: true })).toBeVisible({ timeout: 10_000 });
    } finally {
      await safeCloseSessions(sessions);
    }
  });
});
