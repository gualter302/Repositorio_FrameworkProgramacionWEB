import { useEffect, useState, type FormEvent } from 'react';
import {
  actualizarProducto,
  crearProducto,
  desactivarProducto,
  getProductosAdmin,
} from '../services/api';
import type { Producto } from '../context/CartContext';
import { Edit3, PackagePlus, Power } from 'lucide-react';
import { ConfirmDialog, EmptyState, PageHeader, StatusBadge } from './Ui';

type DatosProducto = Omit<Producto, 'id' | 'activo'>;

const inicial: DatosProducto = {
  nombre: '',
  descripcion: '',
  categoria: '',
  precio: 0,
  img: '',
  galeria: [],
  stock: 0,
};

const AdminProductos = () => {
  const [productos, setProductos] = useState<Producto[]>([]);
  const [formulario, setFormulario] = useState<DatosProducto>(inicial);
  const [editando, setEditando] = useState<number | null>(null);
  const [mensaje, setMensaje] = useState('');
  const [error, setError] = useState('');
  const [cargando, setCargando] = useState(true);
  const [productoADesactivar, setProductoADesactivar] = useState<Producto | null>(null);
  const [desactivando, setDesactivando] = useState(false);

  const cargar = async () => {
    setError('');
    try {
      setProductos(await getProductosAdmin());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudo cargar el catálogo');
    } finally {
      setCargando(false);
    }
  };

  useEffect(() => {
    let activo = true;
    getProductosAdmin()
      .then((datos) => {
        if (activo) setProductos(datos);
      })
      .catch((err: Error) => {
        if (activo) setError(err.message);
      })
      .finally(() => {
        if (activo) setCargando(false);
      });
    return () => {
      activo = false;
    };
  }, []);

  const cambiar = (campo: keyof DatosProducto, valor: string | number) => {
    setFormulario((actual) => ({ ...actual, [campo]: valor }));
  };

  const guardar = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setMensaje('');
    setError('');
    try {
      if (editando === null) {
        await crearProducto(formulario);
        setMensaje('Producto creado correctamente.');
      } else {
        await actualizarProducto(editando, formulario);
        setMensaje('Producto actualizado correctamente.');
      }
      setFormulario(inicial);
      setEditando(null);
      await cargar();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudo guardar el producto');
    }
  };

  const editar = (producto: Producto) => {
    setEditando(producto.id);
    setFormulario({
      nombre: producto.nombre,
      descripcion: producto.descripcion ?? '',
      categoria: producto.categoria ?? '',
      precio: producto.precio,
      img: producto.img,
      galeria: producto.galeria ?? [],
      stock: producto.stock,
    });
    setMensaje('');
  };

  const desactivar = async () => {
    if (!productoADesactivar) return;
    setError('');
    setDesactivando(true);
    try {
      await desactivarProducto(productoADesactivar.id);
      setMensaje('Producto desactivado correctamente.');
      await cargar();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudo desactivar el producto');
    } finally {
      setDesactivando(false);
      setProductoADesactivar(null);
    }
  };

  return (
    <div className="page page-stack">
      <PageHeader kicker="Administración" title="Productos" description="Crea, edita o retira artículos del catálogo público." />

      {mensaje && <div className="alert-success">{mensaje}</div>}
      {error && <div className="alert-error" role="alert">{error}</div>}

      <form onSubmit={guardar} className="glass-panel grid gap-4 md:grid-cols-2">
        <h2 className="section-title md:col-span-2">
          {editando === null ? 'Nuevo producto' : `Editando producto #${editando}`}
        </h2>
        <label><span className="field-label">Nombre</span><input required value={formulario.nombre} onChange={(e) => cambiar('nombre', e.target.value)} className="form-input" /></label>
        <label><span className="field-label">Precio</span><input required type="number" min="0" step="0.01" value={formulario.precio} onChange={(e) => cambiar('precio', Number(e.target.value))} className="form-input" /></label>
        <label><span className="field-label">Categoría</span><input value={formulario.categoria} onChange={(e) => cambiar('categoria', e.target.value)} className="form-input" /></label>
        <label><span className="field-label">Stock</span><input required type="number" min="0" value={formulario.stock} onChange={(e) => cambiar('stock', Number(e.target.value))} className="form-input" /></label>
        <label className="md:col-span-2"><span className="field-label">URL de imagen principal</span><input type="url" value={formulario.img} onChange={(e) => cambiar('img', e.target.value)} placeholder="https://..." className="form-input" /></label>
        <label className="md:col-span-2"><span className="field-label">Descripción</span><textarea value={formulario.descripcion} onChange={(e) => cambiar('descripcion', e.target.value)} className="form-input" rows={3} /></label>
        <div className="flex flex-wrap gap-3 md:col-span-2">
          <button type="submit" className="btn-primary"><PackagePlus size={17} />
            {editando === null ? 'Crear producto' : 'Guardar cambios'}
          </button>
          {editando !== null && (
            <button type="button" onClick={() => { setEditando(null); setFormulario(inicial); }} className="btn-secondary">
              Cancelar
            </button>
          )}
        </div>
      </form>

      {cargando ? <div className="glass-panel animate-pulse text-slate-500">Cargando productos...</div> : productos.length === 0 ? <EmptyState icon={PackagePlus} title="No hay productos" description="Usa el formulario para registrar el primer artículo del catálogo." /> : <div className="data-shell">
          <table className="data-table">
            <thead>
              <tr><th>Producto</th><th>Categoría</th><th>Precio</th><th>Stock</th><th>Estado</th><th>Acciones</th></tr>
            </thead>
            <tbody>
              {productos.map((producto) => (
                <tr key={producto.id}>
                  <td className="font-medium text-white">{producto.nombre}</td>
                  <td>{producto.categoria || 'Sin categoría'}</td>
                  <td className="mono text-cyan-300">${producto.precio.toFixed(2)}</td>
                  <td>{producto.stock}</td>
                  <td><StatusBadge>{producto.activo ? 'Activo' : 'Inactivo'}</StatusBadge></td>
                  <td><div className="flex gap-2">
                    <button onClick={() => editar(producto)} className="btn-secondary min-h-9 px-3 py-1.5 text-xs"><Edit3 size={14} />Editar</button>
                    {producto.activo && <button onClick={() => setProductoADesactivar(producto)} className="btn-danger min-h-9 px-3 py-1.5 text-xs"><Power size={14} />Desactivar</button>}
                  </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
      </div>}
      <ConfirmDialog
        open={productoADesactivar !== null}
        title={`Desactivar ${productoADesactivar?.nombre ?? 'producto'}`}
        description="El producto dejará de aparecer en el catálogo público, pero conservará su información y el historial de pedidos."
        confirmLabel="Sí, desactivar"
        onConfirm={() => void desactivar()}
        onCancel={() => setProductoADesactivar(null)}
        busy={desactivando}
      />
    </div>
  );
};

export default AdminProductos;
