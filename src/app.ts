import { Command } from "./commands/command";
import { StartCommand } from "./commands/start-command";
import { Config } from "./config";
import { DeliveredCommand } from "./commands/delivered-command";
import { OrderCommand } from "./commands/order-command";
import { HelpCommand } from "./commands/help-command";
import { StateMachine } from "./state-machine";
import { MatrixClientFacade } from "./matrix-client-facade";
import { Menu, Orders } from "./menu";

export class App {
  public readonly config: Config;

  private commandList: Command[] = [];
  private stateMachine: StateMachine;
  private matrixClient: MatrixClientFacade;
  public orders: Orders = new Orders();
  public menu: Menu | null = null;

  public constructor(
    config: Config,
    stateMachine: StateMachine,
    matrixClient: MatrixClientFacade
  ) {
    this.stateMachine = stateMachine;
    this.matrixClient = matrixClient;
    this.config = config;
    this.matrixClient.sendMessage(".inder is back!");

    this.matrixClient.listenToRoomEvents(this.processMessage);

    console.log(this.config.getRoomId());

    // register your commands here
    this.commandList.push(new StartCommand(this));
    this.commandList.push(new DeliveredCommand(this));
    this.commandList.push(new OrderCommand(this));
    this.commandList.push(new HelpCommand(this));
  }

  private processMessage(message: string): void {
    const command: Command | undefined = this.commandList.find((c) =>
      c.matcher.test(message)
    );
    if (!command) {
      this.sendMessage(`what is '${message}'???`);
      return;
    }
    this.stateMachine.handleState(command, message);
  }

  /* send a matrix message */
  public sendMessage(message: string): void {
    this.matrixClient.sendMessage(message);
  }
}

require("dotenv").config();
new App(new StateMachine(), new MatrixClientFacade(new Config()));
