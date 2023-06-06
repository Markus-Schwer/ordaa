import * as sdk from "matrix-js-sdk";
import { Command } from "./commands/command";
import { StartCommand } from "./commands/start-command";
import { State } from "./states";
import { Config } from "./config";
import { DeliveredCommand } from "./commands/delivered-command";
import { OrderCommand } from "./commands/order-command";
import { HelpCommand } from "./commands/help-command";

export class App {
  private config: Config = new Config();

  // TODO: use a map here for efficency
  private commandList: Command[] = [];

  private state: State = State.IDLE;

  private matrixClient: sdk.MatrixClient = sdk.createClient({
    baseUrl: this.config.getBaseUrl(),
    accessToken: this.config.getAccessToken(),
    userId: this.config.getUserId(),
  });

  public constructor() {
    this.matrixClient.startClient();
    this.sendMessage(".inder is back!");

    console.log(this.config.getRoomId());

    // register your commands here
    this.commandList.push(new StartCommand(this));
    this.commandList.push(new DeliveredCommand(this));
    this.commandList.push(new OrderCommand(this));
    this.commandList.push(new HelpCommand(this));

    /* listen to matrix messages */
    this.matrixClient.on(
      sdk.RoomEvent.Timeline,
      (event, room, toStartOfTimeline) => {
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
          //sendMessage("[" + event.event.sender + "] " + event.event.content.body);
          this.processMessage(event.event.content.body);
        }
      }
    );
  }

  private processMessage(message: string): void {
    const command: Command | undefined = this.commandList.find(
      (e) => e.command === message
    );

    if (command) {
      command.process();
    }
  }

  /* send a matrix message */
  public sendMessage(message: string): void {
    const app = this; //App.getInstance();

    const content: any = {
      body: message,
      msgtype: "m.text",
    };
    app.matrixClient.sendEvent(
      app.config.getRoomId(),
      "m.room.message",
      content,
      ""
    );
  }

  public setState(newState: State): void {
    this.state = newState;
  }

  public getState(): State {
    return this.state;
  }
}

require("dotenv").config();
new App();
