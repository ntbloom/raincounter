import ip from './ipAddress';

// get the payload from the API for a url
async function getUrl(url: string): Promise<Response> {
  let cors: RequestMode = 'cors';
  const args = {
    mode: cors,
    method: 'GET',
    headers: {
      'content-type': 'application/json',
    },
  };
  const val = await fetch(url, args).then((res: Response) => {
    if (!res.ok) {
      throw new Error(`bad response, error code ${res.status}`);
    }
    return res.json();
  });
  return val;
}

class DataGetter {
  baseUrl: string;
  lastRainURL: string;

  constructor() {
    this.baseUrl = ip;
    this.lastRainURL = `${this.baseUrl}/lastRain`;
  }

  // get last rain value
  async getLastRain(): Promise<object> {
    const url = `${this.baseUrl}/lastRain`;
    console.log(`url=${url}`);
    const data = await getUrl(url);
    console.log(data);
    return data;
  }
}

export default DataGetter;
