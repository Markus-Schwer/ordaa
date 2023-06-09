import { describe, expect, test } from "@jest/globals";
import { StateMachineImpl } from "../state-machine";
import { Command } from "../commands/command";
import { State, Transition } from "../states";
import { App } from "../app";

class MockCommand extends Command {
    constructor() {
        super({} as App);
    }
    internalTransition = Transition.START_ORDER;
    process(): void {}
    get transition() {
        return this.internalTransition;
    }
    matcher = new RegExp("mock");
}

describe("state machine", () => {
  const cmd = new MockCommand();
  test("multiple orders", () => {
      const machine = new StateMachineImpl();
      expect(machine.getCurrentState()).toEqual(State.IDLE);
      cmd.internalTransition = Transition.START_ORDER;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
      cmd.internalTransition = Transition.ADD_ITEM;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
      cmd.internalTransition = Transition.ADD_ITEM;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
      cmd.internalTransition = Transition.FINALIZE;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.ORDERED);
      cmd.internalTransition = Transition.ARRIVED;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.IDLE);
  });
  test("cancel order", () => {
      const machine = new StateMachineImpl();
      expect(machine.getCurrentState()).toEqual(State.IDLE);
      cmd.internalTransition = Transition.START_ORDER;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
      cmd.internalTransition = Transition.ADD_ITEM;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
      cmd.internalTransition = Transition.ADD_ITEM;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
      cmd.internalTransition = Transition.CANCEL;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.IDLE);
  });
  test("invalid transition", () => {
      const machine = new StateMachineImpl();
      expect(machine.getCurrentState()).toEqual(State.IDLE);
      cmd.internalTransition = Transition.START_ORDER;
      machine.handleState(cmd, () => {});
      expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
      cmd.internalTransition = Transition.ARRIVED;
      let wasCalled = false;
      const onError = () => { wasCalled = true; };
      machine.handleState(cmd, onError);
      expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
      expect(wasCalled).toBeTruthy();
  });
});

