export class Config {
  private userId: string;
  private accessToken: string;
  private roomId: string;
  private baseUrl: string;

  public constructor() {
    this.userId = process.env.USER_ID ? process.env.USER_ID : "";
    this.accessToken = process.env.ACCESS_TOKEN ? process.env.ACCESS_TOKEN : "";
    this.roomId = process.env.ROOM_ID ? process.env.ROOM_ID : "";
    this.baseUrl = process.env.BASE_URL ? process.env.BASE_URL : "";

    console.log("+++ config +++");
    console.log(process.env.USER_ID);
    console.log(process.env.ACCESS_TOKEN);
    console.log(process.env.ROOM_ID);
    console.log(process.env.BASE_URL);
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
}
