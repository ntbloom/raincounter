import React from 'react';
import TimeUtils from '../lib/data/timeUtils';
import UrlBuilder from '../lib/data/urlBuilder';

interface LastRainProps {
  url: string;
  lastRainMM: number;
  lastRainDuration: number;
}

interface LastRainState {
  date: string;
  timeSince: string;
}

interface LastRainData {
  timestamp: string;
}

// LastRain shows the last time it rained and how many hours or days ago it was
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
    fetch(this.props.url, UrlBuilder.getInit())
      .then((response) => {
        return response.json();
      })

      .then((data) => {
        const timestamp = (data as LastRainData).timestamp;
        this.setState({
          date: TimeUtils.getMonthDayYear(timestamp),
          timeSince: TimeUtils.getTimeSince(timestamp),
        });
      })

      .catch((err) => {
        console.error(err);
      });
  }

  render() {
    return (
      <p id="lastRain">
        Last Rain Event: {this.state.date} ({this.state.timeSince} ago)
      </p>
    );
  }
}

export { LastRain, LastRainProps };
