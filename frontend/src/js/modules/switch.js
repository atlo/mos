class SwitchController {
  constructor() {
    this.switches = [];
    this.subscribers = [];
  }

  add(value) {
    this.switches.push(value);
  }

  init() {
    this.switches.forEach((item) => {
      item.input.addEventListener('change', () => {
        this.subscribers.forEach(subscriber => (subscriber(item.input.value, item.input.checked)));
      });
    });
  }

  subscribe(subscriber) {
    this.subscribers.push(subscriber);
  }
}

class Switch {
  constructor(options) {
    this.switch = options.switch;
    this.input = null;
    this.label = null;

    this.initBindings();
  }

  initBindings() {
    this.input = this.switch.querySelector('input');
    this.label = this.switch.querySelector('label');

    this.label.addEventListener('click', () => {
      if (this.input.checked) {
        setTimeout(() => {
          this.input.checked = false;
          const event = new Event('change');
          this.input.dispatchEvent(event);
        }, 10);
      }
    });
  }
}

export const switchController = new SwitchController();

export const initSwitches = () => {
  Array.from(document.querySelectorAll('[data-switch]')).forEach((val) => {
    switchController.add(new Switch({ switch: val }));
  });
  switchController.init();
};
