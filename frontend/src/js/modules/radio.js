class RadioController {
  constructor() {
    this.subscribers = [];
  }

  init() {
    this.radioButtons = Array.from(document.querySelectorAll('[data-radio]'));
    this.initBindings();
  }

  initBindings() {
    this.radioButtons.forEach((radioButton) => {
      const input = radioButton.querySelector('input');
      input.addEventListener('change', () => {
        this.subscribers.forEach(subscriber => (subscriber(input.value)));
      });
    });
  }

  subscribe(subscriber) {
    this.subscribers.push(subscriber);
  }
}

export const radioController = new RadioController();
