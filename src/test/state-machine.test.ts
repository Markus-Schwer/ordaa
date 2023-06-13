import { describe, expect, test } from "@jest/globals";
import { Command } from "../commands/command";
import { State, Transition } from "../states";
import { App } from "../app";
import { StateMachine } from "../state-machine";

class MockCommand extends Command {
  constructor() {
    super({} as App);
  }
  internalTransition = Transition.START_ORDER;
  internalOnErrorCalled = false;
  process(rawInput: string, user: string): Promise<void> {
    return Promise.resolve();
  }

  get transition() {
    return this.internalTransition;
  }
  matcher = new RegExp("mock");
  onError(_: string) {
    this.internalOnErrorCalled = true;
  }
}

describe("state machine", () => {
  const cmd = new MockCommand();
  test("multiple orders", () => {
    const machine = new StateMachine();
    expect(machine.getCurrentState()).toEqual(State.IDLE);
    cmd.internalTransition = Transition.START_ORDER;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
    cmd.internalTransition = Transition.ADD_ITEM;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
    cmd.internalTransition = Transition.ADD_ITEM;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
    cmd.internalTransition = Transition.FINALIZE;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.ORDERED);
    cmd.internalTransition = Transition.ARRIVED;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.IDLE);
  });
  test("cancel order", () => {
    const machine = new StateMachine();
    expect(machine.getCurrentState()).toEqual(State.IDLE);
    cmd.internalTransition = Transition.START_ORDER;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
    cmd.internalTransition = Transition.ADD_ITEM;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
    cmd.internalTransition = Transition.ADD_ITEM;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
    cmd.internalTransition = Transition.CANCEL;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.IDLE);
  });
  test("invalid transition", () => {
    const machine = new StateMachine();
    expect(machine.getCurrentState()).toEqual(State.IDLE);
    cmd.internalTransition = Transition.START_ORDER;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
    cmd.internalTransition = Transition.ARRIVED;
    machine.handleState(cmd, "raw", "test-user");
    expect(machine.getCurrentState()).toEqual(State.TAKE_ORDERS);
    expect(cmd.internalOnErrorCalled).toBeTruthy();
  });
});
