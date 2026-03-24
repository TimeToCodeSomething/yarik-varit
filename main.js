/* ═══════════════════════════════════════════
   ЯРИК ВАРИТ — main.js
   ═══════════════════════════════════════════ */

let currentUser = null;

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

function doLogin() {
  const email = document.getElementById('loginEmail').value.trim();
  const pass  = document.getElementById('loginPass').value.trim();
  if (!email || !pass) { alert('Введите email и пароль'); return; }
  loginUser({ name: capitalize(email.split('@')[0]), email });
  closeAuth();
}

function doRegister() {
  const name  = document.getElementById('regName').value.trim();
  const email = document.getElementById('regEmail').value.trim();
  const pass  = document.getElementById('regPass').value.trim();
  if (!name || !email || !pass) { alert('Заполните все поля'); return; }
  if (pass.length < 6) { alert('Пароль минимум 6 символов'); return; }
  loginUser({ name, email });
  closeAuth();
}

function doGoogleLogin() {
  loginUser({ name: 'Алексей Иванов', email: 'alex@gmail.com' });
  closeAuth();
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
  document.getElementById('profileEmail').textContent = user.email;
  document.getElementById('profileAvatar').textContent = user.name.charAt(0).toUpperCase();
}

function doLogout() {
  currentUser = null;
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
