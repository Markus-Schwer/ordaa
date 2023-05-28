import { State } from "../states";

export interface Command {
    readonly command: string;
    process(cmd: string, state: State): void;
}