# Taller 3 - Sidebar colapsable, iconos y diseño responsivo

Extensión del proyecto base *MultiCatálogo* (React + TypeScript + Vite + Tailwind).

## ¿Qué se hizo?

1. **Prop `isCollapsed` en el `Sidebar`**: si es `true` el menú mide **80px** (`w-20`) y muestra solo los iconos; si es `false` mide 256px (`w-64`) y muestra icono + texto.
2. **Iconos en las opciones del menú**: carpeta `src/icons/` con un componente SVG por icono (Dashboard, Catálogo, Mi Red y Menú). Usan `currentColor`, así que toman el color del texto.
3. **Botón Toggle en el `Navbar`** que controla el estado **global** (`SidebarContext`, mismo patrón que `AuthContext` y `CartContext`).
4. **Responsive para celulares**: en pantallas menores a 768px el sidebar es un menú deslizable con fondo oscuro; el Navbar, el Layout, el Carrito y la tabla de Mi Red se adaptan al ancho.

## ¿Cómo funciona?

```
SidebarProvider (src/context/SidebarContext.tsx)  -> guarda isCollapsed y toggleSidebar()
        |
        +--> Navbar  : el botón ☰ llama a toggleSidebar()
        +--> Layout  : lee isCollapsed y lo pasa como prop:  <Sidebar isCollapsed={isCollapsed} />
                          +--> Sidebar : cambia su ancho / muestra u oculta los textos
```

## Ejecutar

```bash
npm install
npm run dev
```

Usuario de prueba: `admin@upse.edu.ec` / `123456`
