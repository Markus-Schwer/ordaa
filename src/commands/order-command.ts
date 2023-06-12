import { Transition } from "../states";
import { Command } from "./command";

export class OrderCommand extends Command {
  public readonly matcher: RegExp = new RegExp("^!order (\\w{0,1}\\d+)$");
  public readonly transition: Transition = Transition.ADD_ITEM;

  public async process(rawInput: string, user: string): Promise<void> {
    const regexResult = this.matcher.exec(rawInput);
    if (!regexResult) {
      console.error(`invalid command string '${rawInput}' in OrderCommand`);
      return;
    }
    if (!this.app.menu) {
      console.error("menu is not loaded, cannot check if ordered item exists");
      return;
    }
    const orderedMenuItem = this.app.menu.items.find(
      (it) => it.id == regexResult[1]
    );
    if (!orderedMenuItem) {
      this.app.sendMessage(`@${user}: [${regexResult[1]}] is not on the menu`);
      return;
    }
    this.app.orders.orderItem(user, orderedMenuItem);
    const msg = this.app.orders.getUserMessage(user);
    if (!msg) {
      console.error(
        "something went wrong, user message after ordering was empty"
      );
      return;
    }
    this.app.sendMessage(msg);
  }
}
