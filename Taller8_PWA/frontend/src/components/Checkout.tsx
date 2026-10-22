import { useRef, useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';
import { ApiError, checkoutRequest, getProductos, type CheckoutRequest } from '../services/api';
import { ArrowLeft, CheckCircle2, CreditCard, MapPin, ShoppingBag } from 'lucide-react';
import { EmptyState, PageHeader } from './Ui';

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
  const busy = useRef(false);

  if (cart.length === 0) {
    return (
      <div className="page"><EmptyState icon={ShoppingBag} title="No hay productos para confirmar" description="Añade al menos un producto antes de completar los datos de envío." action={<Link to="/catalogo" className="btn-primary"><ArrowLeft size={17} />Ir al catálogo</Link>} /></div>
    );
  }

  const actualizarDato = (campo: keyof DatosEnvio, valor: string) => {
    setDatos((actual) => ({ ...actual, [campo]: valor }));
  };

  const confirmarPedido = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (busy.current) return;
    if (!navigator.onLine) { setError('Conéctate a internet para confirmar la compra. Tu carrito sigue guardado.'); return; }
    busy.current = true;
    setProcesando(true);
    setError('');

    try {
      const pendingKey = 'multicatalogo_compra_' + userEmail;
      const pending = localStorage.getItem(pendingKey);
      const purchase: CheckoutRequest = {
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
      };
      let attempt: { key: string; body: CheckoutRequest };
      if (pending) {
        attempt = JSON.parse(pending);
      } else {
        const current = await getProductos(true);
        for (const item of cart) {
          const product = current.find((product) => product.id === item.id);
          if (!product || product.stock < item.cantidad) throw new Error('Algún producto ya no tiene existencias suficientes. Revisa el catálogo antes de comprar.');
          if (product.precio !== item.precio) throw new Error('Un precio cambió. Retira ese producto y vuelve a añadirlo desde el catálogo actualizado.');
        }
        attempt = { key: crypto.randomUUID(), body: purchase };
        localStorage.setItem(pendingKey, JSON.stringify(attempt));
      }
      const pedido = await checkoutRequest(attempt.body, attempt.key);
      localStorage.removeItem(pendingKey);
      clearCart();
      navigate('/confirmacion', { state: { pedido } });
    } catch (err) {
      if (err instanceof ApiError && [400, 409].includes(err.status))
        localStorage.removeItem('multicatalogo_compra_' + userEmail);
      setError((err instanceof Error ? err.message : 'No se pudo registrar el pedido.') + ' Si el envío quedó pendiente, al confirmar otra vez se recuperará la misma compra original. Revisa también Mis pedidos.');
    } finally {
      setProcesando(false);
      busy.current = false;
    }
  };

  return (
    <div className="page page-stack">
      <PageHeader kicker="Paso final" title="Finalizar compra" description="Completa los datos de entrega. Antes de registrar el pedido volveremos a comprobar precio y existencias." />
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <form onSubmit={confirmarPedido} className="lg:col-span-2 space-y-5">
          <section className="glass-panel">
            <div className="mb-5 flex items-center gap-3"><span className="empty-icon h-10 w-10 rounded-xl"><MapPin size={20} /></span><div><p className="page-kicker mb-1">Entrega</p><h2 className="section-title">Datos de envío</h2></div></div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {([
                ['nombre', 'Nombre completo', 'text'],
                ['email', 'Correo electrónico', 'email'],
                ['direccion', 'Dirección', 'text'],
                ['ciudad', 'Ciudad', 'text'],
                ['telefono', 'Teléfono', 'tel'],
              ] as const).map(([campo, etiqueta, tipo]) => (
                <label key={campo} className={campo === 'direccion' ? 'md:col-span-2' : ''}>
                  <span className="field-label">{etiqueta} *</span>
                  <input
                    type={tipo}
                    required
                    value={datos[campo]}
                    onChange={(event) => actualizarDato(campo, event.target.value)}
                    className="form-input"
                  />
                </label>
              ))}
            </div>
          </section>

          <section className="glass-panel">
            <div className="mb-5 flex items-center gap-3"><span className="empty-icon h-10 w-10 rounded-xl"><CreditCard size={20} /></span><div><p className="page-kicker mb-1">Preferencia</p><h2 className="section-title">Método de pago</h2></div></div>
            <div className="space-y-3">
              {METODOS_PAGO.map((metodo) => (
                <label key={metodo} className={`flex cursor-pointer items-center gap-3 rounded-2xl border p-3.5 transition ${datos.metodoPago === metodo ? 'border-cyan-300/45 bg-cyan-300/10' : 'border-cyan-200/10 bg-white/[0.025] hover:border-cyan-200/25'}`}>
                  <input
                    type="radio"
                    name="metodoPago"
                    value={metodo}
                    checked={datos.metodoPago === metodo}
                    onChange={(event) => actualizarDato('metodoPago', event.target.value)}
                    className="h-4 w-4 accent-cyan-400"
                  />
                  <span className="text-sm text-slate-200">{metodo}</span>
                </label>
              ))}
            </div>
            <p className="mt-4 text-xs leading-relaxed text-slate-500">Esta selección se guarda con el pedido. La aplicación todavía no procesa cobros con tarjetas ni transferencias bancarias.</p>
          </section>

          {error && <div className="alert-error" role="alert">{error}</div>}
          <button
            type="submit"
            disabled={procesando}
            className="btn-primary w-full"
          >
            <CheckCircle2 size={18} />
            {procesando ? 'Registrando pedido...' : 'Confirmar Pedido'}
          </button>
        </form>

        <aside className="glass-panel glass-panel--accent h-fit lg:sticky lg:top-4">
          <p className="page-kicker">Pedido</p><h2 className="section-title mb-5">Resumen del Pedido</h2>
          <div className="mb-4 space-y-3 border-b border-cyan-200/10 pb-4">
            {cart.map((item) => (
              <div key={item.id} className="flex justify-between text-sm gap-4">
                <span className="text-slate-400">{item.nombre} × {item.cantidad}</span>
                <span className="mono text-slate-200">${(item.precio * item.cantidad).toFixed(2)}</span>
              </div>
            ))}
          </div>
          <div className="flex justify-between items-center">
            <span className="font-bold text-white">Total estimado</span>
            <span className="mono text-2xl text-cyan-300">${totalPrice.toFixed(2)}</span>
          </div>
          <p className="text-xs text-slate-500 mt-3">El servidor vuelve a calcular el total con los precios guardados.</p>
        </aside>
      </div>
    </div>
  );
};

export default Checkout;
