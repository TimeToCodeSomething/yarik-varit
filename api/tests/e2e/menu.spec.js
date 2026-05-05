const { test, expect } = require('@playwright/test');

test.describe('Меню', () => {

  test('страница открывается', async ({ page }) => {
    await page.goto('');
    await expect(page).toHaveTitle(/Ярик Варит/);
  });

  test('меню загружается из БД', async ({ page }) => {
    await page.goto('');

    // Ждём пока исчезнет статика и появятся данные из API
    const menuItems = page.locator('.menu-item');
    await expect(menuItems.first()).toBeVisible({ timeout: 5000 });

    const count = await menuItems.count();
    expect(count).toBeGreaterThan(0);
  });

  test('переключение категорий работает', async ({ page }) => {
    await page.goto('');

    // Кликаем на "Молочные"
    await page.click('button:has-text("Молочные")');
    const panel = page.locator('#panel-milk');
    await expect(panel).toHaveClass(/active/);
  });

  test('категория эспрессо активна по умолчанию', async ({ page }) => {
    await page.goto('');

    const panel = page.locator('#panel-espresso');
    await expect(panel).toHaveClass(/active/);
  });

});