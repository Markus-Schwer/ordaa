import { App } from "../app";
import { State } from "../states";
import { Command } from "./command";

export class StartCommand implements Command {
  public command: string = ".inder";

  public process(cmd: string, state: State): void {
    if (state != State.IDLE) {
      App.getInstance().sendMessage("Hey! I'm already running.");
      return;
    }

    App.getInstance().sendMessage(
      "Hey, dotinder here. I'm able to take your orders now."
    );
    App.getInstance().setState(State.TAKE_ORDERS);
  }
}
