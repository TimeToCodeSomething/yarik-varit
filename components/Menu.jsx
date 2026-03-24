/* ═══════════════════════════════════════════
   ЯРИК ВАРИТ — components/Menu.jsx
   React-компонент меню с вкладками.
   Подключается через index.html если нужен
   React-рендер, либо используется как
   reference для будущего SPA.
   ═══════════════════════════════════════════ */

import { useState } from 'react';
import { menuCategories, menuItems } from '../data/menu.js';

function MenuItem({ item }) {
  return (
    <div className="menu-item">
      {item.tag && <span className="bakery-tag">{item.tag}</span>}
      <p className="item-number">{item.num}</p>
      <h3 className="item-name">{item.name}</h3>
      <p className="item-desc">{item.desc}</p>
      <div className="item-price">
        {item.price}
        {item.volume && <span>{item.volume}</span>}
      </div>
    </div>
  );
}

export default function Menu() {
  const [active, setActive] = useState('espresso');
  const items = menuItems[active] || [];

  return (
    <section className="menu" id="menu">
      <div className="section-inner">

        <div className="menu-header reveal">
          <div>
            <p className="section-label">Меню</p>
            <h2 className="menu-title"><em>Меню</em></h2>
          </div>
          <div className="menu-cats">
            {menuCategories.map(cat => (
              <button
                key={cat.id}
                className={`cat-btn${active === cat.id ? ' active' : ''}`}
                onClick={() => setActive(cat.id)}
              >
                {cat.label}
              </button>
            ))}
          </div>
        </div>

        <div className="menu-grid reveal">
          {items.map(item => (
            <MenuItem key={item.num} item={item} />
          ))}
        </div>

      </div>
    </section>
  );
}
