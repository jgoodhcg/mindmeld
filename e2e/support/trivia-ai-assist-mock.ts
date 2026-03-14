import type { Page } from '@playwright/test';

export const triviaAiAssistMockTopic = 'city landmarks';

export const triviaAiAssistMockResponse = {
  question_text: 'Which city is known as the City of Light?',
  correct_answer: 'Paris',
  wrong_answer_1: 'Rome',
  wrong_answer_2: 'Madrid',
  wrong_answer_3: 'Vienna',
  source: 'openrouter',
};

export type TriviaAiAssistMockState = {
  lastTopic: string;
  callCount: number;
};

export async function installTriviaAiAssistMock(page: Page): Promise<TriviaAiAssistMockState> {
  const state: TriviaAiAssistMockState = {
    lastTopic: '',
    callCount: 0,
  };

  await page.route('**/lobbies/*/trivia/generate-question', async (route) => {
    state.callCount += 1;
    const body = route.request().postData() ?? '';
    const params = new URLSearchParams(body);
    state.lastTopic = params.get('topic') ?? '';

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(triviaAiAssistMockResponse),
    });
  });

  return state;
}
