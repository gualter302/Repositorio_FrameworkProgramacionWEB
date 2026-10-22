// Servidor exclusivo de pruebas: simula otra versión sin modificar los archivos compilados.
import { createServer } from 'node:http';
import { readFile } from 'node:fs/promises';
import { resolve, extname, sep } from 'node:path';
const root = resolve('dist');
let version = 0;
const types = { '.html': 'text/html', '.js': 'application/javascript', '.css': 'text/css', '.svg': 'image/svg+xml', '.png': 'image/png', '.webmanifest': 'application/manifest+json' };
createServer(async (req, res) => {
  const url = new URL(req.url, 'http://localhost');
  if (url.pathname === '/__test/update' && req.method === 'POST') {
    version++; res.end('ok'); return;
  }
  if (url.pathname.startsWith('/api/')) {
    res.writeHead(503, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: 'API no disponible en este servidor de pruebas' })); return;
  }
  let file = resolve(root, '.' + decodeURIComponent(url.pathname));
  if (!file.startsWith(root + sep) && file !== root) { res.writeHead(403); res.end(); return; }
  try {
    let data;
    try { data = await readFile(file); }
    catch { file = resolve(root, 'index.html'); data = await readFile(file); }
    if (url.pathname === '/sw.js') data = Buffer.concat([data, Buffer.from('\n// versión de prueba ' + version)]);
    res.writeHead(200, { 'Content-Type': types[extname(file)] ?? 'application/octet-stream', 'Cache-Control': 'no-store' });
    res.end(data);
  } catch { res.writeHead(500); res.end(); }
}).listen(4178, 'localhost');
