// Configure theme based on system preference
const isDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
const theme = localStorage.getItem("theme") || (isDark ? "dark" : "light");
document.documentElement.setAttribute("data-theme", theme);

// CSRF Token Configuration for htmx requests
document.addEventListener('htmx:configRequest', (event) => {
  const meta = document.querySelector('meta[name="csrf-token"]');
  if (meta) {
    event.detail.headers['X-Csrf-Token'] = meta.content;
  }
});

// Update CSRF Token on htmx request response
document.addEventListener('htmx:beforeOnLoad', (event) => {
  const newToken = event.detail.xhr.getResponseHeader('X-Csrf-Token');
  if (newToken) {
    const meta = document.querySelector('meta[name="csrf-token"]');
    if (meta) meta.content = newToken;
  }
});

window.closeDialog = () => {
  // Resets dialog to a blank state to clear dialog blur effect
  document.getElementById('dialog').outerHTML = '<div id="dialog"></div>';
};

window.clearFieldErrors = () => {
  document.querySelectorAll('[id^="err-"]').forEach(el => el.textContent = '');
};

window.setFieldErrors = (raw) => {
  window.clearFieldErrors();
  try {
    const { errors } = JSON.parse(raw);
    Object.entries(errors).forEach(([field, msg]) => {
      const el = document.getElementById(`err-${field}`);
      if (el) el.textContent = msg;
    });
  } catch (e) {}
};

document.addEventListener('DOMContentLoaded', () => {
  // Theme toggle logic
  const toggle = document.getElementById('theme-toggle');
  const root = document.documentElement;
  if (toggle) {
    toggle.checked = root.getAttribute("data-theme") === "light";
    toggle.onchange = ({ target }) => {
      const theme = target.checked ? "light" : "dark";
      root.setAttribute("data-theme", theme);
      localStorage.setItem("theme", theme);
    };
  }

  // htmx requests logic
  if (window.htmx) {
    // Clear field errors after swap
    htmx.on('htmx:afterSwap', ({ detail }) => {
      if (detail.target.id === 'dialog') window.clearFieldErrors?.();
    });

    // Update visibility of "No Notes" placeholder after note operations
    htmx.on('htmx:afterRequest', ({ detail: { target } }) => {
      const isNoteTarget = target?.id === 'notes' || target?.id?.startsWith('note-');
      const notes = document.getElementById('notes');

      if (isNoteTarget && notes) {
        document.getElementById('no-notes')?.classList.toggle('hidden', notes.childElementCount > 0);
      }
    });
  }
});
