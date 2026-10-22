import { test, expect } from '@playwright/test';

test('muestra las cuentas de prueba en el inicio de sesión', async ({ page }) => {
  await page.goto('/login');
  await expect(page.getByText('Cuentas de prueba')).toBeVisible();
  await expect(page.getByText('admin@upse.edu.ec', { exact: false })).toBeVisible();
  await expect(page.getByText('cliente@upse.edu.ec', { exact: false })).toBeVisible();
  await expect(page.getByText('123456123456')).toBeVisible();
});

test('instala service worker, conserva carrito y catálogo offline sin cachear datos privados', async ({ page, context }) => {
  await page.route('**/api/**', async route => {
    const path = new URL(route.request().url()).pathname;
    if (path.endsWith('/login')) {
      const claims = Buffer.from(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 })).toString('base64url');
      return route.fulfill({ json: { email: 'pwa@test.local', nombre: 'Prueba', rol: 'cliente', token: 'test.' + claims + '.signature' } });
    }
    if (path.endsWith('/productos')) return route.fulfill({ json: [{ id: 1, nombre: 'Producto PWA', precio: 10, stock: 10, img: '/placeholder.svg' }] });
    return route.fulfill({ json: [] });
  });
  await page.goto('/login');
  await page.getByPlaceholder('admin@upse.edu.ec').fill('pwa@test.local');
  await page.locator('input[type=password]').fill('123456');
  await page.getByRole('button', { name: 'Iniciar Sesión' }).click();
  await page.goto('/catalogo');
  await expect(page.getByText('Producto PWA')).toBeVisible();
  await page.getByRole('button', { name: 'Añadir al Carrito' }).click();
  await page.evaluate(() => navigator.serviceWorker.ready);
  await page.reload();
  await expect.poll(() => page.evaluate(() => Boolean(navigator.serviceWorker.controller))).toBe(true);
  await context.setOffline(true);
  await page.unroute('**/api/**');
  await page.reload();
  await expect(page.getByText('Producto PWA')).toBeVisible();
  await expect(page.getByText('Sin conexión.', { exact: false })).toBeVisible();
  await page.goto('/carrito');
  await expect(page.getByText('Producto PWA')).toBeVisible();
  await page.goto('/checkout');
  await expect(page.getByText('Resumen del Pedido')).toBeVisible();
  for (const [label, value] of [['Nombre completo *', 'Cliente PWA'], ['Dirección *', 'Centro'], ['Ciudad *', 'La Libertad'], ['Teléfono *', '0999999999']]) {
    await page.getByLabel(label, { exact: true }).fill(value);
  }
  await page.getByRole('button', { name: 'Confirmar Pedido' }).click();
  await expect(page.getByText('Conéctate a internet para confirmar la compra.', { exact: false })).toBeVisible();
  expect(await page.evaluate(() => JSON.parse(localStorage.getItem('multicatalogo_carrito_pwa@test.local') ?? '[]').length)).toBe(1);
  const cached = await page.evaluate(async () => {
    const keys = await caches.keys();
    return (await Promise.all(keys.map(async key => (await (await caches.open(key)).keys()).map(r => r.url)))).flat();
  });
  expect(cached.some(url => new URL(url).pathname.startsWith('/api/'))).toBe(false);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
});

test('manifest e iconos disponibles y primera visita sin catálogo muestra un error claro', async ({ page }) => {
  const manifest = await page.request.get('/manifest.webmanifest');
  const data = await manifest.json();
  expect(data.display).toBe('standalone');
  for (const icon of data.icons) expect((await page.request.get(icon.src)).ok()).toBe(true);
  const fotos = ['auriculares-orbit', 'teclado-nova', 'mouse-pulse', 'parlante-echo', 'bateria-volt', 'lampara-aura', 'botella-atlas', 'mochila-urban', 'organizador-zen', 'difusor-bruma', 'banda-active', 'mat-yoga', 'kit-lumen', 'protector-nube', 'reloj-vector', 'lentes-prisma', 'libreta-terra', 'taza-andina', 'hub-nexus', 'soporte-flex'];
  for (const foto of fotos) expect((await page.request.get(`/products/images/${foto}.webp`)).ok()).toBe(true);
  await page.route('**/api/**', route => route.abort());
  await page.goto('/catalogo');
  await expect(page.getByText('No se pudo conectar con el servidor.', { exact: false })).toBeVisible();
});

test('muestra la copia del catálogo sin esperar a que una conexión inestable termine', async ({ page }) => {
  await page.goto('/login');
  await page.evaluate(() => {
    const date = Date.now();
    localStorage.setItem('multicatalogo_catalogo', JSON.stringify({
      date,
      products: [{ id: 9, nombre: 'Producto guardado', precio: 15, stock: 4, img: '/placeholder.svg' }],
    }));
    localStorage.setItem('multicatalogo_catalogo_fecha', String(date));
  });
  await page.route('**/api/productos', () => new Promise(() => {}));

  await page.goto('/catalogo');

  await expect(page.getByText('Producto guardado')).toBeVisible({ timeout: 1000 });
  await expect(page.getByText('Copia disponible 24 h')).toBeVisible();
});

test('cambiar cuenta no mezcla carritos y una sesión vencida permite consultar el catálogo público', async ({ page }) => {
  await page.goto('/login');
  await page.evaluate(() => {
    localStorage.setItem('multicatalogo_carrito_a@test.local', JSON.stringify([{ id: 1, nombre: 'Carrito A', precio: 10, cantidad: 1, stock: 5, img: '/placeholder.svg' }]));
    const payload = btoa(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 }));
    localStorage.setItem('multicatalogo_sesion', JSON.stringify({ email: 'a@test.local', token: 'a.' + payload + '.c', rol: 'cliente' }));
  });
  await page.goto('/carrito');
  await expect(page.getByText('Carrito A')).toBeVisible();
  await page.evaluate(() => {
    const payload = btoa(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 }));
    localStorage.setItem('multicatalogo_sesion', JSON.stringify({ email: 'b@test.local', token: 'a.' + payload + '.c', rol: 'cliente' }));
  });
  await page.reload();
  await expect(page.getByText('Tu carrito está vacío actualmente.')).toBeVisible();
  await page.evaluate(() => {
    const payload = btoa(JSON.stringify({ exp: 1 }));
    localStorage.setItem('multicatalogo_sesion', JSON.stringify({ email: 'a@test.local', token: 'a.' + payload + '.c', rol: 'cliente' }));
  });
  await page.reload();
  await expect(page.getByRole('button', { name: 'Iniciar Sesión' })).toBeVisible();
  await page.goto('/catalogo');
  await expect(page.getByRole('heading', { name: 'Catálogo de Productos' })).toBeVisible();
});
