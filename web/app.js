'use strict';

// Form → params JSON → wifiSign.render (WASM) → live preview.
// Export wiring lives in export.js.

const $ = (id) => document.getElementById(id);

let logoDataURI = '';
let currentSVG = '';

function collectParams() {
  return JSON.stringify({
    ssid: $('ssid').value,
    password: $('password').value,
    auth: $('auth').value,
    hidden: $('hidden').checked,
    accentColor: $('accentColor').value,
    backgroundColor: $('backgroundColor').value,
    logoDataUri: logoDataURI,
    ornament: $('ornament').value,
    tagline: $('tagline').value,
    headline: $('headline').value,
    subtitle: $('subtitle').value,
    footerText: $('footerText').value,
  });
}

function clearFieldErrors() {
  document.querySelectorAll('.field-error').forEach((el) => (el.textContent = ''));
}

function showFieldErrors(errors) {
  for (const err of errors) {
    const el = document.querySelector(`.field-error[data-field="${err.field}"]`);
    if (el) el.textContent = err.message;
  }
}

function setDownloadsEnabled(enabled) {
  for (const id of ['dl-svg', 'dl-png', 'dl-pdf', 'print']) {
    $(id).disabled = !enabled;
  }
}

function refresh() {
  clearFieldErrors();

  if ($('ssid').value === '') {
    // Nothing to render yet; don't nag about a required field being empty.
    $('preview').hidden = true;
    $('preview-hint').hidden = false;
    setDownloadsEnabled(false);
    return;
  }

  const result = JSON.parse(globalThis.wifiSign.render(collectParams()));
  if (result.errors) {
    showFieldErrors(result.errors);
    setDownloadsEnabled(false);
    return;
  }

  currentSVG = result.svg;
  $('preview').innerHTML = result.svg; // self-generated; user text is XML-escaped in Go
  $('preview').hidden = false;
  $('preview-hint').hidden = true;
  setDownloadsEnabled(true);
}

function debounce(fn, ms) {
  let t;
  return () => {
    clearTimeout(t);
    t = setTimeout(fn, ms);
  };
}

const refreshSoon = debounce(refresh, 150);

function readLogo(file) {
  if (!file) return;
  const reader = new FileReader();
  reader.onload = () => {
    logoDataURI = reader.result;
    $('logo-clear').hidden = false;
    refresh();
  };
  reader.readAsDataURL(file);
}

function wireForm() {
  document.querySelectorAll('#sign-form input, #sign-form select').forEach((el) => {
    if (el.type === 'file') return;
    el.addEventListener('input', refreshSoon);
  });
  $('ornament').addEventListener('input', () => {
    $('tagline-label').hidden = $('ornament').value !== 'tagline';
  });
  $('logo').addEventListener('change', (e) => readLogo(e.target.files[0]));
  $('logo-clear').addEventListener('click', () => {
    logoDataURI = '';
    $('logo').value = '';
    $('logo-clear').hidden = true;
    refresh();
  });
  if (globalThis.wireExports) globalThis.wireExports(() => currentSVG);
}

// WASM bootstrap. wifiSignReady is invoked by the Go side once the API exists.
globalThis.wifiSignReady = () => {
  $('loading').hidden = true;
  $('version').textContent = 'v' + globalThis.wifiSign.version;
  wireForm();
  refresh();
};

(async () => {
  try {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(
      fetch('main.wasm'),
      go.importObject
    );
    go.run(result.instance); // resolves wifiSignReady via the Go side
  } catch (err) {
    $('loading').textContent = 'Failed to load the generator: ' + err;
  }
})();
