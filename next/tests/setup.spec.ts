import { test, expect } from '@playwright/test';

test('platform choice keeps downloads and Soloist instructions in sync', async ({ page }) => {
  await page.goto('/docs');
  await page.getByRole('button', { name: 'Mac Apple Silicon' }).click();
  await expect(page.getByRole('button', { name: 'Mac Apple Silicon' })).toHaveAttribute('aria-pressed', 'true');
  await page.getByText('Already downloaded it, or prefer a manual archive?').click();
  await expect(page.getByRole('link', { name: 'Download Apple Silicon archive' })).toHaveAttribute('href', /_darwin_arm64\.tar\.gz$/);
  await expect(page.locator('#audio')).toContainText('soloist_release_osx_arm64.tar.gz');
  await expect(page.locator('#audio')).not.toContainText('soloist_release_x86_64.tar.gz');
  await page.getByRole('button', { name: 'Linux x86_64' }).click();
  await expect(page.getByRole('link', { name: 'Download Linux x86_64 archive' })).toHaveAttribute('href', /_linux_amd64\.tar\.gz$/);
  await expect(page.locator('#audio')).toContainText('soloist_release_x86_64.tar.gz');
  await expect(page.locator('#audio')).toContainText('setup --soloist-path');
  await expect(page.locator('#audio')).not.toContainText('YOUR_PUBLIC_CLIENT_ID');
  const command = await page.locator('.code').filter({ hasText: 'Download and install' }).locator('pre').innerText();
  expect(command).toMatch(/--version \d+\.\d+\.\d+ \\\n/);
  expect(command).toContain('&& \\\n  bash ./spotify-tui-');
  expect(command).not.toContain('\\n');
});

test('all four steps are readable with working checkpoints and active navigation', async ({ page }) => {
  const errors: string[] = [];
  page.on('pageerror', error => errors.push(error.message));
  await page.goto('/docs');
  await expect(page.getByRole('navigation', { name: 'Main navigation' }).getByRole('link', { name: 'Docs', exact: true })).toHaveAttribute('aria-current', 'page');
  await expect(page.locator('.setup-checkpoint')).toHaveCount(4);
  for (const [id, link] of [['sign-in','02 Sign in'], ['audio','03 Pair'], ['play','04 Play']]) {
    const anchor = page.getByRole('navigation', { name: 'Setup sections' }).getByRole('link', { name: link });
    await anchor.focus();
    await page.keyboard.press('Enter');
    await expect(page).toHaveURL(new RegExp('#'+id+'$'));
    await expect(anchor).toHaveAttribute('aria-current', 'location');
    await expect(page.locator('#'+id)).toHaveCSS('opacity', '1');
    await expect(page.locator('#'+id)).toHaveCSS('transform', 'none');
  }
  await expect(page.getByRole('button', { name: 'Copy redirect uri' }).locator('..').locator('..').locator('pre')).toHaveText('http://127.0.0.1:43821/callback');
  await page.getByRole('navigation', { name: 'Main navigation' }).getByRole('link', { name: 'Changelog', exact: true }).click();
  await expect(page).toHaveURL(/\/changelog\/?$/);
  await expect(page.getByRole('navigation', { name: 'Main navigation' }).getByRole('link', { name: 'Changelog', exact: true })).toHaveAttribute('aria-current', 'page');
  expect(errors).toEqual([]);
});

test('clipboard denial selects real text and explains manual copy', async ({ page }) => {
  await page.addInitScript(() => Object.defineProperty(navigator, 'clipboard', {
    configurable: true,
    value: { writeText: async () => { throw new Error('Clipboard denied in test'); } },
  }));
  await page.goto('/docs');
  await page.getByRole('button', { name: 'Copy redirect uri' }).click();
  await expect(page.getByText('Text selected. Press Ctrl+C or ⌘C to copy.')).toBeVisible();
  expect(await page.evaluate(() => window.getSelection()?.toString())).toBe('http://127.0.0.1:43821/callback');
  await expect(page.getByRole('button', { name: 'Copy redirect uri' })).toHaveText('Selected');
});

test('clipboard success copies the exact redirect and reports success', async ({ page, context, browserName }) => {
  test.skip(browserName !== 'chromium', 'Clipboard permissions differ in WebKit; denial path is covered on both engines');
  await context.grantPermissions(['clipboard-read', 'clipboard-write']);
  await page.goto('/docs');
  await page.getByRole('button', { name: 'Copy redirect uri' }).click();
  await expect(page.getByRole('button', { name: 'Copy redirect uri' })).toContainText('Copied');
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe('http://127.0.0.1:43821/callback');
});

test('small-screen setup and reduced-motion reading remain usable', async ({ page }) => {
  await page.setViewportSize({ width: 320, height: 700 });
  await page.emulateMedia({ reducedMotion: 'reduce' });
  await page.goto('/docs');
  await page.getByRole('button', { name: 'Mac Apple Silicon' }).click();
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  for (const id of ['install','sign-in','audio','play']) {
    await expect(page.locator('#'+id)).toHaveCSS('opacity','1');
    await expect(page.locator('#'+id)).toHaveCSS('transform','none');
  }
  await expect(page.locator('.setup-route a').first()).toHaveCSS('transition-duration','0s');
  await page.goto('/');
  await page.getByRole('button', { name: 'Pause animations', exact: true }).click();
  await expect(page.getByRole('button', { name: 'Resume animations', exact: true })).toHaveAttribute('aria-pressed','true');
  await expect(page.locator('.footer-wave-track')).toHaveCSS('animation-play-state','paused');
});

test('setup remains readable without JavaScript', async ({ browser }) => {
  const context = await browser.newContext({ javaScriptEnabled: false });
  const page = await context.newPage();
  await page.goto('/docs');
  await expect(page.locator('#play')).toContainText('Pick a song. Hit Enter.');
  await expect(page.getByRole('link', { name: 'Apple Silicon Mac', exact: true })).toHaveAttribute('href', /_darwin_arm64\.tar\.gz$/);
  await expect(page.locator('#install')).toHaveCSS('opacity','1');
  await context.close();
});
