# Pokédex con React y PokeAPI

Este es un proyecto sencillo hecho con React que muestra los primeros 151 Pokémon.
Los datos salen de una API gratis llamada PokeAPI, así que no hay que guardar nada a
mano: la página los pide sola cada vez que entras.

## ¿Qué tiene la página?

- Un menú a la izquierda para moverte entre **Inicio** y **Pokémons**.
- Una lista de Pokémon en tarjetas, cada una con su número, su foto, su nombre y sus tipos.
- Filtros en la sección de Pokémon:
  - Por **tipo** (fuego, agua, planta, etc.).
  - Por **letra inicial** (si eliges la A, solo aparecen los que empiezan por A).
  - Por **orden del nombre** (de la A a la Z o al revés).
- Se ve bien tanto en computadora como en celular (en el celular el menú se abre con el botón ☰).

## ¿Con qué está hecho?

- **React** con **TypeScript** (la base de la página).
- **Vite** (la herramienta que la levanta y la arma).
- **Tailwind CSS** (para los estilos y colores).
- **React Router** (para cambiar entre páginas sin recargar).

## ¿Cómo hacerlo funcionar?

Necesitas tener instalado **Node.js**. Después, abre una terminal en la carpeta del
proyecto y escribe:

1. Instalar todo lo que necesita el proyecto:

   ```bash
   npm install
   ```

2. Levantar la página en tu computadora:

   ```bash
   npm run dev
   ```

3. Abre en el navegador la dirección que aparece en la terminal (normalmente
   `http://localhost:5173`).

## ¿Cómo está organizado?

```
src/
├── components/            (piezas que se repiten)
│   ├── Layout.tsx         -> el marco: menú + barra de arriba
│   ├── Sidebar.tsx        -> el menú de la izquierda
│   ├── Navbar.tsx         -> la barra de arriba
│   └── TarjetaPokemon.tsx -> la tarjeta de cada Pokémon
├── pages/                 (las páginas)
│   ├── Inicio.tsx         -> la página de bienvenida
│   └── Pokemons.tsx       -> la lista con los filtros
├── App.tsx                -> decide qué página mostrar
└── main.tsx               -> el punto donde arranca todo
```

## ¿De dónde salen los datos?

De la PokeAPI: https://pokeapi.co
La dirección que usamos es: https://pokeapi.co/api/v2/pokemon?limit=151
