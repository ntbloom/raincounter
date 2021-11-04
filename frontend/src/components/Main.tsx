import React from 'react';
import UrlBuilder from '../lib/data/urlBuilder';

import { LastRain, LastRainProps } from './lastRain';

const data = new UrlBuilder();

const main: JSX.Element = (
  <div>
    <LastRain
      url={data.lastRainURL}
      //  // come back to these later?
      // lastRainMM={0}
      // lastRainDuration={0}
    ></LastRain>
  </div>
);

export default main;
