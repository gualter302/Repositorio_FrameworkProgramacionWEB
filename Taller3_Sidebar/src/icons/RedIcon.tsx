import type { SVGProps } from "react";

// Icono "mi red" (SVG en línea). Usa currentColor: toma el color del texto donde se coloque.
const RedIcon = (props: SVGProps<SVGSVGElement>) => (
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
    <circle cx="9" cy="8" r="3.5" />
    <path d="M2.5 20a6.5 6.5 0 0 1 13 0" />
    <circle cx="17.5" cy="9" r="2.5" />
    <path d="M17 14.2A5.5 5.5 0 0 1 21.5 20" />
  </svg>
);

export default RedIcon;
