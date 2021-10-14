import GetData from '../src/lib/getData';
import fetch from 'node-fetch';

describe('GetData', () => {
     let g = new GetData();
     it('get last rain returns data', async () => {
          const data = await g.getLastRain();
          expect(data).toBe(1234);
     });
});
