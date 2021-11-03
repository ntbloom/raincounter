import React from 'react';

type lastRainData = {
  timestamp: Date;
};

interface LastRainProps {
  url: string;
  lastRainMM: number;
  lastRainDuration: number;
}

interface LastRainState {
  date: Date | string;
  //   timeSince: string;
}

class LastRain extends React.Component<LastRainProps, LastRainState> {
  constructor(props: LastRainProps) {
    super(props);
    this.state = {
      date: 'no rain recorded',
      //   timeSince: ""
    };
  }
  componentDidMount() {
    let cors: RequestMode = 'cors';
    const args = {
      mode: cors,
      method: 'GET',
      headers: {
        'content-type': 'application/json',
      },
    };
    const url = this.props.url;
    console.log(url);
    fetch(url, args)
      .then((response) => {
        return response.json();
      })
      .then((data) => {
        console.log(`data=${data}`);
        const val = (data as lastRainData).timestamp;
        this.setState({ date: val });
      })
      .catch((err) => {
        console.error(err);
      });
  }

  render() {
    return (
      <div>
        <p>Last Rain: {this.state.date}</p>
      </div>
    );
  }
}

export { LastRain, LastRainProps };
