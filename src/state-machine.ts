import { Command } from "./commands/command";
import { State } from "./states";

export interface StateMachine {
    getCurrentState(): State;
    // Takes a command from the chat and a 'onError' callback.
    // 'onError' will be called when the transition of the command is not
    // possible in the current state, otherwise the process of the command will
    // be executed.
    // Returns the next state.
    handleState(cmd: Command, onError: () => {}): State;
}
