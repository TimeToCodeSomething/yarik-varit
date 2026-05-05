/* ═══════════════════════════════════════════
   ЯРИК ВАРИТ — main.js
   ═══════════════════════════════════════════ */

const API_URL = 'http://localhost:8080';

let currentUser = null;

/* ── MENU FROM API ── */
async function loadMenu() {
  try {
    const res = await fetch(`${API_URL}/menu`);
    if (!res.ok) throw new Error('Ошибка загрузки меню');
    const items = await res.json();

    // Группируем по категориям
    const byCategory = {};
    items.forEach(item => {
      if (!byCategory[item.category]) byCategory[item.category] = [];
      byCategory[item.category].push(item);
    });

    // Наполняем каждую панель
    Object.entries(byCategory).forEach(([cat, catItems]) => {
      const panel = document.getElementById('panel-' + cat);
      if (!panel) return;
      const grid = panel.querySelector('.menu-grid');
      if (!grid) return;

      grid.innerHTML = catItems.map((item, i) => {
        const num = cat === 'bakery' ? `Б-${String(i + 1).padStart(2, '0')}` : String(i + 1).padStart(2, '0');
        const vol = item.vol > 0 ? `<span>/ ${item.vol} мл</span>` : '';
        return `
          <div class="menu-item">
            <p class="item-number">${num}</p>
            <h3 class="item-name">${item.name}</h3>
            <div class="item-price">${item.price} ₽${vol}</div>
          </div>`;
      }).join('');
    });

    // Переинициализируем курсор для новых элементов
    initCursorTargets();

  } catch (err) {
    console.warn('Меню из БД не загрузилось, показываем статику:', err);
  }
}

document.addEventListener('DOMContentLoaded', loadMenu);

/* ── BURGER MENU ── */
function toggleMenu() {
  const b = document.getElementById('burger');
  const o = document.getElementById('mobileOverlay');
  b.classList.toggle('open');
  o.classList.toggle('open');
  document.body.style.overflow = o.classList.contains('open') ? 'hidden' : '';
}
function closeMenu() {
  document.getElementById('burger').classList.remove('open');
  document.getElementById('mobileOverlay').classList.remove('open');
  document.body.style.overflow = '';
}

/* ── AUTH MODAL ── */
function openAuth() {
  document.getElementById('authModal').classList.add('open');
  document.body.style.overflow = 'hidden';
}
function closeAuth() {
  document.getElementById('authModal').classList.remove('open');
  document.body.style.overflow = '';
}
document.getElementById('authModal').addEventListener('click', function(e) {
  if (e.target === this) closeAuth();
});

function switchAuthTab(tab, el) {
  document.querySelectorAll('.auth-tab').forEach(t => t.classList.remove('active'));
  el.classList.add('active');
  document.getElementById('tab-login').style.display = tab === 'login' ? 'block' : 'none';
  document.getElementById('tab-register').style.display = tab === 'register' ? 'block' : 'none';
}

async function doLogin() {
  const username = document.getElementById('loginEmail').value.trim();
  const pass     = document.getElementById('loginPass').value.trim();
  if (!username || !pass) { alert('Введите логин и пароль'); return; }

  try {
    const res = await fetch(`${API_URL}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password: pass }),
    });

    if (res.status === 401) { alert('Неверный логин или пароль'); return; }
    if (!res.ok) { alert('Ошибка сервера'); return; }

    const data = await res.json();
    localStorage.setItem('token', data.token);
    loginUser({ name: capitalize(username) });
    closeAuth();
  } catch (err) {
    alert('Не удалось подключиться к серверу');
  }
}

async function doRegister() {
  const username = document.getElementById('regName').value.trim();
  const pass     = document.getElementById('regPass').value.trim();
  if (!username || !pass) { alert('Заполните все поля'); return; }
  if (pass.length < 6)    { alert('Пароль минимум 6 символов'); return; }

  try {
    const res = await fetch(`${API_URL}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password: pass }),
    });

    if (res.status === 409) { alert('Логин уже занят'); return; }
    if (!res.ok) { alert('Ошибка регистрации'); return; }

    // После регистрации сразу логиним
    const loginRes = await fetch(`${API_URL}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password: pass }),
    });
    const data = await loginRes.json();
    localStorage.setItem('token', data.token);
    loginUser({ name: capitalize(username) });
    closeAuth();
  } catch (err) {
    alert('Не удалось подключиться к серверу');
  }
}

function doGoogleLogin() {
  alert('Скоро появится');
}

function capitalize(s) {
  return s ? s.charAt(0).toUpperCase() + s.slice(1) : s;
}

function loginUser(user) {
  currentUser = user;
  document.getElementById('btnNavAuth').style.display = 'none';
  document.getElementById('btnNavProfile').style.display = 'block';
  document.getElementById('mobAuthBtn').style.display = 'none';
  document.getElementById('mobProfileBtn').style.display = 'block';
  document.getElementById('profileName').textContent = user.name;
  document.getElementById('profileAvatar').textContent = user.name.charAt(0).toUpperCase();
}

function doLogout() {
  currentUser = null;
  localStorage.removeItem('token');
  closeProfile();
  document.getElementById('btnNavAuth').style.display = '';
  document.getElementById('btnNavProfile').style.display = 'none';
  document.getElementById('mobAuthBtn').style.display = '';
  document.getElementById('mobProfileBtn').style.display = 'none';
}

/* ── PROFILE PANEL ── */
function openProfile() {
  if (!currentUser) { openAuth(); return; }
  document.getElementById('profilePanel').classList.add('open');
  document.getElementById('panelOverlay').classList.add('open');
  document.body.style.overflow = 'hidden';
  setTimeout(() => { document.getElementById('bonusBar').style.width = '42%'; }, 350);
}
function closeProfile() {
  document.getElementById('profilePanel').classList.remove('open');
  document.getElementById('panelOverlay').classList.remove('open');
  document.body.style.overflow = '';
}

/* ── MENU TABS ── */
function switchMenu(id, el) {
  document.querySelectorAll('.cat-btn').forEach(b => b.classList.remove('active'));
  document.querySelectorAll('.menu-panel').forEach(p => p.classList.remove('active'));
  el.classList.add('active');
  document.getElementById('panel-' + id).classList.add('active');
}

/* ── CUSTOM CURSOR ── */
const cursor = document.getElementById('cursor');
const ring   = document.getElementById('cursorRing');
let mx = 0, my = 0, rx = 0, ry = 0;

document.addEventListener('mousemove', e => {
  mx = e.clientX; my = e.clientY;
  cursor.style.left = mx + 'px';
  cursor.style.top  = my + 'px';
});

(function animRing() {
  rx += (mx - rx) * .12;
  ry += (my - ry) * .12;
  ring.style.left = rx + 'px';
  ring.style.top  = ry + 'px';
  requestAnimationFrame(animRing);
})();

// Expand ring on interactive elements
function initCursorTargets() {
  document.querySelectorAll('a, button, .menu-item, .delivery-card').forEach(el => {
    el.addEventListener('mouseenter', () => {
      ring.style.width  = '50px';
      ring.style.height = '50px';
      ring.style.borderColor = 'rgba(201,169,110,.7)';
    });
    el.addEventListener('mouseleave', () => {
      ring.style.width  = '32px';
      ring.style.height = '32px';
      ring.style.borderColor = 'rgba(201,169,110,.4)';
    });
  });
}
initCursorTargets();

// Dark cursor on gold surfaces
function setCursorDark() {
  cursor.style.background  = '#1a1710';
  ring.style.borderColor   = 'rgba(20,18,12,.6)';
  ring.style.background    = 'rgba(20,18,12,.05)';
}
function setCursorLight() {
  cursor.style.background  = '';
  ring.style.borderColor   = '';
  ring.style.background    = '';
}

document.querySelectorAll(
  '.ticker-wrap, .btn-auth-submit, .btn-order, .btn-nav-profile, .bakery-tag'
).forEach(el => {
  el.addEventListener('mouseenter', setCursorDark);
  el.addEventListener('mouseleave', setCursorLight);
});

document.querySelectorAll('.cat-btn, .mob-auth-btn, .btn-google').forEach(el => {
  el.addEventListener('mouseenter', setCursorDark);
  el.addEventListener('mouseleave', setCursorLight);
});

/* ── SCROLL REVEAL ── */
const revealObs = new IntersectionObserver(entries => {
  entries.forEach(e => { if (e.isIntersecting) e.target.classList.add('visible'); });
}, { threshold: .1 });

document.querySelectorAll('.reveal').forEach(el => revealObs.observe(el));
