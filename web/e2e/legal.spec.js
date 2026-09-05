// @ts-check
const { test, expect } = require('@playwright/test');

// Legal pages: Terms of Service and Privacy Policy are static HTML,
// reachable from the landing-page footer.

test.describe('Legal pages', () => {
  for (const [path, title] of [
    ['/terms.html', 'Terms of Service'],
    ['/privacy.html', 'Privacy Policy'],
  ]) {
    test(`${title} renders and links home`, async ({ page }) => {
      await page.goto(path);
      await expect(page.locator('.legal-page h1')).toHaveText(title);
      await expect(page.locator('.legal-page .dateline')).toContainText('Effective Date');

      await page.locator('.site-name a').click();
      await expect(page).toHaveURL(/\/$/);
      await expect(page.locator('#sign-form')).toBeVisible();
    });
  }

  test('landing-page footer links to both policies', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('footer a[href="terms.html"]')).toHaveText('Terms');
    await expect(page.locator('footer a[href="privacy.html"]')).toHaveText('Privacy');
    await expect(page.locator('footer a[href="https://github.com/grocky/qr-codes"]')).toHaveText('Open source');
  });
});
