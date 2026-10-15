import type { SVGProps } from "react";

// Icono "catálogo" (SVG en línea). Usa currentColor: toma el color del texto donde se coloque.
const CatalogoIcon = (props: SVGProps<SVGSVGElement>) => (
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
    <path d="M21 8 12 3 3 8v8l9 5 9-5V8Z" />
    <path d="M3 8l9 5 9-5" />
    <path d="M12 13v8" />
  </svg>
);

export default CatalogoIcon;
