import { appendFile } from "fs";
import { State, Transition } from "../states";
import { Command } from "./command";
import { App } from "../app";

export class OrderCommand extends Command {
  public readonly command: string = ".order";
  public readonly transition: Transition = Transition.ADD_ITEM;

  public process(): void {
    /*if (state != State.TAKE_ORDERS) {
      this.app.sendMessage("Sorry, I'm not able to do that currently.");
      return;
    }*/

    this.app.sendMessage("Are you sure that you want to order now?");
  }
}
