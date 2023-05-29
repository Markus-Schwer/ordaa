import { appendFile } from "fs";
import { State } from "../states";
import { Command } from "./command";
import { App } from "../app";

export class OrderCommand implements Command {
    public command: string = ".order";

    public process(cmd: string, state: State): void {
        if(state != State.TAKE_ORDERS) {
            App.getInstance().sendMessage("Sorry, I'm not able to do that currently.");
            return;
        }

        App.getInstance().sendMessage("Are you sure that you want to order now?");
    }
}