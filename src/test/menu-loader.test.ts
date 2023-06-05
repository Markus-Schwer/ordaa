import { describe, expect, test } from "@jest/globals";
import { Menu } from "../menu";
import * as fs from "fs";

describe("menu scraper", () => {
  const srcString = fs.readFileSync(
    "src/test/resources/test_menu.html",
    "utf8"
  );
  test("parsing", () => {
    const menu = new Menu(srcString);
    expect(menu.items.length).toBe(136);
  });
});
