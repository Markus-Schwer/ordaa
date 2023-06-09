import { App } from "../app";
import { Transition } from "../states";

export abstract class Command {
  protected app: App;

  public abstract readonly matcher: RegExp;
  public abstract readonly transition: Transition;

  public constructor(app: App) {
    this.app = app;
  }

  abstract process(): void;
}
