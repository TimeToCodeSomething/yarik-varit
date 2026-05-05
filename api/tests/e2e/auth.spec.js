const { test, expect } = require('@playwright/test');

// Берём логин/пароль из переменных окружения
const ADMIN_USER = process.env.TEST_ADMIN_USER || 'admin';
const ADMIN_PASS = process.env.TEST_ADMIN_PASS || 'password';

test.describe('Авторизация', () => {

  test('кнопка Войти открывает модальное окно', async ({ page }) => {
    await page.goto('');

    await page.click('#btnNavAuth');
    const modal = page.locator('#authModal');
    await expect(modal).toHaveClass(/open/);
  });

  test('модальное окно закрывается по крестику', async ({ page }) => {
    await page.goto('');

    await page.click('#btnNavAuth');
    await page.click('.modal-close');
    const modal = page.locator('#authModal');
    await expect(modal).not.toHaveClass(/open/);
  });

  test('неверный пароль показывает alert', async ({ page }) => {
    await page.goto('');

    // Перехватываем alert
    let alertMessage = '';
    page.on('dialog', async dialog => {
      alertMessage = dialog.message();
      await dialog.dismiss();
    });

    await page.click('#btnNavAuth');
    await page.fill('#loginEmail', ADMIN_USER);
    await page.fill('#loginPass', 'wrongpassword');
    await page.click('.btn-auth-submit');

    await page.waitForTimeout(500);
    expect(alertMessage).toContain('Неверный');
  });

  test('успешный логин показывает кнопку Профиль', async ({ page }) => {
    await page.goto('');

    await page.click('#btnNavAuth');
    await page.fill('#loginEmail', ADMIN_USER);
    await page.fill('#loginPass', ADMIN_PASS);
    await page.click('.btn-auth-submit');

    await expect(page.locator('#btnNavProfile')).toBeVisible({ timeout: 3000 });
    await expect(page.locator('#btnNavAuth')).toBeHidden();
  });

  test('после логина открывается профиль', async ({ page }) => {
    await page.goto('');

    await page.click('#btnNavAuth');
    await page.fill('#loginEmail', ADMIN_USER);
    await page.fill('#loginPass', ADMIN_PASS);
    await page.click('.btn-auth-submit');

    await page.click('#btnNavProfile');
    const panel = page.locator('#profilePanel');
    await expect(panel).toHaveClass(/open/);
  });

  test('выход из аккаунта возвращает кнопку Войти', async ({ page }) => {
    await page.goto('');

    // Логинимся
    await page.click('#btnNavAuth');
    await page.fill('#loginEmail', ADMIN_USER);
    await page.fill('#loginPass', ADMIN_PASS);
    await page.click('.btn-auth-submit');

    // Открываем профиль и выходим
    await page.click('#btnNavProfile');
    await page.click('.btn-logout');

    await expect(page.locator('#btnNavAuth')).toBeVisible();
    await expect(page.locator('#btnNavProfile')).toBeHidden();
  });

  test('переключение на вкладку Аккаунт работает', async ({ page }) => {
    await page.goto('');

    await page.click('#btnNavAuth');
    await page.click('button:has-text("Аккаунт")');

    await expect(page.locator('#tab-register')).toBeVisible();
    await expect(page.locator('#tab-login')).toBeHidden();
  });

});