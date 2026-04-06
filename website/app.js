/**
 * TermNode — Professional Website Interactions
 * Parallax · Typing Animation · Scroll Reveals · Micro-interactions
 */

document.addEventListener('DOMContentLoaded', () => {
    initNav();
    initScrollProgress();
    initScrollReveal();
    initParallax();
    initParticles();
    initTerminalTyping();
    initInstallTabs();
    initCopyButtons();
    initCounterAnimation();
    initRippleEffect();
    initContributors();
    initChartAnimation();
    initDocsSidebar();
});

/* ═══════════════════════════════════════
   NAVIGATION
   ═══════════════════════════════════════ */
function initNav() {
    const nav = document.getElementById('mainNav');
    const hamburger = document.getElementById('hamburger');
    const mobileMenu = document.getElementById('mobileMenu');

    // Scroll state
    if (!nav.classList.contains('scrolled')) {
        window.addEventListener('scroll', () => {
            if (window.scrollY > 40) {
                nav.classList.add('scrolled');
            } else {
                nav.classList.remove('scrolled');
            }
        }, { passive: true });
    }

    // Hamburger
    if (hamburger && mobileMenu) {
        hamburger.addEventListener('click', () => {
            hamburger.classList.toggle('open');
            mobileMenu.classList.toggle('open');
            document.body.style.overflow = mobileMenu.classList.contains('open') ? 'hidden' : '';
        });

        // Close on link click
        mobileMenu.querySelectorAll('a').forEach(link => {
            link.addEventListener('click', () => closeMobile());
        });
    }
}

function closeMobile() {
    const hamburger = document.getElementById('hamburger');
    const mobileMenu = document.getElementById('mobileMenu');
    if (hamburger) hamburger.classList.remove('open');
    if (mobileMenu) mobileMenu.classList.remove('open');
    document.body.style.overflow = '';
}

/* ═══════════════════════════════════════
   SCROLL PROGRESS
   ═══════════════════════════════════════ */
function initScrollProgress() {
    const bar = document.getElementById('scrollProgress');
    if (!bar) return;

    window.addEventListener('scroll', () => {
        const scrollTop = window.scrollY;
        const docHeight = document.documentElement.scrollHeight - window.innerHeight;
        const progress = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0;
        bar.style.width = progress + '%';
    }, { passive: true });
}

/* ═══════════════════════════════════════
   SCROLL REVEAL
   ═══════════════════════════════════════ */
function initScrollReveal() {
    const reveals = document.querySelectorAll('.reveal');
    if (!reveals.length) return;

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
                observer.unobserve(entry.target);
            }
        });
    }, {
        threshold: 0.1,
        rootMargin: '0px 0px -40px 0px'
    });

    reveals.forEach(el => observer.observe(el));
}

/* ═══════════════════════════════════════
   PARALLAX
   ═══════════════════════════════════════ */
function initParallax() {
    const hero = document.getElementById('hero');
    if (!hero) return;

    const parallaxElements = hero.querySelectorAll('[data-speed]');

    window.addEventListener('scroll', () => {
        const scrollY = window.scrollY;
        const heroHeight = hero.offsetHeight;

        if (scrollY > heroHeight) return;

        parallaxElements.forEach(el => {
            const speed = parseFloat(el.dataset.speed) || 0.2;
            el.style.transform = `translateY(${scrollY * speed}px)`;
        });
    }, { passive: true });
}

/* ═══════════════════════════════════════
   FLOATING PARTICLES
   ═══════════════════════════════════════ */
function initParticles() {
    const heroBg = document.getElementById('heroBg');
    if (!heroBg) return;

    const particleCount = 20;
    for (let i = 0; i < particleCount; i++) {
        const particle = document.createElement('div');
        particle.className = 'particle';
        particle.style.left = Math.random() * 100 + '%';
        particle.style.width = (Math.random() * 3 + 2) + 'px';
        particle.style.height = particle.style.width;
        particle.style.animationDuration = (Math.random() * 15 + 10) + 's';
        particle.style.animationDelay = (Math.random() * 15) + 's';
        particle.style.opacity = 0;
        heroBg.appendChild(particle);
    }
}

/* ═══════════════════════════════════════
   TERMINAL TYPING ANIMATION
   ═══════════════════════════════════════ */
function initTerminalTyping() {
    const terminalBody = document.getElementById('terminalBody');
    if (!terminalBody) return;

    const lines = [
        { type: 'command', cmd: 'help' },
        { type: 'output', text: 'Available commands: battery, services, iot log, config, exit.' },
        { type: 'command', cmd: 'battery' },
        { type: 'output', text: 'Level: 87% | Status: Discharging | Health: Good | Temp: 31.2°C' },
        { type: 'command', cmd: 'services' },
        { type: 'output', text: 'SSH: [RUNNING] | HTTP: [STOPPED] | WakeLock: [ACTIVE]' },
        { type: 'command', cmd: 'iot log' },
        { type: 'output', text: '12:04:22 - ESP32_Temp: 24.5°C\n12:05:01 - LivingRoom_Light: ON' },
    ];

    let lineIndex = 0;
    let charIndex = 0;
    let currentEl = null;

    const observer = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting) {
            observer.disconnect();
            typeNextLine();
        }
    }, { threshold: 0.3 });

    observer.observe(terminalBody);

    function typeNextLine() {
        if (lineIndex >= lines.length) {
            // Add blinking cursor at the end
            const cursorLine = document.createElement('div');
            cursorLine.className = 'term-line visible';
            cursorLine.innerHTML = '<span class="term-prompt">&gt;</span><span class="term-cmd"><span class="typing-cursor"></span></span>';
            terminalBody.appendChild(cursorLine);
            return;
        }

        const line = lines[lineIndex];

        if (line.type === 'command') {
            // Create command line
            const lineEl = document.createElement('div');
            lineEl.className = 'term-line';
            lineEl.innerHTML = `<span class="term-prompt">&gt;</span><span class="term-cmd"></span>`;
            terminalBody.appendChild(lineEl);
            currentEl = lineEl.querySelector('.term-cmd');
            charIndex = 0;
            typeChar(line.cmd);
        } else {
            // Create output
            const outputEl = document.createElement('span');
            outputEl.className = 'term-output';
            outputEl.textContent = line.text;
            terminalBody.appendChild(outputEl);

            requestAnimationFrame(() => {
                outputEl.classList.add('visible');
            });

            lineIndex++;
            setTimeout(typeNextLine, 300);
        }
    }

    function typeChar(text) {
        if (charIndex < text.length) {
            // Show the line element
            currentEl.closest('.term-line').classList.add('visible');
            currentEl.textContent = text.substring(0, charIndex + 1);
            charIndex++;
            setTimeout(() => typeChar(text), 30 + Math.random() * 40);
        } else {
            lineIndex++;
            setTimeout(typeNextLine, 400);
        }
    }
}

/* ═══════════════════════════════════════
   INSTALL TABS
   ═══════════════════════════════════════ */
function initInstallTabs() {
    const tabs = document.querySelectorAll('.install-tab');
    const panels = document.querySelectorAll('.install-panel');

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const target = tab.dataset.tab;

            tabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');

            panels.forEach(p => p.classList.remove('active'));
            const panel = document.getElementById('panel-' + target);
            if (panel) {
                panel.classList.add('active');
                // Re-trigger code copy buttons
                panel.querySelectorAll('.code-copy').forEach(btn => {
                    btn.textContent = 'Copy';
                    btn.classList.remove('copied');
                });
            }
        });
    });
}

/* ═══════════════════════════════════════
   COPY BUTTONS
   ═══════════════════════════════════════ */
function initCopyButtons() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('.code-copy');
        if (!btn) return;

        const codeBlock = btn.closest('.code-block');
        if (!codeBlock) return;

        const code = codeBlock.querySelector('code');
        if (!code) return;

        // Get text content, stripping HTML tags
        const text = code.innerText || code.textContent;

        navigator.clipboard.writeText(text).then(() => {
            btn.textContent = 'Copied!';
            btn.classList.add('copied');

            setTimeout(() => {
                btn.textContent = 'Copy';
                btn.classList.remove('copied');
            }, 2000);
        }).catch(() => {
            // Fallback for older browsers
            const textarea = document.createElement('textarea');
            textarea.value = text;
            textarea.style.position = 'fixed';
            textarea.style.opacity = '0';
            document.body.appendChild(textarea);
            textarea.select();
            document.execCommand('copy');
            document.body.removeChild(textarea);

            btn.textContent = 'Copied!';
            btn.classList.add('copied');
            setTimeout(() => {
                btn.textContent = 'Copy';
                btn.classList.remove('copied');
            }, 2000);
        });
    });
}

/* ═══════════════════════════════════════
   COUNTER ANIMATION
   ═══════════════════════════════════════ */
function initCounterAnimation() {
    const counters = document.querySelectorAll('.stat-value[data-count]');
    if (!counters.length) return;

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                animateCounter(entry.target);
                observer.unobserve(entry.target);
            }
        });
    }, { threshold: 0.5 });

    counters.forEach(el => observer.observe(el));
}

function animateCounter(el) {
    const target = parseInt(el.dataset.count, 10);
    const duration = 1500;
    const start = performance.now();

    function update(now) {
        const elapsed = now - start;
        const progress = Math.min(elapsed / duration, 1);

        // Ease out cubic
        const eased = 1 - Math.pow(1 - progress, 3);
        const current = Math.round(target * eased);

        el.textContent = current.toLocaleString();

        if (progress < 1) {
            requestAnimationFrame(update);
        }
    }

    requestAnimationFrame(update);
}

/* ═══════════════════════════════════════
   RIPPLE EFFECT ON BUTTONS
   ═══════════════════════════════════════ */
function initRippleEffect() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('.btn');
        if (!btn) return;

        const ripple = document.createElement('span');
        ripple.className = 'ripple';

        const rect = btn.getBoundingClientRect();
        const size = Math.max(rect.width, rect.height);
        ripple.style.width = ripple.style.height = size + 'px';
        ripple.style.left = (e.clientX - rect.left - size / 2) + 'px';
        ripple.style.top = (e.clientY - rect.top - size / 2) + 'px';

        btn.appendChild(ripple);
        ripple.addEventListener('animationend', () => ripple.remove());
    });
}

/* ═══════════════════════════════════════
   GITHUB CONTRIBUTORS
   ═══════════════════════════════════════ */
async function initContributors() {
    const slots = document.querySelectorAll('.contributor-slot');
    if (!slots.length) return;

    try {
        const response = await fetch('https://api.github.com/repos/Mr-Dark-debug/termnode/contributors?per_page=5');
        if (!response.ok) throw new Error('Failed');
        const data = await response.json();

        slots.forEach((slot, i) => {
            if (data[i]) {
                slot.style.backgroundImage = `url(${data[i].avatar_url})`;
                slot.style.backgroundSize = 'cover';
                slot.title = data[i].login;
                slot.onclick = () => window.open(data[i].html_url, '_blank');
            }
        });
    } catch {
        slots.forEach((slot, i) => {
            const colors = ['#E6FFFA', '#EBF8FF', '#EEF2FF'];
            const initials = ['M', 'D', '?'];
            slot.style.backgroundColor = colors[i % 3];
            slot.innerText = initials[i % 3] || '?';
            slot.style.display = 'flex';
            slot.style.alignItems = 'center';
            slot.style.justifyContent = 'center';
            slot.style.fontSize = '0.75rem';
            slot.style.fontWeight = '700';
            slot.style.color = 'var(--teal-dark)';
            slot.style.fontFamily = 'var(--font-mono)';
        });
    }
}

/* ═══════════════════════════════════════
   CHART LINE ANIMATION
   ═══════════════════════════════════════ */
function initChartAnimation() {
    const path = document.querySelector('.chart-line');
    if (!path) return;

    const length = path.getTotalLength();
    path.style.strokeDasharray = length;
    path.style.strokeDashoffset = length;

    const area = document.querySelector('.chart-area');
    if (area) {
        const areaLength = area.getTotalLength ? 0 : 0;
        area.style.opacity = '0';
    }

    const observer = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting) {
            path.style.transition = 'stroke-dashoffset 2s ease-out';
            path.style.strokeDashoffset = '0';

            if (area) {
                setTimeout(() => { area.style.transition = 'opacity 1s ease'; area.style.opacity = '0.3'; }, 800);
            }

            observer.disconnect();
        }
    }, { threshold: 0.3 });

    observer.observe(path.closest('.stats-chart') || path);
}

/* ═══════════════════════════════════════
   DOCS SIDEBAR ACTIVE STATE
   ═══════════════════════════════════════ */
function initDocsSidebar() {
    const sidebarLinks = document.querySelectorAll('.docs-sidebar-nav a');
    const sections = [];

    sidebarLinks.forEach(link => {
        const id = link.getAttribute('href');
        if (id && id.startsWith('#')) {
            const section = document.querySelector(id);
            if (section) sections.push({ el: section, link: link });
        }
    });

    if (!sections.length) return;

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                sidebarLinks.forEach(l => l.classList.remove('active'));
                const match = sections.find(s => s.el === entry.target);
                if (match) match.link.classList.add('active');
            }
        });
    }, {
        rootMargin: '-80px 0px -60% 0px',
        threshold: 0
    });

    sections.forEach(s => observer.observe(s.el));
}
