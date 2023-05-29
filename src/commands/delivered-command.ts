import { State } from "../states";
import { Command } from "./command";
import { App } from "../app";

export class DeliveredCommand implements Command {
    public command: string = ".delivered";

    public process(cmd: string, state: State): void {
        App.getInstance().sendMessage("@ALL: Food is here!");
        App.getInstance().sendMessage("Bon appetit!");

        // TODO: Abrechnung
    }
}