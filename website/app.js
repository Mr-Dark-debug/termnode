/* ═══════════════════════════════════════════
   TermNode — Interactive JavaScript
   ═══════════════════════════════════════════ */

(function () {
    'use strict';

    // ── Nav scroll effect ──
    const nav = document.getElementById('nav');
    let lastScroll = 0;

    function onScroll() {
        const y = window.scrollY;
        if (y > 40) {
            nav.classList.add('scrolled');
        } else {
            nav.classList.remove('scrolled');
        }
        lastScroll = y;
    }

    window.addEventListener('scroll', onScroll, { passive: true });

    // ── Mobile menu toggle ──
    const navToggle = document.getElementById('navToggle');
    const navLinks = document.getElementById('navLinks');

    if (navToggle && navLinks) {
        navToggle.addEventListener('click', function () {
            navLinks.classList.toggle('open');
            navToggle.classList.toggle('active');
        });

        // Close on link click
        navLinks.querySelectorAll('a').forEach(function (link) {
            link.addEventListener('click', function () {
                navLinks.classList.remove('open');
                navToggle.classList.remove('active');
            });
        });
    }

    // ── Scroll reveal ──
    const reveals = document.querySelectorAll('.reveal');

    function checkReveal() {
        const windowHeight = window.innerHeight;
        reveals.forEach(function (el) {
            const top = el.getBoundingClientRect().top;
            if (top < windowHeight - 60) {
                el.classList.add('visible');
            }
        });
    }

    window.addEventListener('scroll', checkReveal, { passive: true });
    window.addEventListener('resize', checkReveal, { passive: true });
    // Initial check
    checkReveal();

    // ── Copy to clipboard ──
    document.querySelectorAll('.code-copy').forEach(function (btn) {
        btn.addEventListener('click', function () {
            var code = btn.getAttribute('data-code');
            if (!code) return;

            // Decode HTML entities
            var textarea = document.createElement('textarea');
            textarea.innerHTML = code;
            var decoded = textarea.value;

            navigator.clipboard.writeText(decoded).then(function () {
                btn.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>';
                btn.style.color = '#10B981';
                setTimeout(function () {
                    btn.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>';
                    btn.style.color = '';
                }, 2000);
            });
        });
    });

    // ── GitHub star count ──
    var starEl = document.getElementById('starCount');
    if (starEl) {
        fetch('https://api.github.com/repos/Mr-Dark-debug/termnode', {
            headers: { 'Accept': 'application/vnd.github.v3+json' }
        })
            .then(function (r) { return r.json(); })
            .then(function (data) {
                if (data.stargazers_count !== undefined) {
                    starEl.textContent = '\u2605 ' + data.stargazers_count;
                } else {
                    starEl.textContent = '\u2605 Star us!';
                }
            })
            .catch(function () {
                starEl.textContent = '\u2605 Star us!';
            });
    }

    // ── GitHub contributors ──
    var contribEl = document.getElementById('contributors');
    if (contribEl) {
        fetch('https://api.github.com/repos/Mr-Dark-debug/termnode/contributors?per_page=8', {
            headers: { 'Accept': 'application/vnd.github.v3+json' }
        })
            .then(function (r) { return r.json(); })
            .then(function (users) {
                if (!Array.isArray(users) || users.length === 0) return;
                // Keep the first avatar (owner), add the rest
                var html = '';
                users.forEach(function (u) {
                    html += '<a href="' + u.html_url + '" target="_blank" rel="noopener" class="contributor-avatar">'
                        + '<img src="' + u.avatar_url + '&s=48" alt="' + u.login + '" width="40" height="40" loading="lazy">'
                        + '</a>';
                });
                contribEl.innerHTML = html;
            })
            .catch(function () {
                // Keep existing placeholder
            });
    }

    // ── Terminal typing animation ──
    var termBody = document.getElementById('termBody');
    if (termBody) {
        // Subtle pulse on the battery bar
        var bar = termBody.querySelector('.bar-filled');
        if (bar) {
            var origWidth = bar.style.width;
            setInterval(function () {
                var delta = Math.floor(Math.random() * 5) - 2;
                var base = parseInt(origWidth, 10);
                var newWidth = Math.max(80, Math.min(95, base + delta));
                bar.style.width = newWidth + '%';
            }, 3000);
        }
    }

    // ── Smooth scroll for anchor links ──
    document.querySelectorAll('a[href^="#"]').forEach(function (link) {
        link.addEventListener('click', function (e) {
            var target = document.querySelector(link.getAttribute('href'));
            if (target) {
                e.preventDefault();
                target.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }
        });
    });

})();
