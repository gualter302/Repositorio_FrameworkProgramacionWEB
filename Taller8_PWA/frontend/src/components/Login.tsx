// src/components/Login.tsx
import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { loginRequest } from '../services/api';
import { KeyRound, LogIn, Mail } from 'lucide-react';
import { AuthShell } from './Ui';

const Login = () => {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false); // true mientras se espera la respuesta

  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // Llamada real al backend: POST /api/login (ver src/services/api.ts)
      const data = await loginRequest(email, password);
      login(data);
      navigate('/'); // Redirigimos al Dashboard
    } catch (err) {
      // Muestra el mensaje del backend ("Credenciales incorrectas") o el de conexión
      setError(err instanceof Error ? err.message : 'Error inesperado al iniciar sesión');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthShell
      kicker="Acceso seguro"
      title="Vuelve a tu catálogo"
      description="Administra productos, pedidos y tu red desde un mismo espacio."
      footer={<>¿Aún no tienes cuenta? <Link to="/registro" className="font-semibold text-cyan-300 hover:text-cyan-200">Crear cuenta</Link></>}
    >
        {error && (
          <div className="alert-error mb-6" role="alert">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <label className="block">
            <span className="field-label">Correo electrónico</span>
            <span className="relative block">
              <Mail aria-hidden="true" className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" size={18} />
              <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="form-input form-input--with-icon"
              placeholder="admin@upse.edu.ec"
              autoComplete="email"
              required
              />
            </span>
          </label>

          <label className="block">
            <span className="field-label">Contraseña</span>
            <span className="relative block">
              <KeyRound aria-hidden="true" className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" size={18} />
              <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="form-input form-input--with-icon"
              placeholder="••••••"
              autoComplete="current-password"
              required
              />
            </span>
          </label>

          <button
            type="submit"
            disabled={loading}
            className="btn-primary w-full"
          >
            <LogIn size={18} aria-hidden="true" />
            {loading ? 'Ingresando...' : 'Iniciar Sesión'}
          </button>
        </form>
        <div className="mt-6 rounded-2xl border border-cyan-200/15 bg-cyan-300/[0.05] p-4 text-sm text-slate-300">
          <p className="mb-3 font-semibold text-cyan-200">Cuentas de prueba</p>
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="rounded-xl border border-white/10 bg-slate-950/35 p-3">
              <p className="mb-2 font-semibold text-white">Administrador</p>
              <p className="break-all"><span className="text-slate-500">Correo:</span> admin@upse.edu.ec</p>
              <p><span className="text-slate-500">Contraseña:</span> la definida en Render</p>
            </div>
            <div className="rounded-xl border border-white/10 bg-slate-950/35 p-3">
              <p className="mb-2 font-semibold text-white">Cliente</p>
              <p className="break-all"><span className="text-slate-500">Correo:</span> cliente@upse.edu.ec</p>
              <p><span className="text-slate-500">Contraseña:</span> 123456</p>
            </div>
          </div>
        </div>
    </AuthShell>
  );
};

export default Login;
