export interface User {
  userID: string;
  userName: string;
  nickname: string;

  isLocked: boolean;
  isLockedDisplay: string;

  permission: number;
  createdBy: string;
}

export interface CloudFile {
  fileID: string;
  fileName: string;
  fileURL: string;

  isPublic: boolean;
  isPublicDisplay: string;

  updateTime: number;
  updateTimeDisplay: string;

  createdTime: number;
  createdTimeDisplay: string;
}

export interface ThinkingNote {
  noteID: string;
  writeBy: string;
  topic: string;
  content: string;

  isPublic: boolean;
  isPublicDisplay: string;

  updateTime: number;
  updateTimeDisplay: string;

  createdTime: number;
  createdTimeDisplay: string;

}

const cloudFileURLOrigin = "https://mats9693.cn/cloud-file/";

export function generateCloudFileURL(url: string): string {
  return cloudFileURLOrigin + url
}

export function displayIsPublic(isPublic: boolean): string {
  return isPublic ? "公开" : "非公开";
}

export function displayTime(time: number): string {
  return (new Date(time*1000)).toString();
}
