import ip from './ipaddress';

// get the payload from the API for a url
async function getUrl(url: string): Promise<object> {
     const args = {
          method: 'GET',
          headers: {
               'content-type': 'application/json',
          },
     };
     const val = await fetch(url, args).then((response: Response) => {
          if (!response.ok) {
               throw new Error('bad response');
          }
          return response.json();
     });
     return val;
}

class GetData {
     baseUrl: string;

     constructor() {
          this.baseUrl = ip;
     }

     // get last rain value
     getLastRain(): object {
          const url = `${this.baseUrl}/lastRain`;
          return getUrl(url);
     }
}

export default GetData;
