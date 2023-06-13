import { Menu } from "../menu";
import { Transition } from "../states";
import { Command } from "./command";

export class StartCommand extends Command {
  public transition: Transition = Transition.START_ORDER;
  public matcher = new RegExp("^\\.inder$");

  public async process(_rawInput: string, _user: string): Promise<void> {
    const menuUrl = this.app.config.getMenuUrl();
    if (!menuUrl) {
      console.error("no URL to fetch the menu from is provided");
      return;
    }

    // TODO Johannes: fix fetch call
    /*const response = await fetch(menuUrl);
    if (!response.ok) {
      console.error("response from server was not ok when fetching menu");
      console.error(response);
    }
    const text = await response.text();
    this.app.menu = new Menu(text);
    this.app.sendMessage(this.app.menu.toString());*/
    this.app.sendMessage(
      `@ALL: I will now collect orders, use '!order <item-id>' to add to your order, or '!order reset' to reset your order`
    );
  }
}
