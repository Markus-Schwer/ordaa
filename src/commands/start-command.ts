import { State, Transition } from "../states";
import { Command } from "./command";

export class StartCommand extends Command {
  public transition: Transition = Transition.START_ORDER;
  public matcher = new RegExp("^\\.inder$");

  public process(): void {
    /*if (state != State.IDLE) {
      this.app.sendMessage("Hey! I'm already running.");
      return;
    }*/

    this.app.sendMessage(
      "Hey, dotinder here. I'm able to take your orders now."
    );
    this.app.setState(State.TAKE_ORDERS);
  }
}
