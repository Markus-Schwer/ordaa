import * as sdk from "matrix-js-sdk";
import { Command } from "./commands/command"
import { StartCommand } from "./commands/start-command"
import { State } from "./states";
import { Config } from "./config";

require('dotenv').config();

export class App {
    public config: Config = new Config();

    public App() {
        console.log(this.config.getRoomId());
    }
}

const app: App = new App();


const matrixClient = sdk.createClient({
    baseUrl: app.config.getBaseUrl(),
    accessToken: app.config.getAccessToken(),
    userId: app.config.getUserId(),
});

const commandList: Command[] = [];
commandList.push(new StartCommand());

let currentState = State.IDLE;

/* listen to matrix messages */
matrixClient.on(sdk.RoomEvent.Timeline, (event, room, toStartOfTimeline) => {
    // only listen to messages in the given room which were not sent by this bot
    if (event.getType() !== "m.room.message" || event.event.room_id != app.config.getRoomId() || event.event.sender == app.config.getUserId()) {
        return;
    }

    if(event?.event?.content?.body && event.event.origin_server_ts && (Math.floor(Date.now()) - event.event.origin_server_ts) <= 1000) {
        //sendMessage("[" + event.event.sender + "] " + event.event.content.body);
        processMessage(event.event.content.body);
    }
});

function processMessage(message: string): void {
    const command: Command | undefined = commandList.find(e => e.command === message)

    if(command) {
        command.process(message, currentState);
    }
}

/* send a matrix message */ 
function sendMessage(message: string): void {
    const content: any = {
        body:  message,
        msgtype: "m.text",
    };
    matrixClient.sendEvent(app.config.getRoomId(), "m.room.message", content, "");
}

matrixClient.startClient();
sendMessage(".inder is back!");
