import React from 'react';
import DataGetter from '../lib/data/getData';

import { LastRain, LastRainProps } from './lastRain';

const data = new DataGetter();

const main: JSX.Element = (
  <div>
    <LastRain
      url={data.lastRainURL}
      lastRainMM={0}
      lastRainDuration={0}
    ></LastRain>
  </div>
);

export default main;
