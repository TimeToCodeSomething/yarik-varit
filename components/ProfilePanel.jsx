/* ═══════════════════════════════════════════
   ЯРИК ВАРИТ — components/ProfilePanel.jsx
   ═══════════════════════════════════════════ */

import { useEffect, useRef } from 'react';

const LEVELS = [
  { icon: '☕',  name: 'Любитель', range: '0 — 999 баллов',      cashback: '5%',  active: true  },
  { icon: '◈',  name: 'Мастер',   range: '1000 — 2999 баллов',  cashback: '8%',  active: false },
  { icon: '✦',  name: 'Сомелье',  range: 'от 3000 баллов',      cashback: '12%', active: false },
];

const HISTORY = [
  { name: 'Флэт Уайт + Круассан',    date: '2 мар 2025',  pts: '+18 б' },
  { name: 'Капучино × 2',            date: '28 фев 2025', pts: '+24 б' },
  { name: 'Фильтр Ethiopia',         date: '25 фев 2025', pts: '+12 б' },
  { name: 'Колд Брю + Тост авокадо', date: '20 фев 2025', pts: '+28 б' },
];

export default function ProfilePanel({ isOpen, user, onClose, onLogout }) {
  const barRef = useRef(null);

  useEffect(() => {
    if (isOpen && barRef.current) {
      setTimeout(() => { barRef.current.style.width = '42%'; }, 350);
    } else if (barRef.current) {
      barRef.current.style.width = '0';
    }
  }, [isOpen]);

  return (
    <>
      <div className={`panel-overlay${isOpen ? ' open' : ''}`} onClick={onClose} />
      <div className={`profile-panel${isOpen ? ' open' : ''}`}>

        <div className="panel-head">
          <div className="panel-title">Мой профиль</div>
          <button className="panel-close" onClick={onClose} />
        </div>

        <div className="panel-body">
          <div className="profile-avatar">{user?.name?.charAt(0)?.toUpperCase() ?? 'А'}</div>
          <div className="profile-name">{user?.name ?? 'Гость'}</div>
          <div className="profile-email">{user?.email ?? ''}</div>

          {/* BONUS CARD */}
          <div className="bonus-card">
            <div className="bonus-label">Бонусный счёт</div>
            <div className="bonus-points">840</div>
            <div className="bonus-sub">бонусных баллов · уровень <strong style={{ color: 'var(--gold)' }}>Любитель</strong></div>
            <div className="bonus-bar-wrap">
              <div className="bonus-bar" ref={barRef} />
            </div>
            <div className="bonus-bar-note">До уровня «Мастер» — 160 баллов</div>
          </div>

          {/* HOW BONUSES WORK */}
          <div className="bonus-info">
            <div className="bonus-info-title">Как работают бонусы</div>
            <p>С каждой покупки начисляется <strong>процент от суммы</strong> в виде баллов. Чем выше уровень — тем больше кешбэк.</p>
            <div className="bonus-info-uses">
              {[
                { icon: '☕', title: 'Оплата кофе', desc: '— спишите баллы в счёт любого напитка. 1 балл = 1 ₽' },
                { icon: '✂',  title: 'Скидка на заказ', desc: '— оплатите до 50% от суммы бонусами' },
                { icon: '🎁', title: 'Выпечка в подарок', desc: '— 300 баллов = любая выпечка бесплатно' },
              ].map(row => (
                <div className="bonus-use-row" key={row.title}>
                  <span className="bonus-use-icon">{row.icon}</span>
                  <span className="bonus-use-text"><strong style={{ color: 'var(--white-dim)' }}>{row.title}</strong> {row.desc}</span>
                </div>
              ))}
            </div>
          </div>

          {/* LEVELS */}
          <div className="profile-section-label">Уровни и кешбэк</div>
          <div className="levels">
            {LEVELS.map(lvl => (
              <div key={lvl.name} className={`level-row${lvl.active ? ' active' : ''}`}>
                <div className="level-left">
                  <span className="level-icon">{lvl.icon}</span>
                  <div>
                    <div className="level-name">{lvl.name}</div>
                    <div className="level-threshold">{lvl.range}</div>
                  </div>
                </div>
                <div style={{ textAlign: 'right' }}>
                  <div className="level-cashback">{lvl.cashback}</div>
                  <div style={{ fontSize: '8px', letterSpacing: '.15em', color: 'var(--white-muted)', textTransform: 'uppercase' }}>кешбэк</div>
                </div>
              </div>
            ))}
          </div>

          {/* HISTORY */}
          <div className="profile-section-label">История заказов</div>
          {HISTORY.map(h => (
            <div className="history-item" key={h.name}>
              <div>
                <div className="history-name">{h.name}</div>
                <div className="history-date">{h.date}</div>
              </div>
              <div className="history-pts">{h.pts}</div>
            </div>
          ))}

          <button className="btn-logout" onClick={onLogout}>Выйти из аккаунта</button>
        </div>
      </div>
    </>
  );
}
