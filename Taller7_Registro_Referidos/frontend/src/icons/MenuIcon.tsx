import type { SVGProps } from "react";

// Icono "menú (hamburguesa)" (SVG en línea). Usa currentColor: toma el color del texto donde se coloque.
const MenuIcon = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth={2}
    strokeLinecap="round"
    strokeLinejoin="round"
    width={24}
    height={24}
    aria-hidden="true"
    {...props}
  >
    <path d="M4 6h16" />
    <path d="M4 12h16" />
    <path d="M4 18h16" />
  </svg>
);

export default MenuIcon;
