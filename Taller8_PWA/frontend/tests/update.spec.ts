import { test, expect } from '@playwright/test';

test('actualizar requiere decisión del usuario y conserva el carrito', async ({ page }) => {
  await page.goto('/catalogo');
  await page.evaluate(async () => {
    const payload = btoa(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 }));
    localStorage.setItem('multicatalogo_sesion', JSON.stringify({ email: 'update@test.local', token: 'a.' + payload + '.c', rol: 'cliente' }));
    localStorage.setItem('multicatalogo_carrito_update@test.local', JSON.stringify([{ id: 1, nombre: 'Carrito guardado', precio: 10, cantidad: 1, stock: 5, img: '/placeholder.svg' }]));
    await navigator.serviceWorker.ready;
  });
  await page.reload();
  await expect.poll(() => page.evaluate(() => Boolean(navigator.serviceWorker.controller))).toBe(true);
  await page.goto('/checkout');
  await page.request.post('/__test/update');
  await page.evaluate(async () => { const registration = await navigator.serviceWorker.getRegistration(); await registration!.update(); });
  await expect(page.getByText('Hay una nueva versión.', { exact: false })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Actualizar desde catálogo o carrito' })).toBeDisabled();
  await page.goto('/carrito');
  await expect(page.getByText('Carrito guardado')).toBeVisible();
  await page.getByRole('button', { name: 'Actualizar desde catálogo o carrito' }).click();
  await expect(page.getByText('Hay una nueva versión.', { exact: false })).not.toBeVisible();
  await expect(page.getByText('Carrito guardado')).toBeVisible();
});
