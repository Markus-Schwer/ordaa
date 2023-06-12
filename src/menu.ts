import * as cheerio from "cheerio";
import * as crypto from "crypto";

export class MenuItem {
  public readonly id: string | undefined;
  public readonly name: string;
  public readonly price: number | undefined;
  constructor(id: string | undefined, name: string, price: number | undefined) {
    this.id = id;
    this.name = name;
    this.price = price;
  }

  public toString(): string {
    let segments: string[] = [];
    if (this.id) {
      segments.push(`[${this.id}]`);
    }
    segments.push(this.name);
    segments.push(`(${this.price} Euro)`);
    return segments.join(" ");
  }
}

export class Orders {
  private itemMap: { [user: string]: MenuItem[] } = {};

  public orderItem(user: string, item: MenuItem) {
    if (!this.itemMap[user]) {
      this.itemMap[user] = [];
    }
    this.itemMap[user].push(item);
  }

  public resetItems(user: string) {
    delete this.itemMap[user];
  }

  public getUserMessage(user: string): string | null {
    const items: MenuItem[] = this.itemMap[user];
    if (!items) return null;
    let message = `@${user}: Your current order:`;
    for (let item of items) {
      message += `\n[${item.id}] ${item.name}`;
    }
    return message;
  }
}

export class Menu {
  private lastHash: string = "";
  private menuItemNameRegex = new RegExp(
    "^((?<id>\\w*\\d+)\\s*[-–]{1}\\s*)?(?<name>([\\w\\.-]{2,} ?)+).*$"
  );
  private menuItemList: Array<MenuItem> = [];

  public get items(): Array<MenuItem> {
    return this.menuItemList;
  }

  constructor(rawMenuHtml: string) {
    this.possiblyUpdate(rawMenuHtml);
  }

  public possiblyUpdate(rawMenuHtml: string): boolean {
    const newHash = this.buildHash(rawMenuHtml);
    if (newHash == this.lastHash) {
      return false;
    }
    const menuParser = cheerio.load(rawMenuHtml);
    this.menuItemList = [];
    menuParser(".menuItemBox").map((_, b) => {
      const menuName = menuParser(".menuItemName", b).text();
      const menuPrice = parseFloat(
        menuParser(".menuItemPrice", b).text().replace(",", ".")
      );
      const regexResult = this.menuItemNameRegex.exec(menuName);
      if (
        regexResult == null ||
        regexResult.groups == null ||
        regexResult.groups.name == null
      ) {
        return;
      }
      this.menuItemList.push(
        new MenuItem(
          regexResult.groups.id,
          regexResult.groups.name.trim(),
          menuPrice
        )
      );
    });
    // set this at the end so it's only set when parsing was successful
    this.lastHash = newHash;
    return true;
  }

  private buildHash(src: string): string {
    return crypto.createHash("md5").update(src).digest("hex");
  }

  public toString(): string {
    return this.items.join("/n");
  }
}
