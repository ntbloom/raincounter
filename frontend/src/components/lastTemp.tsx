import React from 'react';

interface LastTempProps {
  url: string;
}

interface LastTempState {
  date: string;
}

interface LastTempData {
  timestamp: string;
}

// LastRain shows the last time it rained and how many hours or days ago it was
class LastRain extends React.Component<LastTempProps, LastTempState> {
  constructor(props: LastTempProps) {
    super(props);
    this.state = {
      date: '',
    };
  }
}

export default LastRain;
