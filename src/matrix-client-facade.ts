import { MatrixEvent, RoomEvent } from "matrix-js-sdk";
import * as sdk from "matrix-js-sdk";
import { Config } from "./config";

export class MatrixClientFacade {
  private matrixClient: sdk.MatrixClient;
  private config: Config;

  constructor(config: Config) {
    this.config = config;
    this.matrixClient = sdk.createClient({
      baseUrl: this.config.getBaseUrl(),
      accessToken: this.config.getAccessToken(),
      userId: this.config.getUserId(),
    });
    this.matrixClient.startClient();
  }

  public listenToRoomEvents(callback: (message: string, user: string) => void) {
    this.matrixClient.on(RoomEvent.Timeline, (event: MatrixEvent, _a, _b) => {
      // only listen to messages in the given room which were not sent by this bot
      if (
        event.getType() !== "m.room.message" ||
        event.event.room_id != this.config.getRoomId() ||
        event.event.sender == this.config.getUserId()
      ) {
        return;
      }

      if (
        event?.event?.content?.body &&
        event.event.origin_server_ts &&
        Math.floor(Date.now()) - event.event.origin_server_ts <= 1000
      ) {
        console.log(event.event.sender);
        callback(event.event.content.body, "");
      }
    });
  }

  public sendMessage(message: string): void {
    const content: any = {
      body: message,
      msgtype: "m.text",
    };
    this.matrixClient.sendEvent(
      this.config.getRoomId(),
      "m.room.message",
      content,
      ""
    );
  }
}
