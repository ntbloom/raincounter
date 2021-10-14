import GetData from '../src/lib/getData';

describe('GetData', () => {
     it('get last rain returns data', async () => {
          let g = new GetData();
          const data = g.getLastRain();
          console.log(data);
          expect(data).toBeInstanceOf('string');
     });
});
