export class Config {
  private userId: string;
  private accessToken: string;
  private roomId: string;
  private baseUrl: string;
  private menuHtmlUrl: string | null;

  public constructor() {
    this.userId = process.env.USER_ID ? process.env.USER_ID : "";
    this.accessToken = process.env.ACCESS_TOKEN ? process.env.ACCESS_TOKEN : "";
    this.roomId = process.env.ROOM_ID ? process.env.ROOM_ID : "";
    this.baseUrl = process.env.BASE_URL ? process.env.BASE_URL : "";
    this.menuHtmlUrl = process.env.MENU_HTML_URL
      ? process.env.MENU_HTML_URL
      : null;
  }

  public getUserId(): string {
    return this.userId;
  }

  public getAccessToken(): string {
    return this.accessToken;
  }

  public getRoomId(): string {
    return this.roomId;
  }

  public getBaseUrl(): string {
    return this.baseUrl;
  }

  public getMenuUrl(): string | null {
    return this.menuHtmlUrl;
  }
}
