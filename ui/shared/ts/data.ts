export interface User {
  userID: string;
  userName: string;
  nickname: string;
  isLocked: boolean;
  permission: number;
  createdBy: string;
}

export interface CloudFile {
  fileID: string;
  fileName: string;
  fileURL: string;
  isPublic: boolean;
  updateTime: number;
  createdTime: number;
}

export interface ThinkingNote {
  noteID: string;
  writeBy: string;
  topic: string;
  content: string;
  isPublic: boolean;
  updateTime: number;
  createdTime: number;
}

const cloudFileURLPrefix = "https://mats9693.cn/cloud-file/";

export function generateCloudFileURL(url: string): string {
  return cloudFileURLPrefix + url
}

export function displayIsLocked(isLocked: boolean): string {
  return isLocked ? "已锁定" : "未锁定";
}

export function displayIsPublic(isPublic: boolean): string {
  return isPublic ? "公开" : "非公开";
}

/**
 * display time string base on timestamp
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
