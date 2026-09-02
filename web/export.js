'use strict';

// SVG / PNG / PDF / print exports. Wired by app.js via wireExports(getSVG),
// where getSVG returns the current rendered SVG string.

// US Letter at 300dpi.
const PNG_WIDTH = 2550;
const PNG_HEIGHT = 3300;
const PDF_WIDTH_PT = 612;
const PDF_HEIGHT_PT = 792;

function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

// Rasterizes the SVG with the browser's own SVG engine at print resolution.
async function svgToPngBlob(svg) {
  const svgURL = URL.createObjectURL(new Blob([svg], { type: 'image/svg+xml' }));
  try {
    const img = new Image();
    img.src = svgURL;
    await img.decode();

    const canvas = document.createElement('canvas');
    canvas.width = PNG_WIDTH;
    canvas.height = PNG_HEIGHT;
    canvas.getContext('2d').drawImage(img, 0, 0, PNG_WIDTH, PNG_HEIGHT);

    return await new Promise((resolve, reject) => {
      canvas.toBlob((b) => (b ? resolve(b) : reject(new Error('PNG export failed'))), 'image/png');
    });
  } finally {
    URL.revokeObjectURL(svgURL);
  }
}

async function pngToPdfBlob(pngBlob) {
  const doc = await PDFLib.PDFDocument.create();
  const png = await doc.embedPng(await pngBlob.arrayBuffer());
  const page = doc.addPage([PDF_WIDTH_PT, PDF_HEIGHT_PT]);
  page.drawImage(png, { x: 0, y: 0, width: PDF_WIDTH_PT, height: PDF_HEIGHT_PT });
  const bytes = await doc.save();
  return new Blob([bytes], { type: 'application/pdf' });
}

globalThis.wireExports = (getSVG) => {
  const errorEl = document.getElementById('export-error');
  const busy = (btn, fn) => async () => {
    btn.disabled = true;
    errorEl.textContent = '';
    try {
      await fn();
    } catch (err) {
      errorEl.textContent = 'Export failed: ' + err.message;
    } finally {
      btn.disabled = false;
    }
  };

  const svgBtn = document.getElementById('dl-svg');
  svgBtn.addEventListener('click', busy(svgBtn, async () => {
    downloadBlob(new Blob([getSVG()], { type: 'image/svg+xml' }), 'wifi-sign.svg');
  }));

  const pngBtn = document.getElementById('dl-png');
  pngBtn.addEventListener('click', busy(pngBtn, async () => {
    downloadBlob(await svgToPngBlob(getSVG()), 'wifi-sign.png');
  }));

  const pdfBtn = document.getElementById('dl-pdf');
  pdfBtn.addEventListener('click', busy(pdfBtn, async () => {
    downloadBlob(await pngToPdfBlob(await svgToPngBlob(getSVG())), 'wifi-sign.pdf');
  }));

  document.getElementById('print').addEventListener('click', () => window.print());
};
