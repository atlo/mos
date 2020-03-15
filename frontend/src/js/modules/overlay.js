const OVERLAY_OPEN = 'overlay--open';

class Overlay {
  constructor(options) {
    this.overlay = options.overlay;
    this.key = null;

    if (this.overlay) {
      this.key = this.overlay.getAttribute('data-overlay');
      this.init();
    }
  }

  init() {
    this.initDOMElements();
    this.initBindings();
  }

  initDOMElements() {
    this.overlayClose = this.overlay.querySelector('[data-overlay-close]');
    this.overlayGraph = this.overlay.querySelector('[data-overlay-graph]');
    this.overlayScrollUp = this.overlay.querySelector('[data-overlay-scroll="up"]');
    this.overlayScrollDown = this.overlay.querySelector('[data-overlay-scroll="down"]');
  }

  initBindings() {
    this.overlayClose.addEventListener('click', this.close.bind(this));

    if (this.overlayGraph) {
      this.overlayScrollDown.addEventListener('click', () => {
        this.overlayGraph.scrollTop = this.overlayGraph.scrollTop + 40;
      });

      this.overlayScrollUp.addEventListener('click', () => {
        this.overlayGraph.scrollTop = this.overlayGraph.scrollTop - 40;
      });
    }
  }

  open() {
    this.overlay.classList.add(OVERLAY_OPEN);
  }

  close() {
    this.overlay.classList.remove(OVERLAY_OPEN);
  }
}

class OverlayController {
  constructor() {
    this.overlays = {};
  }

  init() {
    this.initDOMElements();
    this.collectOverlays();
    this.initBindings();
  }

  initDOMElements() {
    this.overlayElements = Array.from(document.querySelectorAll('[data-overlay]'));
    this.overlayOpenButtons = Array.from(document.querySelectorAll('[data-overlay-open]'));
  }

  collectOverlays() {
    this.overlayElements.forEach((overlay) => {
      const key = overlay.getAttribute('data-overlay');
      this.overlays[key] = new Overlay({ overlay });
    });
  }

  initBindings() {
    this.overlayOpenButtons.forEach((overlayButton) => {
      overlayButton.addEventListener('click', this.openOverlay.bind(this));
    });
  }

  openOverlay(e) {
    e.preventDefault();
    const key = e.target.getAttribute('data-overlay-open');
    this.overlays[key].open();
  }
}

export const overlayController = new OverlayController();
