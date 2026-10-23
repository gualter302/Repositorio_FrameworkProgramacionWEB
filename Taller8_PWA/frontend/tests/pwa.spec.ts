import { test, expect } from '@playwright/test';

test('muestra las cuentas de prueba en el inicio de sesión', async ({ page }) => {
  await page.goto('/login');
  await expect(page.getByText('Cuentas de prueba')).toBeVisible();
  await expect(page.getByText('admin@upse.edu.ec', { exact: false })).toBeVisible();
  await expect(page.getByText('cliente@upse.edu.ec', { exact: false })).toBeVisible();
  await expect(page.getByText('123456123456')).toBeVisible();
  await page.getByRole('button', { name: 'Usar cuenta de administrador' }).click();
  await expect(page.getByPlaceholder('admin@upse.edu.ec')).toHaveValue('admin@upse.edu.ec');
  await expect(page.locator('input[type=password]')).toHaveValue('123456123456');
  await page.getByRole('button', { name: 'Usar cuenta de cliente' }).click();
  await expect(page.getByPlaceholder('admin@upse.edu.ec')).toHaveValue('cliente@upse.edu.ec');
  await expect(page.locator('input[type=password]')).toHaveValue('123456');
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
});

test('recorre registro, red, pedidos, comisiones, catálogo y compra del cliente', async ({ page }) => {
  let registro: Record<string, string> | null = null;
  let compraRegistrada = false;
  const claims = Buffer.from(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 })).toString('base64url');
  const producto = { id: 101, nombre: 'Auriculares Orbit', descripcion: 'Producto de prueba', categoria: 'Tecnología', precio: 39.90, img: '/products/images/auriculares-orbit.webp', galeria: [], stock: 1, activo: true };
  const red = [{
    id: 6, nombre: 'Cliente Demo', codigo: 'DEMO', nivel: 4, ventas: 208.80, comisiones: 8.36, parentId: 4,
    hijos: [
      { id: 7, nombre: 'Sofía Rivera', codigo: 'DEMO-SOFIA', nivel: 5, ventas: 66.85, comisiones: 2.01, parentId: 6, hijos: [{ id: 9, nombre: 'Valentina Cruz', codigo: 'DEMO-VALE', nivel: 6, ventas: 0, comisiones: 0, parentId: 7, hijos: [] }] },
      { id: 8, nombre: 'Mateo Salazar', codigo: 'DEMO-MATEO', nivel: 5, ventas: 0, comisiones: 0, parentId: 6, hijos: [] },
    ],
  }];
  const pedido = { id: 501, numero: 'DEMO-2026-005', total: 92.45, estado: 'confirmado', fecha: new Date().toISOString(), cliente: { nombre: 'Cliente Demo', email: 'cliente@upse.edu.ec', direccion: 'La Milina', ciudad: 'Salinas', telefono: '0991234567' }, usuarioEmail: 'cliente@upse.edu.ec', items: [{ productoId: 101, nombre: 'Auriculares Orbit', cantidad: 1, precio: 39.90, subtotal: 39.90 }] };

  await page.route('**/api/**', route => {
    const url = new URL(route.request().url());
    const method = route.request().method();
    if (url.pathname === '/api/register' && method === 'POST') {
      registro = route.request().postDataJSON() as Record<string, string>;
      return route.fulfill({ status: 201, json: { email: 'nueva@upse.edu.ec', nombre: 'Nueva Persona', rol: 'cliente', token: `test.${claims}.firma` } });
    }
    if (url.pathname === '/api/red') return route.fulfill({ json: red });
    if (url.pathname === '/api/pedidos') return route.fulfill({ json: [pedido] });
    if (url.pathname === '/api/comisiones') return route.fulfill({ json: [{ id: 81, pedidoId: 501, numeroPedido: 'PRUEBA-RED-001', referidoId: 6, referido: 'Cliente Demo', nivel: 1, porcentaje: 0.05, monto: 3.14, estado: 'pagada', fecha: new Date().toISOString() }] });
    if (url.pathname === '/api/productos') return route.fulfill({ json: [producto] });
    if (url.pathname === '/api/checkout' && method === 'POST') {
      compraRegistrada = true;
      return route.fulfill({ status: 201, json: { ...pedido, id: 900, numero: 'MC-000900', total: 39.90, items: pedido.items } });
    }
    return route.fulfill({ status: 404, json: { error: `Ruta simulada no encontrada: ${method} ${url.pathname}` } });
  });

  await page.goto('/registro?ref=demo');
  await page.getByLabel('Nombre completo').fill('Nueva Persona');
  await page.getByLabel('Correo electrónico').fill('nueva@upse.edu.ec');
  await page.getByLabel('Contraseña').fill('123456');
  await expect(page.getByLabel('Código de referido (opcional)')).toHaveValue('DEMO');
  await page.getByRole('button', { name: 'Registrarme' }).click();
  await expect(page).toHaveURL(/\/mi-red$/);
  expect(registro).toMatchObject({ ref: 'DEMO', email: 'nueva@upse.edu.ec' });
  await expect(page.locator('tbody')).toContainText('Sofía Rivera');
  await expect(page.locator('tbody')).not.toContainText('Cliente Demo');

  await page.goto('/');
  await expect(page.getByText('$208.80')).toBeVisible();
  await expect(page.getByText('$8.36')).toBeVisible();
  await expect(page.locator('.metric-card').filter({ hasText: 'Referidos activos' })).toContainText('3');
  await page.goto('/pedidos');
  await expect(page.getByText('DEMO-2026-005')).toBeVisible();
  await page.goto('/comisiones');
  await expect(page.getByText('PRUEBA-RED-001')).toBeVisible();

  await page.goto('/catalogo');
  await page.getByRole('button', { name: 'Añadir al carrito' }).click();
  await expect(page.getByRole('button', { name: 'Stock agregado' })).toBeDisabled();
  await page.goto('/checkout');
  await page.getByLabel('Nombre completo *', { exact: true }).fill('Nueva Persona');
  await page.getByLabel('Dirección *', { exact: true }).fill('Barrio Central');
  await page.getByLabel('Ciudad *', { exact: true }).fill('Santa Elena');
  await page.getByLabel('Teléfono *', { exact: true }).fill('0999999999');
  await page.getByRole('button', { name: 'Confirmar Pedido' }).click();
  await expect(page).toHaveURL(/\/confirmacion$/);
  await expect(page.getByText('MC-000900')).toBeVisible();
  expect(compraRegistrada).toBe(true);
  await page.goto('/carrito');
  await expect(page.getByText('Tu carrito está vacío actualmente.')).toBeVisible();
});

test('los reportes cortos aprovechan el espacio sin mostrar desplazamiento horizontal', async ({ page }) => {
  await page.route('**/api/admin/reportes/general', route => route.fulfill({ json: {
    totalVentas: 3563.85, totalPedidos: 12, ticketPromedio: 296.99, totalComisiones: 92.40, pendientePagar: 48.20,
    productosPrincipales: [{ nombre: 'Banda Active', cantidad: 2, ventas: 29.80 }, { nombre: 'Libreta Terra', cantidad: 2, ventas: 25.20 }],
    referidosPrincipales: [{ nombre: 'Ana Torres', ventas: 1567.95, comisiones: 31.36 }, { nombre: 'Administrador', ventas: 367.95, comisiones: 7.36 }],
  } }));
  await page.goto('/login');
  await page.evaluate(() => {
    const payload = btoa(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 }));
    localStorage.setItem('multicatalogo_sesion', JSON.stringify({ email: 'admin@upse.edu.ec', token: `a.${payload}.c`, rol: 'admin' }));
  });

  for (const viewport of [{ width: 1440, height: 900 }, { width: 375, height: 812 }]) {
    await page.setViewportSize(viewport);
    await page.goto('/admin/reportes');
    await expect(page.getByRole('heading', { name: 'Productos principales' })).toBeVisible();
    expect(await page.locator('.data-shell--fit').evaluateAll(items => items.every(item => item.scrollWidth <= item.clientWidth))).toBe(true);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  }
});

test('recorre productos, pedidos, liquidaciones y reportes del administrador', async ({ page }) => {
  let productoDesactivado = false;
  let productoCreado = false;
  let productoActualizado = false;
  let estadoActualizado = false;
  let pedidoCancelado = false;
  let liquidacionCreada = false;
  const producto = { id: 101, nombre: 'Auriculares Orbit', descripcion: 'Producto de prueba', categoria: 'Tecnología', precio: 39.90, img: '/products/images/auriculares-orbit.webp', galeria: [], stock: 10, activo: true };
  const pedidos = [
    { id: 501, numero: 'PRUEBA-ADMIN-001', total: 75.20, estado: 'confirmado', fecha: new Date().toISOString(), cliente: { nombre: 'Administrador', email: 'admin@upse.edu.ec', direccion: 'UPSE', ciudad: 'Santa Elena', telefono: '0990000001' }, usuarioEmail: 'admin@upse.edu.ec', items: [] },
    { id: 502, numero: 'PRUEBA-ADMIN-002', total: 68.40, estado: 'enviado', fecha: new Date().toISOString(), cliente: { nombre: 'Administrador', email: 'admin@upse.edu.ec', direccion: 'UPSE', ciudad: 'Santa Elena', telefono: '0990000001' }, usuarioEmail: 'admin@upse.edu.ec', items: [] },
  ];
  const comisiones = [
    { id: 31, pedidoId: 501, numeroPedido: 'DEMO-001', referidoId: 6, referido: 'Cliente Demo', nivel: 1, porcentaje: 0.05, monto: 3.25, estado: 'pendiente', fecha: new Date().toISOString() },
    { id: 32, pedidoId: 502, numeroPedido: 'DEMO-002', referidoId: 6, referido: 'Cliente Demo', nivel: 1, porcentaje: 0.05, monto: 4.10, estado: 'pendiente', fecha: new Date().toISOString() },
    { id: 33, pedidoId: 503, numeroPedido: 'DEMO-003', referidoId: 2, referido: 'Ana Torres', nivel: 2, porcentaje: 0.03, monto: 2.20, estado: 'pendiente', fecha: new Date().toISOString() },
  ];
  await page.route('**/api/**', route => {
    const url = new URL(route.request().url());
    const method = route.request().method();
    if (url.pathname === '/api/admin/productos' && method === 'GET') return route.fulfill({ json: [producto] });
    if (url.pathname === '/api/admin/productos' && method === 'POST') { productoCreado = true; return route.fulfill({ status: 201, json: { ...producto, id: 121 } }); }
    if (url.pathname === '/api/admin/productos/101' && method === 'PUT') { productoActualizado = route.request().postDataJSON().img.startsWith('/products/'); return route.fulfill({ json: producto }); }
    if (url.pathname === '/api/admin/productos/101' && method === 'DELETE') {
      productoDesactivado = true;
      return route.fulfill({ status: 200, json: null });
    }
    if (url.pathname === '/api/admin/pedidos' && method === 'GET') return route.fulfill({ json: { datos: pedidos, pagina: 1, limite: 10, total: pedidos.length } });
    if (url.pathname === '/api/admin/pedidos/501/estado' && method === 'PATCH') { estadoActualizado = true; return route.fulfill({ status: 204 }); }
    if (url.pathname === '/api/admin/pedidos/501/cancelar' && method === 'POST') { pedidoCancelado = true; return route.fulfill({ status: 204 }); }
    if (url.pathname === '/api/admin/comisiones' && method === 'GET') return route.fulfill({ json: comisiones });
    if (url.pathname === '/api/admin/liquidaciones' && method === 'POST') { liquidacionCreada = true; return route.fulfill({ status: 201, json: { id: 77, monto: 7.35, fecha: new Date().toISOString() } }); }
    if (url.pathname === '/api/admin/reportes/general') return route.fulfill({ json: { totalVentas: 143.60, totalPedidos: 2, ticketPromedio: 71.80, totalComisiones: 9.55, pendientePagar: 9.55, productosPrincipales: [{ nombre: 'Auriculares Orbit', cantidad: 2, ventas: 79.80 }], referidosPrincipales: [{ nombre: 'Cliente Demo', ventas: 208.80, comisiones: 8.36 }] } });
    return route.fulfill({ status: 404, json: { error: `Ruta simulada no encontrada: ${method} ${url.pathname}` } });
  });
  await page.goto('/login');
  await page.evaluate(() => {
    const payload = btoa(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 }));
    localStorage.setItem('multicatalogo_sesion', JSON.stringify({ email: 'admin@upse.edu.ec', token: `a.${payload}.c`, rol: 'admin' }));
  });
  await page.goto('/admin/productos');

  await page.getByLabel('Nombre').fill('Producto nuevo');
  await page.getByLabel('Precio').fill('20');
  await page.getByLabel('Categoría').fill('Prueba');
  await page.getByLabel('Stock').fill('5');
  await page.getByRole('button', { name: 'Crear producto' }).click();
  await expect.poll(() => productoCreado).toBe(true);
  await expect(page.getByText('Producto creado correctamente.')).toBeVisible();

  await page.getByRole('button', { name: 'Editar' }).click();
  await expect(page.getByLabel('Ruta o URL de imagen principal')).toHaveValue('/products/images/auriculares-orbit.webp');
  await page.getByRole('button', { name: 'Guardar cambios' }).click();
  await expect.poll(() => productoActualizado).toBe(true);
  await expect(page.getByText('Producto actualizado correctamente.')).toBeVisible();

  await page.getByRole('button', { name: 'Desactivar' }).click();
  const modal = page.getByRole('alertdialog');
  await expect(modal).toBeVisible();
  await expect(modal.getByRole('heading', { name: 'Desactivar Auriculares Orbit' })).toBeVisible();
  await expect(modal.getByText('dejará de aparecer en el catálogo público', { exact: false })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(modal).toBeHidden();
  expect(productoDesactivado).toBe(false);

  await page.getByRole('button', { name: 'Desactivar' }).click();
  await modal.getByRole('button', { name: 'Sí, desactivar' }).click();
  await expect.poll(() => productoDesactivado).toBe(true);
  await expect(page.getByText('Producto desactivado correctamente.')).toBeVisible();

  await page.goto('/admin/pedidos');
  await expect(page.getByRole('button', { name: 'Cancelar' })).toHaveCount(1);
  await page.getByLabel('Estado de PRUEBA-ADMIN-001').selectOption('preparando');
  await expect.poll(() => estadoActualizado).toBe(true);
  await page.getByRole('button', { name: 'Cancelar' }).click();
  await page.getByRole('alertdialog').getByRole('button', { name: 'Sí, cancelar pedido' }).click();
  await expect.poll(() => pedidoCancelado).toBe(true);

  await page.goto('/admin/comisiones');
  await page.getByLabel('Seleccionar comisión 31').check();
  await expect(page.getByLabel('Seleccionar comisión 32')).toBeEnabled();
  await expect(page.getByLabel('Seleccionar comisión 33')).toBeDisabled();
  await page.getByLabel('Seleccionar comisión 32').check();
  await page.getByRole('button', { name: 'Crear liquidación (2)' }).click();
  await expect(page.getByRole('alertdialog')).toContainText('$7.35');
  await page.getByRole('alertdialog').getByRole('button', { name: 'Confirmar liquidación' }).click();
  await expect.poll(() => liquidacionCreada).toBe(true);
  await expect(page.getByText('Liquidación #77 creada por $7.35.')).toBeVisible();

  await page.goto('/admin/reportes');
  await expect(page.getByRole('heading', { name: 'Reportes generales' })).toBeVisible();
  await expect(page.getByText('$143.60')).toBeVisible();
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
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
