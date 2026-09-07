import { browser } from '$app/environment';

let imageZoom = $state(1.0);
let headerSize = $state(1.0);
let titleSize = $state(1.0);
let contentSize = $state(1.0);
let marginLeft = $state(0.0);
let marginRight = $state(0.0);

let cardsPerPage = $state(6);

if (browser) {
  const saved = localStorage.getItem('print-settings');
  if (saved) {
    const s = JSON.parse(saved);
    imageZoom = s.imageZoom ?? 1.0;
    headerSize = s.headerSize ?? 1.0;
    titleSize = s.titleSize ?? 1.0;
    contentSize = s.contentSize ?? 1.0;
    marginLeft = s.marginLeft ?? 0.0;
    marginRight = s.marginRight ?? 0.0;
    cardsPerPage = s.cardsPerPage ?? 6;
  }
}

function save() {
  if (!browser) return;
  localStorage.setItem('print-settings', JSON.stringify({ imageZoom, headerSize, titleSize, contentSize, cardsPerPage, marginLeft, marginRight }));
}

function stepUp(v: number) { return Math.min(2.0, +(v + 0.1).toFixed(1)); }
function stepDown(v: number) { return Math.max(0.5, +(v - 0.1).toFixed(1)); }
function stepUpMargin(v: number) { return Math.min(50.0, +(v + 1.0).toFixed(1)); }
function stepDownMargin(v: number) { return Math.max(0.0, +(v - 1.0).toFixed(1)); }

export const printSettings = {
  get imageZoom() { return imageZoom; },
  get headerSize() { return headerSize; },
  get titleSize() { return titleSize; },
  get contentSize() { return contentSize; },
  get cardsPerPage() { return cardsPerPage; },
  get marginLeft() { return marginLeft; },
  get marginRight() { return marginRight; },

  setCardsPerPage(n: number) { cardsPerPage = n; save(); },

  zoomIn() { imageZoom = stepUp(imageZoom); save(); },
  zoomOut() { imageZoom = stepDown(imageZoom); save(); },

  headerUp() { headerSize = stepUp(headerSize); save(); },
  headerDown() { headerSize = stepDown(headerSize); save(); },
  titleUp() { titleSize = stepUp(titleSize); save(); },
  titleDown() { titleSize = stepDown(titleSize); save(); },
  contentUp() { contentSize = stepUp(contentSize); save(); },
  contentDown() { contentSize = stepDown(contentSize); save(); },
  
  marginLeftUp() { marginLeft = stepUpMargin(marginLeft); save(); },
  marginLeftDown() { marginLeft = stepDownMargin(marginLeft); save(); },
  marginRightUp() { marginRight = stepUpMargin(marginRight); save(); },
  marginRightDown() { marginRight = stepDownMargin(marginRight); save(); },

  reset() {
    imageZoom = 1.0;
    headerSize = 1.0;
    titleSize = 1.0;
    contentSize = 1.0;
    marginLeft = 0.0;
    marginRight = 0.0;
    cardsPerPage = 6;
    save();
  }
};
