import { State, Transition } from "../states";
import { Command } from "./command";
import { App } from "../app";

export class DeliveredCommand extends Command {
  public readonly transition: Transition = Transition.ARRIVED;
  public readonly command: string = ".delivered";

  public process(): void {
    this.app.sendMessage("@ALL: Food is here!");
    this.app.sendMessage("Bon appetit!");

    // TODO: Abrechnung
  }
}
