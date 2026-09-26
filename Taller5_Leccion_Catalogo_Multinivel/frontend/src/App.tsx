// src/App.tsx
import { type ReactNode } from 'react';
import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './components/Dashboard';
import Catalogo from './components/Catalogo';
import Storefront from './components/Storefront';
import DetalleProducto from './components/DetalleProducto';
import MiRed from './components/MiRed';
import Carrito from './components/Carrito';
import Login from './components/Login';
import { CartProvider } from './context/CartContext';
import { AuthProvider, useAuth } from './context/AuthContext';
import { SidebarProvider } from './context/SidebarContext';
import { INICIO_CLIENTE } from './utils/rutas';

// Componente para proteger las rutas privadas: exige estar autenticado (cualquier rol)
const ProtectedRoute = () => {
  const { isAuthenticated } = useAuth();

  // Si no está autenticado, lo enviamos al login
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Si está autenticado, renderiza las rutas hijas (Outlet)
  return <Outlet />;
};

// Componente para rutas exclusivas del administrador (Tema 5: roles)
const AdminRoute = () => {
  const { isAuthenticated, user } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Si el usuario no es admin, lo enviamos a su vista por defecto
  if (user?.rol !== 'admin') {
    return <Navigate to={INICIO_CLIENTE} replace />;
  }

  return <Outlet />;
};

// Proveedor intermedio del carrito: al cambiar de usuario se remonta con key,
// de modo que cada cuenta tenga SU propio carrito (se persiste en el Paso 10).
const CartBoundary = ({ children }: { children: ReactNode }) => {
  const { user } = useAuth();
  return <CartProvider key={user?.email ?? 'anonimo'}>{children}</CartProvider>;
};

function App() {
  return (
    <AuthProvider> {/* Proveedor de Autenticación */}
      <CartBoundary> {/* Proveedor del Carrito (uno por usuario) */}
        <SidebarProvider> {/* Proveedor del estado global del Sidebar */}
          <BrowserRouter>
            <Routes>
              {/* Ruta pública */}
              <Route path="/login" element={<Login />} />

              {/* Rutas protegidas (requieren login) */}
              <Route element={<ProtectedRoute />}>
                <Route path="/" element={<Layout />}>
                  {/* Solo administrador: Dashboard y Red multinivel */}
                  <Route element={<AdminRoute />}>
                    <Route index element={<Dashboard />} />
                    <Route path="mi-red" element={<MiRed />} />
                  </Route>

                  {/* Ambos roles: flujo de compra en línea */}
                  <Route path="tienda" element={<Storefront />} />
                  <Route path="catalogo" element={<Catalogo />} />
                  <Route path="producto/:id" element={<DetalleProducto />} />
                  <Route path="carrito" element={<Carrito />} />
                </Route>
              </Route>

              {/* Ruta comodín para capturar 404 y redirigir */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </BrowserRouter>
        </SidebarProvider>
      </CartBoundary>
    </AuthProvider>
  );
}

export default App;
