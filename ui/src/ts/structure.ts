export class SubPageNavItem {
  name: string = "";
  routerName: string = "";
  // control if tag is display,
  // when permission of current user >= 'permission', tag will display
  permission: number = 0;
}
