import { App } from "../app";
import { State, Transition } from "../states";

export abstract class Command {
  protected app: App;

  public abstract readonly command: string;
  public abstract readonly transition: Transition;

  public constructor(app: App) {
    this.app = app;
  }

  abstract process(): void;
}
