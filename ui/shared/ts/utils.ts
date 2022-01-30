// @ts-ignore
import sha256 from "crypto-js/sha256";
import { taskStatus } from "./const";

export function calcSHA256(message: string): string {
  return sha256(message);
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

export function displayTaskStatus(status: number): string {
  let res = "未知状态："+status;
  if (0 <= status && status < taskStatus.length) {
    res = taskStatus[status];
  }

  return res;
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

/**
 *  compareOnStringSliceNotStrict compare two string-slice, return if they contain same values
 *
 *  tips:
 *  1. ignore data orders
 *  2. can handle duplicated data
 */
export function compareOnStringSliceNotStrict(a: Array<string>, b: Array<string>): boolean {
  if (a.length != b.length) {
    return false;
  }

  let aMap: Map<string, number> = new Map();
  for (let i = 0; i < a.length; i++) {
    let v = aMap.get(a[i]); // v: undefined if key is not exist
    aMap = aMap.set(a[i], v ? v+1 : 1);
  }

  let isEqual = true;
  for (let i = 0; i < b.length; i++) {
    let v = aMap.get(b[i]);
    if (!v) {
      isEqual = false;
      break;
    }

    if (v > 1) {
      aMap = aMap.set(b[i], v-1);
    } else {
      aMap.delete(b[i]);
    }
  }

  return isEqual;
}
