import { Transition } from "../states";
import { Command } from "./command";

export class OrderCommand extends Command {
  public readonly matcher: RegExp = new RegExp("^!order (\\w{0,1}\\d+)$");
  public readonly transition: Transition = Transition.ADD_ITEM;

  public process(): void {
    /*if (state != State.TAKE_ORDERS) {
      this.app.sendMessage("Sorry, I'm not able to do that currently.");
      return;
    }*/

    this.app.sendMessage("Are you sure that you want to order now?");
  }
}
