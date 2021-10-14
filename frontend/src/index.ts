import GetData from './lib/getData';

test('get last rain returns data', () => {
     let g = new GetData();
     return g.getLastRain().then((data) => {
          console.log(data);
     });
});
