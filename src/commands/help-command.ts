import { Transition } from "../states";
import { Command } from "./command";

export class HelpCommand extends Command {
  public readonly transition: Transition = Transition.HELP;
  public readonly matcher = new RegExp("^.*help.*$");

  public process(): void {
    this.app.sendMessage("+++ HELP +++");

    // TODO: add help message
  }
}
