# Configuración opcional de Codex

Codex carga automáticamente las instrucciones de `AGENTS.md`. No necesitas activar nada adicional
para que se apliquen las reglas de este repositorio.

Playwright y memoria son servicios opcionales. Si los necesitas y tienes Codex CLI, puedes
registrarlos en tu configuración personal con estos comandos:

```bash
codex mcp add playwright -- npx -y @playwright/mcp
codex mcp add memory -- npx -y @modelcontextprotocol/server-memory
```

Comprueba los servicios registrados con:

```bash
codex mcp list
```

La extensión de Codex para VS Code y Codex CLI comparten la configuración personal. No se activa
MCP dentro del repositorio para evitar instalaciones automáticas en equipos que no lo necesitan.
También puedes copiar el contenido de `CODEX_CONFIG_EXAMPLE.toml` a tu archivo personal
`~/.codex/config.toml`.
