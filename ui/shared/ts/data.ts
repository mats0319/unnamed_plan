export interface User {
  userID: string;
  userName: string;
  nickname: string;
  isLocked: boolean;
  permission: number;
}

export function newUser(): User {
  return {
    userID: "",
    userName: "",
    nickname: "",
    isLocked: false,
    permission: 0,
  };
}

export function deepCopyUser(data: User): User {
  return {
    userID: data.userID,
    userName: data.userName,
    nickname: data.nickname,
    isLocked: data.isLocked,
    permission: data.permission,
  };
}

export interface CloudFile {
  fileID: string;
  fileName: string;
  fileURL: string;
  isPublic: boolean;
  updateTime: number;
  createdTime: number;
}

export function newCloudFile(): CloudFile {
  return {
    fileID: "",
    fileName: "",
    fileURL: "",
    isPublic: false,
    updateTime: 0,
    createdTime: 0,
  }
}

export function deepCopyCloudFile(data: CloudFile): CloudFile {
  return {
    fileID: data.fileID,
    fileName: data.fileName,
    fileURL: data.fileURL,
    isPublic: data.isPublic,
    updateTime: data.updateTime,
    createdTime: data.createdTime,
  }
}

export interface Note {
  noteID: string;
  writeBy: string;
  topic: string;
  content: string;
  isPublic: boolean;
  updateTime: number;
  createdTime: number;
}

export function newNote(): Note {
  return {
    noteID: "",
    writeBy: "",
    topic: "",
    content: "",
    isPublic: false,
    updateTime: 0,
    createdTime: 0,
  }
}

export function deepCopyNote(data: Note): Note {
  return {
    noteID: data.noteID,
    writeBy: data.writeBy,
    topic: data.topic,
    content: data.content,
    isPublic: data.isPublic,
    updateTime: data.updateTime,
    createdTime: data.createdTime,
  }
}

export interface Task {
  taskID: string;
  taskName: string;
  description: string;
  preTaskIDs: Array<string>;
  preTasks: string; // for display
  status: number;
  updateTime: number;
  createdTime: number;
}

export function newTask(): Task {
  return {
    taskID: "",
    taskName: "",
    description: "",
    preTaskIDs: new Array<string>(),
    preTasks: "",
    status: 0,
    updateTime: 0,
    createdTime: 0,
  }
}

export function deepCopyTask(data: Task): Task {
  return {
    taskID: data.taskID,
    taskName: data.taskName,
    description: data.description,
    preTaskIDs: data.preTaskIDs,
    preTasks: data.preTasks,
    status: data.status,
    updateTime: data.updateTime,
    createdTime: data.createdTime,
  }
}
