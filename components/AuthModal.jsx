/* ═══════════════════════════════════════════
   ЯРИК ВАРИТ — components/AuthModal.jsx
   ═══════════════════════════════════════════ */

import { useState } from 'react';

const GoogleIcon = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
    <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
    <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
    <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05"/>
    <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
  </svg>
);

export default function AuthModal({ isOpen, onClose, onLogin }) {
  const [tab, setTab]         = useState('login');
  const [email, setEmail]     = useState('');
  const [pass, setPass]       = useState('');
  const [name, setName]       = useState('');
  const [regEmail, setRegEmail] = useState('');
  const [regPass, setRegPass]   = useState('');

  const capitalize = s => s ? s.charAt(0).toUpperCase() + s.slice(1) : s;

  const handleLogin = () => {
    if (!email || !pass) { alert('Введите email и пароль'); return; }
    onLogin({ name: capitalize(email.split('@')[0]), email });
    onClose();
  };

  const handleRegister = () => {
    if (!name || !regEmail || !regPass) { alert('Заполните все поля'); return; }
    if (regPass.length < 6) { alert('Пароль минимум 6 символов'); return; }
    onLogin({ name, email: regEmail });
    onClose();
  };

  const handleGoogle = () => {
    onLogin({ name: 'Алексей Иванов', email: 'alex@gmail.com' });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className={`modal-backdrop${isOpen ? ' open' : ''}`} onClick={e => e.target === e.currentTarget && onClose()}>
      <div className="modal">
        <button className="modal-close" onClick={onClose} />
        <p className="modal-eyebrow">Добро пожаловать</p>
        <h2 className="modal-title">Ярик<br /><em>узнаёт своих</em></h2>

        <div className="auth-tabs">
          <button className={`auth-tab${tab === 'login' ? ' active' : ''}`} onClick={() => setTab('login')}>Войти</button>
          <button className={`auth-tab${tab === 'register' ? ' active' : ''}`} onClick={() => setTab('register')}>Аккаунт</button>
        </div>

        {tab === 'login' ? (
          <div>
            <div className="auth-form">
              <input type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} />
              <input type="password" placeholder="Пароль" value={pass} onChange={e => setPass(e.target.value)} />
              <button className="btn-auth-submit" onClick={handleLogin}>Войти</button>
            </div>
            <div className="auth-divider"><span>или</span></div>
            <button className="btn-google" onClick={handleGoogle}><GoogleIcon />Войти через Google</button>
          </div>
        ) : (
          <div>
            <div className="auth-form">
              <input type="text" placeholder="Имя" value={name} onChange={e => setName(e.target.value)} />
              <input type="email" placeholder="Email" value={regEmail} onChange={e => setRegEmail(e.target.value)} />
              <input type="password" placeholder="Пароль (мин. 6 символов)" value={regPass} onChange={e => setRegPass(e.target.value)} />
              <button className="btn-auth-submit" onClick={handleRegister}>Создать аккаунт</button>
            </div>
            <div className="auth-divider"><span>или</span></div>
            <button className="btn-google" onClick={handleGoogle}><GoogleIcon />Через Google</button>
          </div>
        )}
      </div>
    </div>
  );
}
