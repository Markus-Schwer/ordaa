import { Transition } from "../states";
import { Command } from "./command";

export class HelpCommand extends Command {
  public readonly transition: Transition = Transition.HELP;
  public readonly matcher = new RegExp("^.*help.*$");

  public async process(rawInput: string, user: string): Promise<void> {
    this.app.sendMessage("+++ HELP +++");

    // TODO: add help message
  }
}
