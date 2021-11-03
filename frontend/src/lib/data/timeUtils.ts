const SECONDS = 60;
const MINUTES = SECONDS * 60;
const HOURS = MINUTES * 60;
const DAYS = HOURS * 24;

// how many days/months/years since an event occurred
function timeSince(date: Date): string {
  let duration = new Date().getUTCSeconds() - date.getUTCSeconds();
  let days: number, minutes: number, hours: number, seconds: number;
  duration > DAYS ? (days = duration - DAYS) : 0;
  if (days != 0) duration -= DAYS;
  return '';
}

export default timeSince;
