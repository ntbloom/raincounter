const MINUTE = 60;
const HOUR = MINUTE * 60;
const DAY = HOUR * 24;

class TimeUtils {
  // gets a human-readable string of time passed since the date
  static timeSince(timestamp: string): string {
    const now = new Date().getTime();
    const dateAsSeconds = Math.floor(Date.parse(timestamp));
    return this.secondsToString((now - dateAsSeconds) / 1000);
  }

  // parse seconds into human-readable string
  static secondsToString(seconds: number): string {
    // just return "hour" for small increments
    if (seconds < HOUR) {
      return `<1 hour`;
    }
    // return "hours" for less than 1 day
    let unit: string;
    if (seconds < DAY) {
      const hours = Math.floor(seconds / HOUR);
      hours == 1 ? (unit = 'hour') : (unit = 'hours');
      return `${hours} ${unit}`;
    }
    // return days for the rest
    const days = Math.floor(seconds / DAY);
    days == 1 ? (unit = 'day') : (unit = 'days');
    return `${days} ${unit}`;
  }

  static getMonth(idx: number) {
    const months = [
      'January',
      'February',
      'March',
      'April',
      'May',
      'June',
      'July',
      'August',
      'September',
      'October',
      'November',
      'December',
    ];
    return months[idx];
  }
}
export default TimeUtils;
