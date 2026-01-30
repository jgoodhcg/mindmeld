/**
 * Smoke Tests
 *
 * Basic tests to verify the app is running and key pages load.
 * Run with: npm test (from e2e directory)
 */

import { test, expect } from '@playwright/test';

test.describe('Smoke Tests', () => {
  test('homepage loads', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/mindmeld/i);
  });

  test('can navigate to create lobby', async ({ page }) => {
    await page.goto('/');

    // Look for any create/start button
    const createLink = page.locator('a, button').filter({ hasText: /create|start|new/i }).first();

    if (await createLink.isVisible()) {
      await createLink.click();
      await page.waitForLoadState('networkidle');

      // Should be on a lobby creation or lobby page
      expect(page.url()).toMatch(/lobby|create/i);
    }
  });

  test('invalid lobby shows error or redirect', async ({ page }) => {
    const response = await page.goto('/lobby/INVALID123');

    // Should either show an error page or redirect
    const is404 = response?.status() === 404;
    const hasError = await page.locator('text=/not found|invalid|error/i').isVisible();
    const redirected = page.url() !== page.context().pages()[0]?.url();

    expect(is404 || hasError || redirected).toBeTruthy();
  });
});
