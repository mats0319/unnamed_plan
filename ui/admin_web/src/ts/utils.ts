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
