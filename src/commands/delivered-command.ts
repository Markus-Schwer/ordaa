import { Transition } from "../states";
import { Command } from "./command";

export class DeliveredCommand extends Command {
    public readonly transition: Transition = Transition.ARRIVED;
    public readonly matcher = new RegExp("^\\.delivered$");

    public process(): void {
        this.app.sendMessage("@ALL: Food is here!");
        this.app.sendMessage("Bon appetit!");

        // TODO: Abrechnung
    }
}
