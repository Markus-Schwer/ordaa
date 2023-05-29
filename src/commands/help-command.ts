import { State } from "../states";
import { Command } from "./command";
import { App } from "../app";

export class HelpCommand implements Command {
    public command: string = ".help";

    public process(cmd: string, state: State): void {
        App.getInstance().sendMessage("+++ HELP +++");

        // TODO: add help message
    }
}