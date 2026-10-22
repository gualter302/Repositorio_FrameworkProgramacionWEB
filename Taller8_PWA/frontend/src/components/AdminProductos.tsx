import { useEffect, useState, type FormEvent } from 'react';
import {
  actualizarProducto,
  crearProducto,
  desactivarProducto,
  getProductosAdmin,
} from '../services/api';
import type { Producto } from '../context/CartContext';

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

  const desactivar = async (producto: Producto) => {
    if (!window.confirm(`¿Desactivar "${producto.nombre}"?`)) return;
    setError('');
    try {
      await desactivarProducto(producto.id);
      setMensaje('Producto desactivado correctamente.');
      await cargar();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudo desactivar el producto');
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-800">Administrar productos</h1>
        <p className="text-slate-500 mt-1">Crea, edita o desactiva productos del catálogo.</p>
      </div>

      {mensaje && <div className="bg-green-50 text-green-700 p-3 rounded-lg">{mensaje}</div>}
      {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg">{error}</div>}

      <form onSubmit={guardar} className="bg-white border border-slate-200 rounded-lg p-5 grid gap-4 md:grid-cols-2">
        <h2 className="font-semibold text-slate-800 md:col-span-2">
          {editando === null ? 'Nuevo producto' : `Editando producto #${editando}`}
        </h2>
        <input required value={formulario.nombre} onChange={(e) => cambiar('nombre', e.target.value)} placeholder="Nombre" className="border rounded-lg p-3" />
        <input required type="number" min="0" step="0.01" value={formulario.precio} onChange={(e) => cambiar('precio', Number(e.target.value))} placeholder="Precio" className="border rounded-lg p-3" />
        <input value={formulario.categoria} onChange={(e) => cambiar('categoria', e.target.value)} placeholder="Categoría" className="border rounded-lg p-3" />
        <input required type="number" min="0" value={formulario.stock} onChange={(e) => cambiar('stock', Number(e.target.value))} placeholder="Stock" className="border rounded-lg p-3" />
        <input value={formulario.img} onChange={(e) => cambiar('img', e.target.value)} placeholder="URL de imagen principal" className="border rounded-lg p-3 md:col-span-2" />
        <textarea value={formulario.descripcion} onChange={(e) => cambiar('descripcion', e.target.value)} placeholder="Descripción" className="border rounded-lg p-3 md:col-span-2" rows={3} />
        <div className="flex gap-3 md:col-span-2">
          <button type="submit" className="bg-indigo-600 text-white px-5 py-2 rounded-lg hover:bg-indigo-700">
            {editando === null ? 'Crear producto' : 'Guardar cambios'}
          </button>
          {editando !== null && (
            <button type="button" onClick={() => { setEditando(null); setFormulario(inicial); }} className="border px-5 py-2 rounded-lg">
              Cancelar
            </button>
          )}
        </div>
      </form>

      <div className="bg-white border border-slate-200 rounded-lg overflow-x-auto">
        {cargando ? <p className="p-5 text-slate-500">Cargando productos...</p> : (
          <table className="w-full text-sm">
            <thead className="bg-slate-100 text-left">
              <tr><th className="p-3">Producto</th><th className="p-3">Categoría</th><th className="p-3">Precio</th><th className="p-3">Stock</th><th className="p-3">Estado</th><th className="p-3">Acciones</th></tr>
            </thead>
            <tbody>
              {productos.map((producto) => (
                <tr key={producto.id} className="border-t">
                  <td className="p-3 font-medium">{producto.nombre}</td>
                  <td className="p-3">{producto.categoria || 'Sin categoría'}</td>
                  <td className="p-3">${producto.precio.toFixed(2)}</td>
                  <td className="p-3">{producto.stock}</td>
                  <td className="p-3">{producto.activo ? 'Activo' : 'Inactivo'}</td>
                  <td className="p-3 flex gap-2">
                    <button onClick={() => editar(producto)} className="text-indigo-600 hover:underline">Editar</button>
                    {producto.activo && <button onClick={() => void desactivar(producto)} className="text-red-600 hover:underline">Desactivar</button>}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default AdminProductos;
