import React from 'react';
import TimeUtils from '../lib/data/timeUtils';

type lastRainData = {
  timestamp: string;
};

interface LastRainProps {
  url: string;
  lastRainMM: number;
  lastRainDuration: number;
}

interface LastRainState {
  date: string;
  timeSince: string;
}

class LastRain extends React.Component<LastRainProps, LastRainState> {
  constructor(props: LastRainProps) {
    super(props);
    this.state = {
      date: 'no rain recorded',
      timeSince: '',
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
        const since = TimeUtils.timeSince(val);
        this.setState({ date: val, timeSince: since });
      })
      .catch((err) => {
        console.error(err);
      });
  }

  render() {
    return (
      <div>
        <p>
          Last Rain Event: {this.state.date} ({this.state.timeSince} ago)
        </p>
      </div>
    );
  }
}

export { LastRain, LastRainProps };
