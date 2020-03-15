/* eslint-disable */
const sliderRefs = [];

class Slider {
    constructor(options) {
        this.slider = options.slider;

        if (this.slider) {
            this._id = this.slider.getAttribute('data-slider');
            this.sliderControl = this.slider.querySelector('[data-slider-control]');
            this.isMoving = false;
            this.width = this.slider.getBoundingClientRect().width;
            this.prevPos = 0;
            this.prevX = 0;
            this.index = 0;
            this.subscribers = [];
            this.currentYear = 1998;

            this.onMouseDown = this.onMouseDown.bind(this);
            this.onMouseMove = this.onMouseMove.bind(this);
            this.onMouseUp = this.onMouseUp.bind(this);
            this.onResize = this.onResize.bind(this);

            this.init();
        }
    }

    get id() {
      return this._id;
    }

    init() {
        this.initBindings();
    }

    initBindings() {
        this.slider.addEventListener('mousedown', this.onMouseDown);
        window.addEventListener('mousemove', this.onMouseMove);
        window.addEventListener('mouseup', this.onMouseUp);
    }

    onMouseDown(e) {
        this.isMoving = true;
        this.prevX = e.clientX;
    }

    onMouseMove(e) {
        if (!this.isMoving) return;

        const absoluteWidth = this.width - 2 * 15;
        const itemNum = this.slider.querySelectorAll('.slider__item').length + 1;
        const step = absoluteWidth / (itemNum - 1);
        const currX = e.clientX;

        const diff = currX - this.prevX;
        const dir = diff > 0 ? 1 : -1;
        const nextX = this.prevPos + step * dir;
        
        if (Math.abs(diff) > step && nextX < step * itemNum && nextX >= 0) {
            this.sliderControl.setAttribute('style', `transform: translate(${this.prevPos + step * dir}px, -50%)`);
            this.sliderControl.innerHTML = `<div class="slider__control-year">${this.currentYear + dir}</div>`;
            this.index += dir;
            this.prevPos += step * dir;
            this.prevX = currX;
            this.currentYear += dir;
        }
    }

    setYear(year) {
        const absoluteWidth = this.width - 2 * 15;
        const itemNum = this.slider.querySelectorAll('.slider__item').length + 1;
        const step = absoluteWidth / (itemNum - 1);

        this.sliderControl.setAttribute('style', `transform: translate(${step * (year - 1998)}px, -50%)`);
        this.sliderControl.innerHTML = `<div class="slider__control-year">${year}</div>`;
        this.index = (year - 1998);
        this.prevPos = step * (year - 1998);
        this.currentYear = year;
    }

    subscribe(subscriber) {
        this.subscribers.push(subscriber);
    }

    notify() {
        console.log(this.subscribers);
        this.subscribers.forEach(subscriber => {
            subscriber(this.currentYear);
        });
    }

    onMouseUp() {
        if (!this.isMoving) return;

        this.isMoving = false;
        this.notify();
    }

    onResize() {
        this.width = this.slider.getBoundingClientRect().width;
    }
}

export const initSliders = () => {
    const sliderNodes = document.querySelectorAll('[data-slider]');
    console.log(sliderNodes);
    Array.from(sliderNodes).forEach(slider => (
        sliderRefs.push(new Slider({ slider }))
    ));
}

export const sliders = sliderRefs;
