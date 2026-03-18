import { expect, test, type Page } from '@playwright/test';

import {
  clickAndWaitForIdle,
  connectedPages,
  createLobby,
  createPlayerSession,
  disconnectGracePeriodMs,
  disconnectPlayer,
  hostBadge,
  expectHostPlayer,
  expectLobbyPlayerCount,
  expectPlayerState,
  joinLobby,
  reconnectPlayer,
  requirePage,
  safeCloseSessions,
  transferHost,
  type PlayerSession,
} from '../support/multiplayer.js';

type TriviaQuestionSpec = {
  author: string;
  text: string;
  correct: string;
  wrongs: [string, string, string];
};

const triviaQuestions: TriviaQuestionSpec[] = [
  {
    author: 'Host',
    text: 'Host calibration question?',
    correct: 'Host correct',
    wrongs: ['Host wrong A', 'Host wrong B', 'Host wrong C'],
  },
  {
    author: 'Player2',
    text: 'Player2 orbit question?',
    correct: 'Player2 correct',
    wrongs: ['Player2 wrong A', 'Player2 wrong B', 'Player2 wrong C'],
  },
  {
    author: 'Player3',
    text: 'Player3 relay question?',
    correct: 'Player3 correct',
    wrongs: ['Player3 wrong A', 'Player3 wrong B', 'Player3 wrong C'],
  },
  {
    author: 'Player4',
    text: 'Player4 signal question?',
    correct: 'Player4 correct',
    wrongs: ['Player4 wrong A', 'Player4 wrong B', 'Player4 wrong C'],
  },
  {
    author: 'Player5',
    text: 'Player5 vector question?',
    correct: 'Player5 correct',
    wrongs: ['Player5 wrong A', 'Player5 wrong B', 'Player5 wrong C'],
  },
];

const questionByAuthor = new Map(triviaQuestions.map((question) => [question.author, question]));

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

async function findVisibleQuestion(page: Page): Promise<TriviaQuestionSpec> {
  for (const question of triviaQuestions) {
    if (await page.getByText(question.text, { exact: true }).isVisible()) {
      return question;
    }
  }

  throw new Error(`unable to identify visible trivia question on ${page.url()}`);
}

async function answerTriviaQuestion(session: PlayerSession, answerText: string) {
  const page = requirePage(session);
  const answerButton = page.locator('button[name="answer"]').filter({ hasText: answerText }).first();
  await expect(answerButton).toBeVisible({ timeout: 10_000 });
  await Promise.all([page.waitForLoadState('networkidle'), answerButton.click()]);
}

async function waitForRevealedResults(pages: Page[], question: TriviaQuestionSpec) {
  for (const page of pages) {
    await expect(page.getByText(question.text, { exact: true })).toBeVisible({
      timeout: disconnectGracePeriodMs() + 10_000,
    });
    await expect(page.getByText('Correct', { exact: true })).toBeVisible({
      timeout: disconnectGracePeriodMs() + 10_000,
    });
  }
}

async function waitForQuestionVisible(pages: Page[]) {
  for (const page of pages) {
    await expect
      .poll(async () => {
        for (const question of triviaQuestions) {
          if (await page.getByText(question.text, { exact: true }).isVisible()) {
            return question.text;
          }
        }
        return '';
      }, { timeout: 10_000 })
      .not.toBe('');
  }
}

function pickSession(
  sessions: PlayerSession[],
  preferredNames: string[],
  excludedNames: string[],
): PlayerSession {
  const excluded = new Set(excludedNames);
  let match: PlayerSession | undefined;
  for (const name of preferredNames) {
    const session = sessions.find((candidate) => candidate.name === name);
    if (session && !excluded.has(session.name)) {
      match = session;
      break;
    }
  }

  if (!match) {
    throw new Error(`unable to select session from ${preferredNames.join(', ')}`);
  }

  return match;
}

function answerersForQuestion(
  sessions: PlayerSession[],
  question: TriviaQuestionSpec,
  excludedNames: string[] = [],
): PlayerSession[] {
  const excluded = new Set([question.author, ...excludedNames]);
  return sessions.filter((session) => !excluded.has(session.name));
}

test.describe('Trivia Resilience', () => {
  test('handles a 6-player session with late join, reconnect, and grace-expiry continuation', async ({
    browser,
  }) => {
    test.setTimeout(Math.max(240_000, disconnectGracePeriodMs() * 7));

    const sessions = await Promise.all(
      ['Host', 'Player2', 'Player3', 'Player4', 'Player5', 'Player6'].map((name) =>
        createPlayerSession(browser, name),
      ),
    );

    const [host, player2, player3, player4, player5, player6] = sessions;
    const initialPlayers = [host, player2, player3, player4, player5];
    const hostPage = requirePage(host);

    try {
      const code = await createLobby(host, '/trivia', `Trivia Resilience ${Date.now()}`);

      for (const session of [player2, player3, player4, player5]) {
        await joinLobby(session, code);
      }

      await expectLobbyPlayerCount(hostPage, 5);
      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/start"] button:has-text("START GAME")');

      for (const session of initialPlayers) {
        const question = questionByAuthor.get(session.name);
        if (!question) {
          throw new Error(`missing trivia question for ${session.name}`);
        }
        await submitTriviaQuestion(session, question);
      }

      await expect(hostPage.getByText(/5\s*\/\s*5 submitted/i)).toBeVisible({ timeout: 10_000 });
      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/advance"] button:has-text("START ROUND")');
      await waitForQuestionVisible(connectedPages(initialPlayers));

      const roundOneQuestion = await findVisibleQuestion(hostPage);
      const sameRoundReconnect = pickSession(
        initialPlayers,
        ['Player5', 'Player4', 'Player3', 'Player2'],
        [roundOneQuestion.author],
      );

      const lateJoinPage = await joinLobby(player6, code);
      await expectLobbyPlayerCount(hostPage, 6);
      await expect(lateJoinPage.getByText(roundOneQuestion.text, { exact: true })).toBeVisible({
        timeout: 10_000,
      });

      await disconnectPlayer(sameRoundReconnect);
      await expectPlayerState(hostPage, sameRoundReconnect.name, 'reconnecting');

      for (const session of answerersForQuestion(sessions, roundOneQuestion, [sameRoundReconnect.name])) {
        await answerTriviaQuestion(session, roundOneQuestion.wrongs[0]);
      }

      const reconnectedSameRound = await reconnectPlayer(sameRoundReconnect, code);
      await expect(reconnectedSameRound.getByText(roundOneQuestion.text, { exact: true })).toBeVisible({
        timeout: 10_000,
      });
      await expectPlayerState(hostPage, sameRoundReconnect.name, 'connected');
      await answerTriviaQuestion(sameRoundReconnect, roundOneQuestion.wrongs[0]);

      await waitForRevealedResults(connectedPages(sessions), roundOneQuestion);

      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/next-question"] button:has-text("NEXT QUESTION")');
      await waitForQuestionVisible(connectedPages(sessions));

      const roundTwoQuestion = await findVisibleQuestion(hostPage);
      const persistentDisconnect = pickSession(
        initialPlayers,
        ['Player4', 'Player3', 'Player2', 'Player5'],
        [roundTwoQuestion.author, sameRoundReconnect.name],
      );

      await disconnectPlayer(persistentDisconnect);
      await expectPlayerState(hostPage, persistentDisconnect.name, 'reconnecting');

      for (const session of answerersForQuestion(sessions, roundTwoQuestion, [persistentDisconnect.name])) {
        await answerTriviaQuestion(session, roundTwoQuestion.wrongs[0]);
      }

      await expectPlayerState(hostPage, persistentDisconnect.name, 'disconnected');
      await waitForRevealedResults(
        connectedPages(sessions.filter((session) => session.name !== persistentDisconnect.name)),
        roundTwoQuestion,
      );

      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/next-question"] button:has-text("NEXT QUESTION")');
      await waitForQuestionVisible(
        connectedPages(sessions.filter((session) => session.name !== persistentDisconnect.name)),
      );

      const reconnectedPersistent = await reconnectPlayer(persistentDisconnect, code);
      const roundThreeQuestion = await findVisibleQuestion(hostPage);
      await expect(reconnectedPersistent.getByText(roundThreeQuestion.text, { exact: true })).toBeVisible({
        timeout: 10_000,
      });
      await expectPlayerState(hostPage, persistentDisconnect.name, 'connected');

      for (const session of answerersForQuestion(sessions, roundThreeQuestion)) {
        await answerTriviaQuestion(session, roundThreeQuestion.wrongs[0]);
      }

      await waitForRevealedResults(connectedPages(sessions), roundThreeQuestion);
    } finally {
      await safeCloseSessions(sessions);
    }
  });

  test('transfers host after grace expiry and lets the fallback host advance trivia', async ({
    browser,
  }) => {
    test.setTimeout(Math.max(150_000, disconnectGracePeriodMs() * 3));

    const sessions = await Promise.all(
      ['Host', 'Player2', 'Player3'].map((name) => createPlayerSession(browser, name)),
    );

    const [host, player2, player3] = sessions;
    const hostPage = requirePage(host);

    try {
      const code = await createLobby(host, '/trivia', `Trivia Host Fallback ${Date.now()}`);

      await joinLobby(player2, code);
      await joinLobby(player3, code);

      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/start"] button:has-text("START GAME")');

      for (const session of sessions) {
        const question = questionByAuthor.get(session.name);
        if (!question) {
          throw new Error(`missing trivia question for ${session.name}`);
        }
        await submitTriviaQuestion(session, question);
      }

      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/advance"] button:has-text("START ROUND")');

      const player2Page = requirePage(player2);
      const player3Page = requirePage(player3);
      await waitForQuestionVisible([hostPage, player2Page, player3Page]);

      const currentQuestion = await findVisibleQuestion(hostPage);
      for (const session of answerersForQuestion(sessions, currentQuestion)) {
        await answerTriviaQuestion(session, currentQuestion.wrongs[0]);
      }
      await waitForRevealedResults([hostPage, player2Page, player3Page], currentQuestion);

      await disconnectPlayer(host);
      await expectPlayerState(player2Page, host.name, 'reconnecting');
      await expectPlayerState(player2Page, host.name, 'disconnected');
      await expectHostPlayer(player2Page, player2.name);
      await expectHostPlayer(player3Page, player2.name);

      await expect(
        player2Page.locator('form[action$="/trivia/next-question"] button:has-text("NEXT QUESTION")'),
      ).toBeVisible({ timeout: 10_000 });
      await clickAndWaitForIdle(
        player2Page,
        'form[action$="/trivia/next-question"] button:has-text("NEXT QUESTION")',
      );

      await waitForQuestionVisible([player2Page, player3Page]);

      const reconnectedHost = await reconnectPlayer(host, code);
      await expectHostPlayer(reconnectedHost, player2.name);
      await expect(hostBadge(reconnectedHost, host.name)).not.toBeVisible();
    } finally {
      await safeCloseSessions(sessions);
    }
  });

  test('allows the host to hand off control between revealed trivia questions', async ({
    browser,
  }) => {
    test.setTimeout(120_000);

    const sessions = await Promise.all(
      ['Host', 'Player2', 'Player3'].map((name) => createPlayerSession(browser, name)),
    );

    const [host, player2, player3] = sessions;
    const hostPage = requirePage(host);

    try {
      const code = await createLobby(host, '/trivia', `Trivia Manual Transfer ${Date.now()}`);

      await joinLobby(player2, code);
      await joinLobby(player3, code);

      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/start"] button:has-text("START GAME")');

      for (const session of sessions) {
        const question = questionByAuthor.get(session.name);
        if (!question) {
          throw new Error(`missing trivia question for ${session.name}`);
        }
        await submitTriviaQuestion(session, question);
      }

      await clickAndWaitForIdle(hostPage, 'form[action$="/trivia/advance"] button:has-text("START ROUND")');

      const player2Page = requirePage(player2);
      const player3Page = requirePage(player3);
      await waitForQuestionVisible([hostPage, player2Page, player3Page]);

      const currentQuestion = await findVisibleQuestion(hostPage);
      for (const session of answerersForQuestion(sessions, currentQuestion)) {
        await answerTriviaQuestion(session, currentQuestion.wrongs[0]);
      }
      await waitForRevealedResults([hostPage, player2Page, player3Page], currentQuestion);

      await expect(
        hostPage.locator('form[action$="/host-transfer"] select[name="target_player_id"]'),
      ).toBeVisible({ timeout: 10_000 });
      await transferHost(hostPage, player2.name);

      await expectHostPlayer(hostPage, player2.name);
      await expectHostPlayer(player2Page, player2.name);
      await expectHostPlayer(player3Page, player2.name);

      await expect(
        hostPage.locator('form[action$="/trivia/next-question"] button:has-text("NEXT QUESTION")'),
      ).not.toBeVisible();
      await expect(
        player2Page.locator('form[action$="/trivia/next-question"] button:has-text("NEXT QUESTION")'),
      ).toBeVisible({ timeout: 10_000 });

      await clickAndWaitForIdle(
        player2Page,
        'form[action$="/trivia/next-question"] button:has-text("NEXT QUESTION")',
      );
      await waitForQuestionVisible([hostPage, player2Page, player3Page]);
    } finally {
      await safeCloseSessions(sessions);
    }
  });
});
