import { expect, type Browser, type BrowserContext, type Page } from '@playwright/test';

export type PlayerSession = {
  name: string;
  context: BrowserContext;
  page: Page | null;
};

const viewport = { width: 1280, height: 800 };
const defaultDisconnectGracePeriodMs = 25_000;

function escapeRegex(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

export async function createPlayerSession(browser: Browser, name: string): Promise<PlayerSession> {
  const context = await browser.newContext({ viewport });
  const page = await context.newPage();
  return { name, context, page };
}

export async function ensurePage(session: PlayerSession): Promise<Page> {
  if (!session.page || session.page.isClosed()) {
    session.page = await session.context.newPage();
  }
  return session.page;
}

export function requirePage(session: PlayerSession): Page {
  if (!session.page || session.page.isClosed()) {
    throw new Error(`player ${session.name} does not have an active page`);
  }
  return session.page;
}

export async function createLobby(
  host: PlayerSession,
  gamePath: '/trivia' | '/cluster',
  lobbyName: string,
): Promise<string> {
  const page = requirePage(host);
  await page.goto(gamePath, { waitUntil: 'networkidle' });
  await expect(page.locator('form[action="/lobbies"]')).toBeVisible({ timeout: 10_000 });
  await page.fill('form[action="/lobbies"] input[name="name"]', lobbyName);
  await page.fill('form[action="/lobbies"] input[name="nickname"]', host.name);
  await Promise.all([
    page.waitForURL(/\/lobbies\/[A-Z0-9]{6}$/, { timeout: 10_000 }),
    page.click('form[action="/lobbies"] button'),
  ]);

  const match = page.url().match(/\/lobbies\/([A-Z0-9]{6})$/);
  if (!match) {
    throw new Error(`failed to extract lobby code from ${page.url()}`);
  }

  return match[1];
}

export async function joinLobby(session: PlayerSession, code: string): Promise<Page> {
  const page = await ensurePage(session);
  await page.goto(`/lobbies/${code}`, { waitUntil: 'networkidle' });
  await expect(page.locator('input[name="nickname"]')).toBeVisible({ timeout: 10_000 });
  await page.fill('input[name="nickname"]', session.name);
  await Promise.all([
    page.waitForURL(new RegExp(`/lobbies/${code}$`), { timeout: 10_000 }),
    page.click('button:has-text("CONNECT")'),
  ]);
  await expect(page.locator('#player-list')).toBeVisible({ timeout: 10_000 });
  return page;
}

export async function reconnectPlayer(session: PlayerSession, code: string): Promise<Page> {
  const page = await ensurePage(session);
  await page.goto(`/lobbies/${code}`, { waitUntil: 'networkidle' });
  await expect(page.locator('#player-list')).toBeVisible({ timeout: 10_000 });
  return page;
}

export async function disconnectPlayer(session: PlayerSession): Promise<void> {
  if (!session.page || session.page.isClosed()) {
    session.page = null;
    return;
  }
  await session.page.close();
  session.page = null;
}

export function connectedPages(sessions: PlayerSession[]): Page[] {
  return sessions
    .map((session) => session.page)
    .filter((page): page is Page => page !== null && !page.isClosed());
}

export async function expectLobbyPlayerCount(page: Page, count: number): Promise<void> {
  await expect(page.locator('#player-list li')).toHaveCount(count, { timeout: 10_000 });
}

export function playerRow(page: Page, playerName: string) {
  return page
    .locator('#player-list li')
    .filter({
      has: page.locator('span.text-text', {
        hasText: new RegExp(`^${escapeRegex(playerName)}$`),
      }),
    })
    .first();
}

export function hostBadge(page: Page, playerName: string) {
  return playerRow(page, playerName).locator('span.text-xs').filter({ hasText: 'Host' }).first();
}

export async function expectHostPlayer(
  page: Page,
  playerName: string,
  timeout = disconnectGracePeriodMs() + 10_000,
): Promise<void> {
  await expect(hostBadge(page, playerName)).toBeVisible({ timeout });
}

export async function expectPlayerState(
  page: Page,
  playerName: string,
  state: 'connected' | 'reconnecting' | 'disconnected',
  timeout = disconnectGracePeriodMs() + 10_000,
): Promise<void> {
  const row = playerRow(page, playerName);
  await expect(row).toBeVisible({ timeout });

  if (state === 'connected') {
    await expect(row).not.toContainText('Reconnecting', { timeout });
    await expect(row).not.toContainText('Disconnected', { timeout });
    return;
  }

  await expect(row).toContainText(state === 'reconnecting' ? 'Reconnecting' : 'Disconnected', {
    timeout,
  });
}

export async function clickAndWaitForIdle(page: Page, locator: string): Promise<void> {
  await Promise.all([page.waitForLoadState('networkidle'), page.click(locator)]);
}

export async function transferHost(page: Page, playerName: string): Promise<void> {
  await page.selectOption('form[action$="/host-transfer"] select[name="target_player_id"]', {
    label: playerName,
  });
  page.once('dialog', (dialog) => dialog.accept());
  await clickAndWaitForIdle(page, 'form[action$="/host-transfer"] button:has-text("TRANSFER HOST")');
}

export function disconnectGracePeriodMs(): number {
  const raw = process.env.DISCONNECT_GRACE_PERIOD;
  if (!raw) {
    return defaultDisconnectGracePeriodMs;
  }

  const trimmed = raw.trim().toLowerCase();
  if (/^\d+$/.test(trimmed)) {
    return Number.parseInt(trimmed, 10);
  }

  const match = trimmed.match(/^(\d+)(ms|s|m)$/);
  if (!match) {
    return defaultDisconnectGracePeriodMs;
  }

  const value = Number.parseInt(match[1], 10);
  switch (match[2]) {
    case 'ms':
      return value;
    case 's':
      return value * 1_000;
    case 'm':
      return value * 60_000;
    default:
      return defaultDisconnectGracePeriodMs;
  }
}

export async function safeCloseSession(session: PlayerSession): Promise<void> {
  try {
    await session.context.close();
  } catch {
    // Ignore cleanup errors during failed test teardown.
  }
}

export async function safeCloseSessions(sessions: PlayerSession[]): Promise<void> {
  for (const session of sessions) {
    await safeCloseSession(session);
  }
}
