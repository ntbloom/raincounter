import React from 'react';
import TimeUtils from '../lib/data/timeUtils';
import UrlBuilder from '../lib/data/urlBuilder';

interface LastRainProps {
  url: string;
  //   // come back to these later?
  //   lastRainMM: number;
  //   lastRainDuration: number;
}

interface LastRainState {
  date: string;
  timeSince: string;
  rawStamp: string;
}

interface LastRainData {
  timestamp: string;
}

// LastRain shows the last time it rained and how many hours or days ago it was
class LastRain extends React.Component<LastRainProps, LastRainState> {
  constructor(props: LastRainProps) {
    super(props);
    this.state = {
      date: '',
      timeSince: '',
      rawStamp: '',
    };
  }

  // get last rain from API
  apiCall() {
    const data = UrlBuilder.apiCall(this.props.url).then((data) => {
      const timestamp = (data as LastRainData).timestamp;
      console.log(`data=${data}, timestamp=${timestamp}`);
      this.setState({
        rawStamp: timestamp,
        date: TimeUtils.getMonthDayYear(timestamp),
        timeSince: `(${TimeUtils.getTimeSince(timestamp)} ago)`,
      });
    });
  }

  componentDidMount() {
    this.apiCall();
  }

  componentDidUpdate(_: LastRainProps, prevState: LastRainState) {
    if (prevState.rawStamp !== this.state.rawStamp) {
      this.apiCall();
    }
  }

  render() {
    return (
      <p id="lastRain">
        Last Rain Event: {this.state.date} {this.state.timeSince}
      </p>
    );
  }
}

export { LastRain, LastRainProps };
