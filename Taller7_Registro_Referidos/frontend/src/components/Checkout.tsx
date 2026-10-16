import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';
import { checkoutRequest } from '../services/api';

interface DatosEnvio {
  nombre: string;
  email: string;
  direccion: string;
  ciudad: string;
  telefono: string;
  metodoPago: string;
}

const METODOS_PAGO = ['Tarjeta de crédito', 'Transferencia bancaria', 'Efectivo contra entrega'];

const Checkout = () => {
  const { cart, totalPrice, clearCart } = useCart();
  const { userEmail } = useAuth();
  const navigate = useNavigate();
  const [datos, setDatos] = useState<DatosEnvio>({
    nombre: '',
    email: userEmail ?? '',
    direccion: '',
    ciudad: '',
    telefono: '',
    metodoPago: METODOS_PAGO[0],
  });
  const [procesando, setProcesando] = useState(false);
  const [error, setError] = useState('');

  if (cart.length === 0) {
    return (
      <div className="bg-white p-10 rounded-lg border border-slate-200 text-center">
        <p className="text-slate-600 text-lg mb-4">No tienes productos en el carrito.</p>
        <Link to="/catalogo" className="text-indigo-600 font-semibold hover:underline">
          Ir al catálogo
        </Link>
      </div>
    );
  }

  const actualizarDato = (campo: keyof DatosEnvio, valor: string) => {
    setDatos((actual) => ({ ...actual, [campo]: valor }));
  };

  const confirmarPedido = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setProcesando(true);
    setError('');

    try {
      const pedido = await checkoutRequest({
        usuarioEmail: userEmail ?? datos.email,
        cliente: {
          nombre: datos.nombre,
          email: datos.email,
          direccion: datos.direccion,
          ciudad: datos.ciudad,
          telefono: datos.telefono,
        },
        metodoPago: datos.metodoPago,
        items: cart.map((item) => ({ productoId: item.id, cantidad: item.cantidad })),
      });
      clearCart();
      navigate('/confirmacion', { state: { pedido } });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudo registrar el pedido.');
    } finally {
      setProcesando(false);
    }
  };

  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-800 mb-6">Finalizar Compra</h1>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <form onSubmit={confirmarPedido} className="lg:col-span-2 space-y-5">
          <div className="bg-white p-6 rounded-lg border border-slate-200 shadow-sm">
            <h2 className="text-lg font-bold text-slate-800 mb-4">Datos de Envío</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {([
                ['nombre', 'Nombre completo', 'text'],
                ['email', 'Correo electrónico', 'email'],
                ['direccion', 'Dirección', 'text'],
                ['ciudad', 'Ciudad', 'text'],
                ['telefono', 'Teléfono', 'tel'],
              ] as const).map(([campo, etiqueta, tipo]) => (
                <label key={campo} className={campo === 'direccion' ? 'md:col-span-2' : ''}>
                  <span className="block text-sm font-medium text-slate-700 mb-1">{etiqueta} *</span>
                  <input
                    type={tipo}
                    required
                    value={datos[campo]}
                    onChange={(event) => actualizarDato(campo, event.target.value)}
                    className="w-full px-4 py-2.5 rounded-lg border border-slate-300 focus:ring-2 focus:ring-indigo-600 outline-none"
                  />
                </label>
              ))}
            </div>
          </div>

          <div className="bg-white p-6 rounded-lg border border-slate-200 shadow-sm">
            <h2 className="text-lg font-bold text-slate-800 mb-4">Método de Pago</h2>
            <div className="space-y-3">
              {METODOS_PAGO.map((metodo) => (
                <label key={metodo} className="flex items-center gap-3 p-3 rounded-lg border border-slate-200 cursor-pointer">
                  <input
                    type="radio"
                    name="metodoPago"
                    value={metodo}
                    checked={datos.metodoPago === metodo}
                    onChange={(event) => actualizarDato('metodoPago', event.target.value)}
                    className="accent-indigo-600"
                  />
                  <span className="text-sm text-slate-700">{metodo}</span>
                </label>
              ))}
            </div>
          </div>

          {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg border border-red-200">{error}</div>}
          <button
            type="submit"
            disabled={procesando}
            className="w-full bg-indigo-600 text-white py-3 rounded-lg font-semibold hover:bg-indigo-700 disabled:opacity-50"
          >
            {procesando ? 'Registrando pedido...' : 'Confirmar Pedido'}
          </button>
        </form>

        <div className="bg-white p-6 rounded-lg border border-slate-200 shadow-sm h-fit">
          <h2 className="text-lg font-bold text-slate-800 mb-4">Resumen del Pedido</h2>
          <div className="space-y-3 border-b border-slate-100 pb-4 mb-4">
            {cart.map((item) => (
              <div key={item.id} className="flex justify-between text-sm gap-4">
                <span className="text-slate-600">{item.nombre} × {item.cantidad}</span>
                <span className="font-semibold">${(item.precio * item.cantidad).toFixed(2)}</span>
              </div>
            ))}
          </div>
          <div className="flex justify-between items-center">
            <span className="font-bold">Total estimado</span>
            <span className="text-2xl font-bold text-indigo-600">${totalPrice.toFixed(2)}</span>
          </div>
          <p className="text-xs text-slate-500 mt-3">El servidor vuelve a calcular el total con los precios guardados.</p>
        </div>
      </div>
    </div>
  );
};

export default Checkout;
