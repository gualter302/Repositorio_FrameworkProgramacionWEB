// src/App.tsx
import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './components/Dashboard';
import Catalogo from './components/Catalogo';
import MiRed from './components/MiRed';
import Carrito from './components/Carrito';
import Login from './components/Login';
import Checkout from './components/Checkout';
import Confirmacion from './components/Confirmacion';
import Registro from './components/Registro';
import AdminProductos from './components/AdminProductos';
import Pedidos from './components/Pedidos';
import AdminPedidos from './components/AdminPedidos';
import Comisiones from './components/Comisiones';
import AdminComisiones from './components/AdminComisiones';
import AdminReportes from './components/AdminReportes';
import PwaControls from './components/PwaControls';
import { CartProvider } from './context/CartContext';
import { AuthProvider, useAuth } from './context/AuthContext';
import { SidebarProvider } from './context/SidebarContext';

// Componente para proteger las rutas privadas
const ProtectedRoute = () => {
  const { isAuthenticated } = useAuth();

  // Si no está autenticado, lo enviamos al login
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Si está autenticado, renderiza las rutas hijas (Outlet)
  return <Outlet />;
};

const AdminRoute = () => {
  const { role } = useAuth();
  return role === 'admin' ? <Outlet /> : <Navigate to="/" replace />;
};

function App() {
  return (
    <AuthProvider> {/* Proveedor de Autenticación */}
      <UserCart>
        <SidebarProvider> {/* Proveedor del estado global del Sidebar */}
          <BrowserRouter>
            <PwaControls />
            <Routes>
              {/* Ruta pública */}
              <Route path="/login" element={<Login />} />
              <Route path="/registro" element={<Registro />} />
              <Route path="/catalogo" element={<Layout />}>
                <Route index element={<Catalogo />} />
              </Route>

              {/* Rutas protegidas */}
              <Route element={<ProtectedRoute />}>
                <Route path="/" element={<Layout />}>
                  <Route index element={<Dashboard />} />
                  <Route path="mi-red" element={<MiRed />} />
                  <Route path="carrito" element={<Carrito />} />
                  <Route path="checkout" element={<Checkout />} />
                  <Route path="confirmacion" element={<Confirmacion />} />
                    <Route path="pedidos" element={<Pedidos />} />
                    <Route path="comisiones" element={<Comisiones />} />
                  <Route element={<AdminRoute />}>
                    <Route path="admin/productos" element={<AdminProductos />} />
                    <Route path="admin/pedidos" element={<AdminPedidos />} />
                    <Route path="admin/comisiones" element={<AdminComisiones />} />
                    <Route path="admin/reportes" element={<AdminReportes />} />
                  </Route>
                </Route>
              </Route>

              {/* Ruta comodín para capturar 404 y redirigir */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </BrowserRouter>
        </SidebarProvider>
      </UserCart>
    </AuthProvider>
  );
}

export default App;

function UserCart({ children }: { children: React.ReactNode }) {
  const { userEmail } = useAuth();
  return <CartProvider key={userEmail ?? 'guest'}>{children}</CartProvider>;
}
