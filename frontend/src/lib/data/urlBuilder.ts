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

class UrlBuilder {
  baseUrl: string;
  lastRainURL: string;

  constructor() {
    this.baseUrl = ip;
    this.lastRainURL = `${this.baseUrl}/lastRain`;
  }

  // get init args for all cors/json fetch GET requests
  static getInit(): RequestInit {
    const cors: RequestMode = 'cors';
    const args = {
      mode: cors,
      method: 'GET',
      headers: {
        'content-type': 'application/json',
      },
    };
    return args;
  }

  static apiCall(url: string): Promise<any> {
    const data = fetch(url, UrlBuilder.getInit()).then(async (response) => {
      try {
        return response.json();
      } catch (err) {
        console.log(err);
      }
    });
    return data;
  }
}

export default UrlBuilder;
