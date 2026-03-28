/**
 * TermNode — Marketing Site Interactivity
 * Focus: Clean, fast, and subtle animations
 */

document.addEventListener('DOMContentLoaded', () => {
    initCopyButtons();
    initScrollAnimations();
    initStatsChart();
    initContributors();
});

/**
 * Copy to Clipboard for code blocks
 */
function initCopyButtons() {
    const buttons = document.querySelectorAll('.code-copy');
    
    buttons.forEach(btn => {
        btn.addEventListener('click', () => {
            const code = btn.nextElementSibling.innerText;
            navigator.clipboard.writeText(code).then(() => {
                const originalText = btn.innerText;
                btn.innerText = 'Copied!';
                btn.style.background = 'var(--teal)';
                btn.style.color = '#fff';
                
                setTimeout(() => {
                    btn.innerText = originalText;
                    btn.style.background = '';
                    btn.style.color = '';
                }, 2000);
            });
        });
    });
}

/**
 * Simple Scroll Reveal Animations
 */
function initScrollAnimations() {
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);

    // Apply reveal styles & observe
    const reveals = document.querySelectorAll('.feature-panel, .community-card, .terminal-window');
    reveals.forEach(el => {
        el.style.opacity = '0';
        el.style.transform = 'translateY(20px)';
        el.style.transition = 'all 0.6s cubic-bezier(0.16, 1, 0.3, 1)';
        
        // Helper class for transition
        const style = document.createElement('style');
        style.innerHTML = `
            .visible {
                opacity: 1 !important;
                transform: translateY(0) !important;
            }
        `;
        document.head.appendChild(style);
        
        observer.observe(el);
    });
}

/**
 * Stats Chart Animation
 */
function initStatsChart() {
    const path = document.querySelector('.chart-line');
    if (!path) return;
    
    const length = path.getTotalLength();
    path.style.strokeDasharray = length;
    path.style.strokeDashoffset = length;
    
    const observer = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting) {
            path.style.transition = 'stroke-dashoffset 2s ease-out';
            path.style.strokeDashoffset = '0';
        }
    });
    
    observer.observe(path.parentElement);
}

/**
 * Mock Contributor Avatars (GitHub Integration Placeholder)
 */
async function initContributors() {
    const slots = document.querySelectorAll('.contributor-slot');
    const repo = 'Mr-Dark-debug/termnode';
    
    try {
        const response = await fetch(`https://api.github.com/repos/${repo}/contributors?per_page=5`);
        if (!response.ok) throw new Error('Failed to fetch');
        const data = await response.json();
        
        slots.forEach((slot, i) => {
            if (data[i]) {
                slot.style.backgroundImage = `url(${data[i].avatar_url})`;
                slot.style.backgroundSize = 'cover';
                slot.title = data[i].login;
                slot.style.cursor = 'pointer';
                slot.onclick = () => window.open(data[i].html_url, '_blank');
            }
        });
    } catch (err) {
        console.warn('Could not fetch contributors:', err);
        // Fallback to placeholders
        slots.forEach((slot, i) => {
            slot.style.backgroundColor = ['#E6FFFA', '#EBF8FF', '#FAF5FF'][i % 3];
            slot.innerText = ['JD', 'AS', 'MK'][i % 3] || '??';
            slot.style.display = 'flex';
            slot.style.alignItems = 'center';
            slot.style.justifyContent = 'center';
            slot.style.fontSize = '0.7rem';
            slot.style.fontWeight = '700';
            slot.style.color = 'var(--teal)';
        });
    }
}
