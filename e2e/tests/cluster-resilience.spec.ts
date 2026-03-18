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

const coordinatePlan: Record<string, Array<[string, string]>> = {
  Host: [
    ['0.12', '0.18'],
    ['0.18', '0.30'],
    ['0.24', '0.42'],
  ],
  Player2: [
    ['0.26', '0.34'],
    ['0.32', '0.48'],
    ['0.38', '0.62'],
  ],
  Player3: [
    ['0.40', '0.58'],
    ['0.44', '0.22'],
    ['0.48', '0.74'],
  ],
  Player4: [
    ['0.54', '0.46'],
    ['0.58', '0.66'],
    ['0.62', '0.28'],
  ],
  Player5: [
    ['0.68', '0.72'],
    ['0.72', '0.40'],
    ['0.76', '0.56'],
  ],
  Player6: [
    ['0.82', '0.26'],
    ['0.86', '0.54'],
    ['0.88', '0.38'],
  ],
};

async function selectCoordinate(page: Page, x: string, y: string) {
  const plane = page.locator('#cluster-plane-input');
  await expect(plane).toBeVisible({ timeout: 10_000 });
  await plane.scrollIntoViewIfNeeded();

  let box = await plane.boundingBox();
  if (!box) {
    await page.waitForTimeout(250);
    box = await plane.boundingBox();
  }
  if (!box) {
    throw new Error('cluster input plane is missing a bounding box');
  }

  const xNum = Number.parseFloat(x);
  const yNum = Number.parseFloat(y);
  const clickX = box.x + box.width * xNum;
  const clickY = box.y + box.height * (1 - yNum);
  await page.mouse.click(clickX, clickY);
}

async function submitCoordinate(session: PlayerSession, roundNumber: number) {
  const page = requirePage(session);
  const [x, y] = coordinatePlan[session.name][roundNumber - 1];
  await selectCoordinate(page, x, y);
  await clickAndWaitForIdle(page, 'button:has-text("SUBMIT POINT")');
  await expect(page.getByText('COORDINATE LOCKED').or(page.getByText('CENTROID REVEALED'))).toBeVisible({
    timeout: disconnectGracePeriodMs() + 10_000,
  });
}

async function waitForRound(pages: Page[], roundNumber: number) {
  for (const page of pages) {
    await expect(page.getByText(`ROUND ${roundNumber}`)).toBeVisible({ timeout: 10_000 });
    await expect(page.locator('#cluster-plane-input')).toBeVisible({ timeout: 10_000 });
  }
}

async function waitForReveal(pages: Page[]) {
  for (const page of pages) {
    await expect(page.getByText('CENTROID REVEALED')).toBeVisible({
      timeout: disconnectGracePeriodMs() + 10_000,
    });
    await expect(page.locator('[title="Centroid target"]')).toBeVisible({ timeout: 10_000 });
  }
}

test.describe('Cluster Resilience', () => {
  test('handles a 6-player session with late join, reconnect, and grace-expiry continuation', async ({
    browser,
  }) => {
    test.setTimeout(Math.max(180_000, disconnectGracePeriodMs() * 6));

    const sessions = await Promise.all(
      ['Host', 'Player2', 'Player3', 'Player4', 'Player5', 'Player6'].map((name) =>
        createPlayerSession(browser, name),
      ),
    );

    const [host, player2, player3, player4, player5, player6] = sessions;
    const initialPlayers = [host, player2, player3, player4, player5];

    try {
      const code = await createLobby(host, '/cluster', `Cluster Resilience ${Date.now()}`);

      for (const session of [player2, player3, player4, player5]) {
        await joinLobby(session, code);
      }

      const hostPage = requirePage(host);
      await expectLobbyPlayerCount(hostPage, 5);
      await expect(hostPage.getByText(/Active now:\s*5/i)).toBeVisible({ timeout: 10_000 });

      await clickAndWaitForIdle(
        hostPage,
        'form[action$="/cluster/start"] button:has-text("START CLUSTER")',
      );
      await waitForRound(connectedPages(initialPlayers), 1);

      await submitCoordinate(host, 1);
      await submitCoordinate(player2, 1);
      await submitCoordinate(player3, 1);
      await expect(hostPage.getByText(/3\s*\/\s*5 submitted/i)).toBeVisible({ timeout: 10_000 });

      const lateJoinPage = await joinLobby(player6, code);
      await expectLobbyPlayerCount(hostPage, 6);
      await expect(lateJoinPage.getByText('ROUND 1')).toBeVisible({ timeout: 10_000 });
      await expect(lateJoinPage.locator('#cluster-plane-input')).toBeVisible({ timeout: 10_000 });
      await expect(hostPage.getByText(/3\s*\/\s*6 submitted/i)).toBeVisible({ timeout: 10_000 });

      await disconnectPlayer(player5);
      await expectPlayerState(hostPage, player5.name, 'reconnecting');

      await submitCoordinate(player4, 1);
      await submitCoordinate(player6, 1);
      await expect(hostPage.getByText(/5\s*\/\s*6 submitted/i)).toBeVisible({ timeout: 10_000 });

      const reconnectedPlayer5 = await reconnectPlayer(player5, code);
      await expect(reconnectedPlayer5.getByText('ROUND 1')).toBeVisible({ timeout: 10_000 });
      await expectPlayerState(hostPage, player5.name, 'connected');
      await submitCoordinate(player5, 1);

      await waitForReveal(connectedPages(sessions));

      await clickAndWaitForIdle(hostPage, 'form[action$="/cluster/next"] button');
      await waitForRound(connectedPages(sessions), 2);

      await disconnectPlayer(player4);
      await expectPlayerState(hostPage, player4.name, 'reconnecting');

      for (const session of [host, player2, player3, player5, player6]) {
        await submitCoordinate(session, 2);
      }

      await expectPlayerState(hostPage, player4.name, 'disconnected');
      await waitForReveal(connectedPages([host, player2, player3, player5, player6]));

      await clickAndWaitForIdle(hostPage, 'form[action$="/cluster/next"] button');
      await waitForRound(connectedPages([host, player2, player3, player5, player6]), 3);

      const reconnectedPlayer4 = await reconnectPlayer(player4, code);
      await expect(reconnectedPlayer4.getByText('ROUND 3')).toBeVisible({ timeout: 10_000 });
      await expect(reconnectedPlayer4.locator('#cluster-plane-input')).toBeVisible({ timeout: 10_000 });
      await expectPlayerState(hostPage, player4.name, 'connected');

      await submitCoordinate(player4, 3);
      await expect(hostPage.getByText(/1\s*\/\s*6 submitted/i)).toBeVisible({ timeout: 10_000 });
    } finally {
      await safeCloseSessions(sessions);
    }
  });

  test('transfers host after grace expiry and lets the fallback host continue the session', async ({
    browser,
  }) => {
    test.setTimeout(Math.max(120_000, disconnectGracePeriodMs() * 3));

    const sessions = await Promise.all(
      ['Host', 'Player2', 'Player3'].map((name) => createPlayerSession(browser, name)),
    );

    const [host, player2, player3] = sessions;

    try {
      const code = await createLobby(host, '/cluster', `Cluster Host Fallback ${Date.now()}`);

      await joinLobby(player2, code);
      await joinLobby(player3, code);

      const hostPage = requirePage(host);
      const player2Page = requirePage(player2);
      const player3Page = requirePage(player3);

      await clickAndWaitForIdle(
        hostPage,
        'form[action$="/cluster/start"] button:has-text("START CLUSTER")',
      );
      await waitForRound([hostPage, player2Page, player3Page], 1);

      await submitCoordinate(host, 1);
      await submitCoordinate(player2, 1);
      await submitCoordinate(player3, 1);
      await waitForReveal([hostPage, player2Page, player3Page]);

      await disconnectPlayer(host);
      await expectPlayerState(player2Page, host.name, 'reconnecting');
      await expectPlayerState(player2Page, host.name, 'disconnected');
      await expectHostPlayer(player2Page, player2.name);
      await expectHostPlayer(player3Page, player2.name);

      await expect(player2Page.locator('form[action$="/cluster/next"] button')).toBeVisible({
        timeout: 10_000,
      });
      await clickAndWaitForIdle(player2Page, 'form[action$="/cluster/next"] button');

      await waitForRound([player2Page, player3Page], 2);

      const reconnectedHost = await reconnectPlayer(host, code);
      await expectHostPlayer(reconnectedHost, player2.name);
      await expect(hostBadge(reconnectedHost, host.name)).not.toBeVisible();
    } finally {
      await safeCloseSessions(sessions);
    }
  });

  test('allows the host to hand off control before the session starts', async ({ browser }) => {
    test.setTimeout(90_000);

    const sessions = await Promise.all(
      ['Host', 'Player2', 'Player3'].map((name) => createPlayerSession(browser, name)),
    );

    const [host, player2, player3] = sessions;

    try {
      const code = await createLobby(host, '/cluster', `Cluster Manual Transfer ${Date.now()}`);

      await joinLobby(player2, code);
      await joinLobby(player3, code);

      const hostPage = requirePage(host);
      const player2Page = requirePage(player2);
      const player3Page = requirePage(player3);

      await expect(
        hostPage.locator('form[action$="/host-transfer"] select[name="target_player_id"]'),
      ).toBeVisible({ timeout: 10_000 });

      await transferHost(hostPage, player2.name);

      await expectHostPlayer(hostPage, player2.name);
      await expectHostPlayer(player2Page, player2.name);
      await expectHostPlayer(player3Page, player2.name);

      await expect(
        hostPage.locator('form[action$="/cluster/start"] button:has-text("START CLUSTER")'),
      ).not.toBeVisible();
      await expect(
        player2Page.locator('form[action$="/cluster/start"] button:has-text("START CLUSTER")'),
      ).toBeVisible({ timeout: 10_000 });

      await clickAndWaitForIdle(
        player2Page,
        'form[action$="/cluster/start"] button:has-text("START CLUSTER")',
      );
      await waitForRound([hostPage, player2Page, player3Page], 1);
    } finally {
      await safeCloseSessions(sessions);
    }
  });
});
