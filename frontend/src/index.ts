import GetData from './lib/getData';

async function lastRain() {
     const g: GetData = new GetData();
     const rain = await g.getLastRain();
     console.log(rain);
}
