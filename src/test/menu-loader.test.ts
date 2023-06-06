import { describe, expect, test } from "@jest/globals";
import { Menu, MenuItem } from "../menu";
import * as fs from "fs";

describe("menu scraper", () => {
  const srcString = fs.readFileSync(
    "src/test/resources/test_menu.html",
    "utf8"
  );
  test("find all items", () => {
    const menu = new Menu(srcString);
    expect(menu.items).toHaveLength(136);
  });
  test("entries", () => {
    const menu = new Menu(srcString);
    expect(menu.items).toContainEqual(new MenuItem("62", "Butter Chicken"));
    expect(menu.items).toContainEqual(new MenuItem("42", "Chicken Tikka"));
    expect(menu.items).toContainEqual(new MenuItem("174", "Nan"));
  });
});
