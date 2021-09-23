export class User {
  userID!: string;
  nickname!: string;
  isLocked!: boolean;
  isLockedDisplay!: string;
  permission!: number;
}

export function displayUserIsLocked(isLocked: boolean): string {
  return isLocked ? "已锁定" : "未锁定";
}
