import { expect, test, type Page } from '@playwright/test';

async function joinLobby(page: Page, code: string, nickname: string) {
  await page.goto(`/lobbies/${code}`, { waitUntil: 'networkidle' });
  await expect(page.locator('input[name="nickname"]')).toBeVisible({ timeout: 10000 });
  await page.fill('input[name="nickname"]', nickname);
  await Promise.all([
    page.waitForURL(new RegExp(`/lobbies/${code}$`), { timeout: 10000 }),
    page.click('button:has-text("CONNECT")'),
  ]);
}

async function submitCoordinate(page: Page, x: string, y: string) {
  await selectCoordinate(page, x, y);
  await submitSelectedCoordinate(page);
}

async function selectCoordinate(page: Page, x: string, y: string) {
  const plane = page.locator('#cluster-plane-input');
  await expect(plane).toBeVisible({ timeout: 10000 });

  const box = await plane.boundingBox();
  if (!box) {
    throw new Error('cluster input plane is missing a bounding box');
  }

  const xNum = Number.parseFloat(x);
  const yNum = Number.parseFloat(y);
  const clickX = box.x + (box.width * xNum);
  const clickY = box.y + (box.height * (1 - yNum));
  await page.mouse.click(clickX, clickY);
}

async function submitSelectedCoordinate(page: Page) {
  await Promise.all([
    page.waitForLoadState('networkidle'),
    page.click('button:has-text("SUBMIT POINT")'),
  ]);
}

async function safeCloseContext(context: { close: () => Promise<void> }) {
  try {
    await context.close();
  } catch {
    // Ignore close errors on already-disposed contexts during timeouts/failures.
  }
}

test.describe('Cluster Multiplayer', () => {
  test('creates a lobby and completes a 3-player round', async ({ browser }) => {
    test.setTimeout(120000);

    const hostContext = await browser.newContext({ viewport: { width: 1280, height: 800 } });
    const player2Context = await browser.newContext({ viewport: { width: 1280, height: 800 } });
    const player3Context = await browser.newContext({ viewport: { width: 1280, height: 800 } });

    const hostPage = await hostContext.newPage();
    const player2Page = await player2Context.newPage();
    const player3Page = await player3Context.newPage();

    try {
      await hostPage.goto('/cluster', { waitUntil: 'networkidle' });
      await expect(hostPage.locator('form[action="/lobbies"]')).toBeVisible();

      const sessionName = `Cluster E2E ${Date.now()}`;
      await hostPage.fill('form[action="/lobbies"] input[name="name"]', sessionName);
      await hostPage.fill('form[action="/lobbies"] input[name="nickname"]', 'Host');
      await Promise.all([
        hostPage.waitForURL(/\/lobbies\/[A-Z0-9]{6}$/, { timeout: 10000 }),
        hostPage.click('form[action="/lobbies"] button:has-text("INITIALIZE")'),
      ]);

      const match = hostPage.url().match(/\/lobbies\/([A-Z0-9]{6})$/);
      expect(match).not.toBeNull();
      const code = match![1];

      await joinLobby(player2Page, code, 'Player2');
      await joinLobby(player3Page, code, 'Player3');

      await expect(hostPage.locator('#player-list li')).toHaveCount(3, { timeout: 10000 });
      await expect(hostPage.getByText(/Connected now:\s*3/i)).toBeVisible({ timeout: 10000 });

      const startButton = hostPage.locator('form[action$="/cluster/start"] button:has-text("START CLUSTER")');
      await expect(startButton).toBeVisible({ timeout: 10000 });
      await expect(startButton).toBeEnabled({ timeout: 10000 });
      await Promise.all([
        hostPage.waitForLoadState('networkidle'),
        startButton.click(),
      ]);

      for (const page of [hostPage, player2Page, player3Page]) {
        await expect(page.getByText('ROUND 1')).toBeVisible({ timeout: 10000 });
        await expect(page.locator('#cluster-plane-input')).toBeVisible({ timeout: 10000 });
      }

      await selectCoordinate(player2Page, '0.85', '0.65');
      await expect(player2Page.locator('#cluster-x')).toHaveValue('0.85');
      await expect(player2Page.locator('#cluster-y')).toHaveValue('0.65');

      await submitCoordinate(hostPage, '0.10', '0.20');
      await expect(hostPage.getByText('COORDINATE LOCKED')).toBeVisible({ timeout: 10000 });
      await expect(hostPage.getByText(/1\s*\/\s*3 submitted/i)).toBeVisible({ timeout: 10000 });

      await expect(player2Page.locator('#cluster-x')).toHaveValue('0.85');
      await expect(player2Page.locator('#cluster-y')).toHaveValue('0.65');
      await expect(player2Page.locator('#cluster-coordinate-readout')).toHaveText('Selected point: (0.70, 0.30)');
      await submitSelectedCoordinate(player2Page);
      await expect(hostPage.getByText(/2\s*\/\s*3 submitted/i)).toBeVisible({ timeout: 10000 });

      await submitCoordinate(player3Page, '0.50', '0.45');

      for (const page of [hostPage, player2Page, player3Page]) {
        await expect(page.getByText('CENTROID REVEALED')).toBeVisible({ timeout: 10000 });
        await expect(page.locator('[title="Centroid target"]')).toBeVisible({ timeout: 10000 });
        await expect(page.getByText('Avg/round')).toBeVisible({ timeout: 10000 });
        await expect(page.locator('#game-content').getByText('Host', { exact: true })).toBeVisible();
        await expect(page.locator('#game-content').getByText('Player2', { exact: true })).toBeVisible();
        await expect(page.locator('#game-content').getByText('Player3', { exact: true })).toBeVisible();
      }

      const nextRoundButton = hostPage.locator('form[action$="/cluster/next"] button');
      await expect(nextRoundButton).toBeVisible({ timeout: 10000 });
      await Promise.all([
        hostPage.waitForLoadState('networkidle'),
        nextRoundButton.click(),
      ]);

      for (const page of [hostPage, player2Page, player3Page]) {
        await expect(page.getByText('ROUND 2')).toBeVisible({ timeout: 10000 });
        await expect(page.locator('#cluster-plane-input')).toBeVisible({ timeout: 10000 });
      }
    } finally {
      await safeCloseContext(hostContext);
      await safeCloseContext(player2Context);
      await safeCloseContext(player3Context);
    }
  });
});
