import React from 'react';
import TimeUtils from '../lib/data/timeUtils';
import UrlBuilder from '../lib/data/urlBuilder';

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
    const date = new Date();
    this.state = {
      date: '',
      timeSince: '',
    };
  }

  componentDidMount() {
    const init = UrlBuilder.getInit();
    const url = this.props.url;
    console.log(url);
    fetch(url, init)
      .then((response) => {
        return response.json();
      })
      .then((data) => {
        console.log(`data=${data}`);
        const val = (data as lastRainData).timestamp;
        const since = TimeUtils.timeSince(val);
        const d = new Date(Date.parse(val));
        this.setState({
          date: `${d.toLocaleString('default', {
            month: 'long',
            day: '2-digit',
            year: 'numeric',
          })}`,

          timeSince: since,
        });
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
