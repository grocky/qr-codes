// @ts-check
const { test, expect } = require('@playwright/test');

// End-to-end smoke: form → WASM render → preview → downloads.
// Build web/main.wasm first (make build-wasm).

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  // WASM ready: the loading indicator hides once wifiSignReady fires.
  await expect(page.locator('#loading')).toBeHidden({ timeout: 15_000 });
});

test('renders a live preview from credentials', async ({ page }) => {
  await expect(page.locator('#preview-hint')).toBeVisible();

  await page.fill('#ssid', 'Potomac Poker');
  await page.fill('#password', 'all-in-2026');

  const svg = page.locator('#preview svg');
  await expect(svg).toBeVisible();
  await expect(svg.locator('text', { hasText: 'WI-FI' })).toBeVisible();
  // Vector QR path present.
  await expect(page.locator('#preview g[shape-rendering="crispEdges"] path')).toHaveCount(1);
});

test('shows per-field errors and disables downloads', async ({ page }) => {
  await page.fill('#ssid', 'net');
  await page.fill('#password', ''); // WPA requires a password
  await expect(page.locator('.field-error[data-field="password"]')).not.toBeEmpty();
  await expect(page.locator('#dl-svg')).toBeDisabled();
});

test('downloads SVG, PNG, and PDF', async ({ page }) => {
  await page.fill('#ssid', 'Potomac Poker');
  await page.fill('#password', 'all-in-2026');
  await expect(page.locator('#preview svg')).toBeVisible();

  for (const [button, suffix, magic] of [
    ['#dl-svg', '.svg', '<svg'],
    ['#dl-png', '.png', '\x89PNG'],
    ['#dl-pdf', '.pdf', '%PDF'],
  ]) {
    const downloadPromise = page.waitForEvent('download');
    await page.click(button);
    const download = await downloadPromise;
    expect(download.suggestedFilename()).toMatch(new RegExp(`\\${suffix}$`));

    const stream = await download.createReadStream();
    const head = await new Promise((resolve, reject) => {
      stream.once('readable', () => resolve(stream.read(8).toString('latin1')));
      stream.once('error', reject);
    });
    expect(head.startsWith(magic), `${suffix} magic bytes`).toBeTruthy();
  }
});

test('customization reaches the preview', async ({ page }) => {
  await page.fill('#ssid', 'net');
  await page.fill('#password', 'password');
  await page.fill('#headline', 'GUEST WI-FI');

  const preview = page.locator('#preview');
  await expect(preview.locator('svg text', { hasText: 'GUEST WI-FI' })).toBeVisible();
  await expect(preview.locator('#suit-row')).toHaveCount(1); // default ornament

  // Switching to a tagline reveals the field and replaces the suits.
  await expect(page.locator('#tagline-label')).toBeHidden();
  await page.selectOption('#ornament', 'tagline');
  await expect(page.locator('#tagline-label')).toBeVisible();
  await expect(page.locator('.field-error[data-field="tagline"]')).not.toBeEmpty();
  await page.fill('#tagline', 'Good luck, have fun');
  await expect(preview.locator('#ornament-tagline')).toBeVisible();
  await expect(preview.locator('#suit-row')).toHaveCount(0);

  await page.selectOption('#ornament', 'none');
  await expect(preview.locator('#ornament-tagline')).toHaveCount(0);
  await expect(preview.locator('#suit-row')).toHaveCount(0);
});
