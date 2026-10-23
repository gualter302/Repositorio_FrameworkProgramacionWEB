// src/components/Login.tsx
import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { loginRequest } from '../services/api';
import { Check, KeyRound, LogIn, Mail, UserRound } from 'lucide-react';
import { AuthShell } from './Ui';

const cuentasPrueba = [
  { id: 'admin', nombre: 'Administrador', email: 'admin@upse.edu.ec', password: '123456123456' },
  { id: 'cliente', nombre: 'Cliente', email: 'cliente@upse.edu.ec', password: '123456' },
] as const;

const Login = () => {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false); // true mientras se espera la respuesta
  const [cuentaElegida, setCuentaElegida] = useState<string | null>(null);

  const navigate = useNavigate();
  const { login } = useAuth();

  const elegirCuenta = (cuenta: (typeof cuentasPrueba)[number]) => {
    setEmail(cuenta.email);
    setPassword(cuenta.password);
    setCuentaElegida(cuenta.id);
    setError('');
  };

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
              onChange={(e) => { setEmail(e.target.value); setCuentaElegida(null); }}
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
              onChange={(e) => { setPassword(e.target.value); setCuentaElegida(null); }}
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
          <p className="font-semibold text-cyan-200">Cuentas de prueba</p>
          <p className="mb-3 mt-1 text-xs leading-relaxed text-slate-500">Selecciona una cuenta para completar el formulario automáticamente.</p>
          <div className="grid gap-3">
            {cuentasPrueba.map((cuenta) => {
              const seleccionada = cuentaElegida === cuenta.id;
              return <button
                key={cuenta.id}
                type="button"
                aria-pressed={seleccionada}
                aria-label={`Usar cuenta de ${cuenta.nombre.toLowerCase()}`}
                onClick={() => elegirCuenta(cuenta)}
                className={`group min-w-0 rounded-xl border p-3 text-left transition ${seleccionada ? 'border-cyan-300/60 bg-cyan-300/10 shadow-[0_0_24px_rgba(48,214,223,0.08)]' : 'border-white/10 bg-slate-950/35 hover:border-cyan-300/35 hover:bg-cyan-300/[0.06]'}`}
              >
                <span className="flex items-center justify-between gap-3">
                  <span className="flex min-w-0 items-center gap-2 font-semibold text-white"><UserRound size={17} className="shrink-0 text-cyan-300" />{cuenta.nombre}</span>
                  <span className={`flex shrink-0 items-center gap-1 text-xs font-semibold ${seleccionada ? 'text-cyan-200' : 'text-slate-500 group-hover:text-cyan-200'}`}>{seleccionada && <Check size={14} />} {seleccionada ? 'Seleccionada' : 'Usar cuenta'}</span>
                </span>
                <span className="mt-2 block break-all text-xs text-slate-400">{cuenta.email}</span>
                <span className="mt-1 block text-xs text-slate-500">Contraseña: <span className="mono text-slate-300">{cuenta.password}</span></span>
              </button>;
            })}
          </div>
        </div>
    </AuthShell>
  );
};

export default Login;
