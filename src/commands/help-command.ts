import { State, Transition } from "../states";
import { Command } from "./command";
import { App } from "../app";

export class HelpCommand extends Command {
  public readonly command: string = ".help";
  public readonly transition: Transition = Transition.HELP;

  public process(): void {
    this.app.sendMessage("+++ HELP +++");

    // TODO: add help message
  }
}
