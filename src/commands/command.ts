import { State, Transition } from "../states";

export interface Command {
  readonly command: string;
  readonly transition: Transition;
  process(cmd: string, state: State): void;
}
