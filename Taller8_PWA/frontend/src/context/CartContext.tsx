// src/context/CartContext.tsx
import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { useAuth } from './AuthContext';

// 1. Definimos las interfaces (tipos de datos)
export interface Producto {
  id: number;
  nombre: string;
  descripcion?: string;
  categoria?: string;
  precio: number;
  img: string;
  galeria?: string[];
  stock: number;
  activo?: boolean;
}

export interface CartItem extends Producto {
  cantidad: number;
}

interface CartContextType {
  cart: CartItem[];
  addToCart: (producto: Producto) => void;
  removeFromCart: (id: number) => void;
  clearCart: () => void;
  totalItems: number;
  totalPrice: number;
}

// 2. Creamos el contexto indicando que puede ser CartContextType o undefined
const CartContext = createContext<CartContextType | undefined>(undefined);

// 3. Hook personalizado con validación de tipo
export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error("useCart debe ser usado dentro de un CartProvider");
  }
  return context;
};

// 4. Tipamos los props del Provider
interface CartProviderProps {
  children: ReactNode;
}

// 5. El Provider
export const CartProvider = ({ children }: CartProviderProps) => {
  // Le decimos a useState que este arreglo contendrá objetos de tipo CartItem
  const { userEmail } = useAuth();
  const storageKey = userEmail ? 'multicatalogo_carrito_' + userEmail : null;
  const [cart, setCart] = useState<CartItem[]>(() => {
    try {
      const saved: unknown = storageKey ? JSON.parse(localStorage.getItem(storageKey) ?? '[]') : [];
      return Array.isArray(saved) ? saved.filter((item): item is CartItem =>
        item && Number.isInteger(item.id) && item.id > 0 && typeof item.nombre === 'string' &&
        typeof item.precio === 'number' && Number.isFinite(item.precio) && item.precio >= 0 &&
        Number.isInteger(item.cantidad) && item.cantidad > 0) : [];
    } catch { return []; }
  });
  useEffect(() => {
    if (storageKey) localStorage.setItem(storageKey, JSON.stringify(cart));
  }, [cart, storageKey]);

  const addToCart = (producto: Producto) => {
    setCart((prevCart) => {
      const itemExists = prevCart.find((item) => item.id === producto.id);
      if (itemExists) {
        if (itemExists.cantidad >= producto.stock) return prevCart;
        return prevCart.map((item) =>
          item.id === producto.id
            ? { ...item, cantidad: item.cantidad + 1 }
            : item
        );
      }
      return producto.stock > 0 ? [...prevCart, { ...producto, cantidad: 1 }] : prevCart;
    });
  };

  const removeFromCart = (id: number) => {
    setCart((prevCart) => prevCart.filter((item) => item.id !== id));
  };

  const clearCart = () => setCart([]);

  const totalItems = cart.reduce((total, item) => total + item.cantidad, 0);
  const totalPrice = cart.reduce((total, item) => total + item.precio * item.cantidad, 0);

  return (
    <CartContext.Provider value={{ cart, addToCart, removeFromCart, clearCart, totalItems, totalPrice }}>
      {children}
    </CartContext.Provider>
  );
};
