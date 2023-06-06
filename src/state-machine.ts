import { Command } from "./commands/command";
import { Transition, State } from "./states";

export interface StateMachine {
    getCurrentState(): State;
    // Takes a command from the chat and a 'onError' callback.
    // 'onError' will be called when the transition of the command is not
    // possible in the current state, otherwise the process of the command will
    // be executed.
    // Returns the next state.
    handleState(cmd: Command, onError: () => void): State;
}

type TransitionTable = {
    [s in State]: {
        [t in Transition]: State | null
    }
}

export class StateMachineImpl implements StateMachine {
    private currentState: State = State.IDLE;

    getCurrentState(): State {
        return this.currentState;
    }

    private t: TransitionTable = {
        [State.IDLE]: {
            [Transition.START_ORDER]: State.TAKE_ORDERS,
            [Transition.ADD_ITEM]: null,
            [Transition.CANCEL]: null,
            [Transition.FINALIZE]: null,
            [Transition.ARRIVED]: null,
        },
        [State.TAKE_ORDERS]: {
            [Transition.START_ORDER]: null,
            [Transition.ADD_ITEM]: State.TAKE_ORDERS,
            [Transition.CANCEL]: State.IDLE,
            [Transition.FINALIZE]: State.ORDERED,
            [Transition.ARRIVED]: null,
        },
        [State.ORDERED]: {
            [Transition.START_ORDER]: null,
            [Transition.ADD_ITEM]: null,
            [Transition.CANCEL]: null,
            [Transition.FINALIZE]: null,
            [Transition.ARRIVED]: State.DELIVERED,
        },
        [State.DELIVERED]: {
            [Transition.START_ORDER]: null,
            [Transition.ADD_ITEM]: null,
            [Transition.CANCEL]: null,
            [Transition.FINALIZE]: State.IDLE,
            [Transition.ARRIVED]: null,
        },
    };

    handleState(cmd: Command, onError: () => void): State {
        const nextState = this.t[this.currentState][cmd.transition];
        if (nextState == null) {
            onError();
            return this.currentState;
        }
        cmd.process();
        this.currentState = nextState;
        return this.currentState;
    }
}

