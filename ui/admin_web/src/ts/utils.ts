import { Task } from "shared/ts/data";

export function displayPreTasks(preTaskIDs: Array<string>, tasks: Array<Task>): string {
  if (!preTaskIDs || !tasks) { // if one param is null
    return "";
  }

  let taskMap: Map<string, string> = new Map(); // task id - task name
  for (let i = 0; i < tasks.length; i++) {
    taskMap = taskMap.set(tasks[i].taskID, tasks[i].taskName);
  }

  let preTasks = "";
  for (let i = 0; i < preTaskIDs.length; i++) {
    const taskName = taskMap.get(preTaskIDs[i]);

    if (taskName) {
      preTasks += " " + taskName + `\n`;
    } else {
      preTasks += " 未知任务ID：" + preTaskIDs[i] + `\n`;
    }
  }

  if (preTasks.endsWith("\n")) {
    preTasks = preTasks.slice(0, preTasks.length - 1);
  }

  return trimSpace(preTasks);
}

export function sortTasks(tasks: Array<Task>): Array<Task> {
  const res: Array<Task> = new Array<Task>();

  // record tasks without pre-tasks
  let length = tasks.length
  for (let i = 0; i < length; i++) {
    if (tasks[i].preTaskIDs.length > 0) {
      continue;
    }

    res.push(tasks[i]);

    tasks[i] = tasks[0];
    tasks.shift();
    length = tasks.length;
  }

  // record tasks which pre-tasks are all in 'res', loop
  while (length > 0) {
    const lengthBackup = length;

    for (let i = 0; i < length; i++) {
      if (!contains(res, tasks[i].preTaskIDs)) {
        continue;
      }
      
      res.push(tasks[i]);
      
      tasks[i] = tasks[0];
      tasks.shift();
      length = tasks.length;
    }

    if (lengthBackup == length) {
      break;
    }
  }

  // handle tasks with unexpected pre-tasks
  for (let i = 0; i < length; i++) {
    res.push(tasks[i]);
  }

  return res;
}

function trimSpace(str: string): string {
  let i = 0;
  while (i < str.length && str[i] == ' ') {
    i++;
  }

  str = str.slice(i);

  i = str.length - 1
  while (i >= 0 && str[i] == ' ') {
    i--
  }

  str = str.slice(0, i + 1);

  return str
}

// contains return if sliceA contains all items in sliceB
function contains(sliceA: Array<Task>, sliceB: Array<string>): boolean {
  if (sliceA.length < sliceB.length) {
    return false;
  }
  
  let isContained = true;
  for (let i = 0; i < sliceB.length; i++) {
    let isContainedInner = false;
    for (let j = 0; j < sliceA.length; j++) {
      if (sliceA[j].taskID == sliceB[i]) {
        isContainedInner = true;
        break;
      }
    }
    
    if (!isContainedInner) {
     isContained = false;
     break;
    }
  }
  
  return isContained
}
