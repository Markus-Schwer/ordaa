import { State } from "../states";
import { Command } from "./command";

export class StartCommand implements Command {
    public command: string = ".inder";

    public process(cmd: string, state: State): void {
        console.error("Start command executed!");
    }
}