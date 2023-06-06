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
      const menuPrice = parseFloat(menuParser(".menuItemPrice", b).text().replace(",", "."));
      const regexResult = this.menuItemNameRegex.exec(menuName);
      if (
        regexResult == null ||
        regexResult.groups == null ||
        regexResult.groups.name == null
      ) {
        return;
      }
      this.menuItemList.push(
        new MenuItem(regexResult.groups.id, regexResult.groups.name.trim(), menuPrice)
      );
    });
    // set this at the end so it's only set when parsing was successful
    this.lastHash = newHash;
    return true;
  }

  private buildHash(src: string): string {
    return crypto.createHash("md5").update(src).digest("hex");
  }
}
