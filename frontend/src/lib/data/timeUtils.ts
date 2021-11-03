const MINUTE = 60;
const HOUR = MINUTE * 60;
const DAY = HOUR * 24;

class TimeUtils {
  // how many days/months/years since an event occurred
  static timeSince(date: Date): string {
    let duration = new Date().getUTCSeconds() - date.getUTCSeconds();

    return '';
  }

  static secondsToString(seconds: number): string {
    if (seconds < MINUTE * 2) {
      return '1 minute';
    }
    if (seconds < HOUR) {
      const minutes = Math.floor(seconds / 60);
      return `${minutes} minutes`;
    }
    return '';
  }
}
export default TimeUtils;
