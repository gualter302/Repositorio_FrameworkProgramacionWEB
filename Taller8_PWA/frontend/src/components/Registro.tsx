import { useState, type FormEvent } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { registerRequest } from '../services/api';
import { Mail, Network, UserPlus, UserRound, KeyRound } from 'lucide-react';
import { AuthShell } from './Ui';

const Registro = () => {
  const [parametros] = useSearchParams();
  const [nombre, setNombre] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [ref, setRef] = useState((parametros.get('ref') ?? '').toUpperCase());
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
    <AuthShell
      kicker="Nueva cuenta"
      title="Forma parte de la red"
      description="Crea tu perfil y recibe un código propio para compartir MultiCatálogo."
      footer={<>¿Ya tienes cuenta? <Link to="/login" className="font-semibold text-cyan-300 hover:text-cyan-200">Iniciar sesión</Link></>}
    >
        {error && <div className="alert-error mb-5" role="alert">{error}</div>}

        <form onSubmit={registrar} className="space-y-4">
          <label className="block">
            <span className="field-label">Nombre completo</span>
            <span className="relative block"><UserRound aria-hidden="true" className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" size={18} /><input required value={nombre} onChange={(event) => setNombre(event.target.value)} className="form-input form-input--with-icon" autoComplete="name" /></span>
          </label>
          <label className="block">
            <span className="field-label">Correo electrónico</span>
            <span className="relative block"><Mail aria-hidden="true" className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" size={18} /><input type="email" required value={email} onChange={(event) => setEmail(event.target.value)} className="form-input form-input--with-icon" autoComplete="email" /></span>
          </label>
          <label className="block">
            <span className="field-label">Contraseña</span>
            <span className="relative block"><KeyRound aria-hidden="true" className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" size={18} /><input type="password" required minLength={6} value={password} onChange={(event) => setPassword(event.target.value)} className="form-input form-input--with-icon" autoComplete="new-password" /></span>
          </label>
          <label className="block">
            <span className="field-label">Código de referido (opcional)</span>
            <span className="relative block"><Network aria-hidden="true" className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" size={18} /><input value={ref} onChange={(event) => setRef(event.target.value.toUpperCase())} className="form-input form-input--with-icon uppercase" placeholder="Ejemplo: ANA" /></span>
          </label>
          <button type="submit" disabled={guardando} className="btn-primary w-full">
            <UserPlus size={18} aria-hidden="true" />
            {guardando ? 'Creando cuenta...' : 'Registrarme'}
          </button>
        </form>
    </AuthShell>
  );
};

export default Registro;
