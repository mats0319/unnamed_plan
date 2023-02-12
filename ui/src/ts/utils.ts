export function displayIsPublic(isPublic: boolean): string {
  return isPublic ? "公开" : "非公开";
}

/**
 * display time string base on timestamp, format: 'yyyy-mm-dd hh:mm:ss'
 * @param time timestamp, unit: second
 */
export function displayTime(time: number): string {
  const date: Date = new Date(time * 1000);

  const year = date.getFullYear();
  const month = formatTimeDigit(date.getMonth() + 1);
  const day = formatTimeDigit(date.getDate());

  const hour = formatTimeDigit(date.getHours());
  const minute = formatTimeDigit(date.getMinutes());
  const second = formatTimeDigit(date.getSeconds());

  return year + "-" + month + "-" + day + " " + hour + ":" + minute + ":" + second;
}

function formatTimeDigit(time: number): string {
  return time >= 10 ? time.toString() : "0" + time;
}

export function timestampMSToS(timestamp: number): number {
  return parseInt((timestamp / 1000).toString())
}
