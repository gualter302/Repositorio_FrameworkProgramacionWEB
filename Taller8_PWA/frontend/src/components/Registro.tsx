import { useState, type FormEvent } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { registerRequest } from '../services/api';

const Registro = () => {
  const [parametros] = useSearchParams();
  const [nombre, setNombre] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [ref, setRef] = useState(parametros.get('ref') ?? '');
  const [error, setError] = useState('');
  const [guardando, setGuardando] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const registrar = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    setGuardando(true);
    try {
      const sesion = await registerRequest({ nombre, email, password, ref: ref || undefined });
      login(sesion);
      navigate('/mi-red');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudo crear la cuenta');
    } finally {
      setGuardando(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-100 px-4 py-8">
      <div className="max-w-md w-full bg-white rounded-xl shadow-md p-8 border border-slate-200">
        <h1 className="text-3xl font-bold text-slate-900 text-center">Crear cuenta</h1>
        <p className="text-slate-500 text-center mt-2 mb-7">Regístrate y forma parte de la red MultiCatálogo.</p>

        {error && <div className="bg-red-50 text-red-700 p-3 rounded-lg mb-5 border border-red-200">{error}</div>}

        <form onSubmit={registrar} className="space-y-4">
          <label className="block">
            <span className="block text-sm font-medium text-slate-700 mb-1">Nombre completo</span>
            <input required value={nombre} onChange={(event) => setNombre(event.target.value)} className="w-full px-4 py-3 rounded-lg border border-slate-300" />
          </label>
          <label className="block">
            <span className="block text-sm font-medium text-slate-700 mb-1">Correo electrónico</span>
            <input type="email" required value={email} onChange={(event) => setEmail(event.target.value)} className="w-full px-4 py-3 rounded-lg border border-slate-300" />
          </label>
          <label className="block">
            <span className="block text-sm font-medium text-slate-700 mb-1">Contraseña</span>
            <input type="password" required minLength={6} value={password} onChange={(event) => setPassword(event.target.value)} className="w-full px-4 py-3 rounded-lg border border-slate-300" />
          </label>
          <label className="block">
            <span className="block text-sm font-medium text-slate-700 mb-1">Código de referido (opcional)</span>
            <input value={ref} onChange={(event) => setRef(event.target.value.toUpperCase())} className="w-full px-4 py-3 rounded-lg border border-slate-300" placeholder="Ejemplo: ANA" />
          </label>
          <button type="submit" disabled={guardando} className="w-full bg-indigo-600 text-white font-bold py-3 rounded-lg hover:bg-indigo-700 disabled:opacity-60">
            {guardando ? 'Creando cuenta...' : 'Registrarme'}
          </button>
        </form>

        <p className="text-center text-sm text-slate-600 mt-6">
          ¿Ya tienes cuenta? <Link to="/login" className="text-indigo-600 font-semibold hover:underline">Iniciar sesión</Link>
        </p>
      </div>
    </div>
  );
};

export default Registro;
