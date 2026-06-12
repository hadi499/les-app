import { browser } from '$app/environment';

let imageZoom = $state(1.0);
let headerSize = $state(1.0);
let titleSize = $state(1.0);
let contentSize = $state(1.0);

let cardsPerPage = $state(6);

if (browser) {
  const saved = localStorage.getItem('print-settings');
  if (saved) {
    const s = JSON.parse(saved);
    imageZoom = s.imageZoom ?? 1.0;
    headerSize = s.headerSize ?? 1.0;
    titleSize = s.titleSize ?? 1.0;
    contentSize = s.contentSize ?? 1.0;
    cardsPerPage = s.cardsPerPage ?? 6;
  }
}

function save() {
  if (!browser) return;
  localStorage.setItem('print-settings', JSON.stringify({ imageZoom, headerSize, titleSize, contentSize, cardsPerPage }));
}

function stepUp(v: number) { return Math.min(2.0, +(v + 0.1).toFixed(1)); }
function stepDown(v: number) { return Math.max(0.5, +(v - 0.1).toFixed(1)); }

export const printSettings = {
  get imageZoom() { return imageZoom; },
  get headerSize() { return headerSize; },
  get titleSize() { return titleSize; },
  get contentSize() { return contentSize; },
  get cardsPerPage() { return cardsPerPage; },

  setCardsPerPage(n: number) { cardsPerPage = n; save(); },

  zoomIn() { imageZoom = stepUp(imageZoom); save(); },
  zoomOut() { imageZoom = stepDown(imageZoom); save(); },

  headerUp() { headerSize = stepUp(headerSize); save(); },
  headerDown() { headerSize = stepDown(headerSize); save(); },
  titleUp() { titleSize = stepUp(titleSize); save(); },
  titleDown() { titleSize = stepDown(titleSize); save(); },
  contentUp() { contentSize = stepUp(contentSize); save(); },
  contentDown() { contentSize = stepDown(contentSize); save(); },

  reset() {
    imageZoom = 1.0;
    headerSize = 1.0;
    titleSize = 1.0;
    contentSize = 1.0;
    cardsPerPage = 6;
    save();
  }
};
