import type { SVGProps } from "react";

// Icono "tienda" (bolsa de compras, SVG en línea). Usa currentColor: toma el color del texto donde se coloque.
const TiendaIcon = (props: SVGProps<SVGSVGElement>) => (
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
    <path d="M6 7h12l1 14H5L6 7Z" />
    <path d="M9 7V6a3 3 0 0 1 6 0v1" />
  </svg>
);

export default TiendaIcon;
