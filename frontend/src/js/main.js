import BlendSVG from 'svg-shape-blend/src/index';
import '../scss/main.scss';
import './helpers/vendor-import';
import Graph from './modules/graph';
import { initSliders } from './modules/slider';
import { initSwitches } from './modules/switch';
import { radioController } from './modules/radio';
import { overlayController } from './modules/overlay';

/* eslint-disable no-new */
class App {
  constructor() {
    try {
      initSliders();
      initSwitches();
      overlayController.init();
      radioController.init();
      new Graph();
      new BlendSVG();
    } catch (err) {
      /* eslint-disable no-console */
      console.log(err);
    }
  }
}

window.onload = () => (new App());
