import { parseLog } from "./parse-log";

test("parse string that includes fg colors", () => {
  expect(
    parseLog("\u001b[30m A \u001b[31m B \u001b[32m C \u001b[33m D \u001b[0m")
  ).toMatchSnapshot();

  expect(
    parseLog("\u001b[34m E \u001b[35m F \u001b[36m G \u001b[37m H \u001b[0m")
  ).toMatchSnapshot();
});

test("parse string that includes bg colors", () => {
  expect(
    parseLog("\u001b[40m A \u001b[41m B \u001b[42m C \u001b[43m D \u001b[0m")
  ).toMatchSnapshot();

  expect(
    parseLog("\u001b[44m A \u001b[45m B \u001b[46m C \u001b[47m D \u001b[0m")
  ).toMatchSnapshot();
});

test("parse string that includes modifier", () => {
  [
    "\u001b[1m BOLD \u001b[0m\u001b[4m Underline \u001b[0m\u001b[7m Reversed \u001b[0m",
    "\u001b[1m\u001b[4m\u001b[7m BOLD Underline Reversed \u001b[0m",
    "\u001b[1m\u001b[31m Red Bold \u001b[0m",
    "\u001b[4m\u001b[44m Blue Background Underline \u001b[0m",
  ].forEach((str) => {
    expect(parseLog(str)).toMatchSnapshot();
  });
});

test("parse bright colors", () => {
  [
    "\u001b[1;30m A \u001b[1;31m B \u001b[1;32m C \u001b[1;33m D \u001b[0m",
  ].forEach((str) => {
    expect(parseLog(str)).toMatchSnapshot();
  });
});

test("multi lines contained in a single log block", () => {
  [
    "\n\u001b[1m BOLD \u001b[0m\u001b[4m Underline \n\u001b[0m\u001b[7m Reversed \u001b[0m",
  ].forEach((str) => {
    expect(parseLog(str)).toMatchSnapshot();
  });
});

// Indexes of the colors in TERM_COLORS.
const BLACK = 0;
const RED = 1;
const GREEN = 2;
const YELLOW = 3;
const WHITE = 7;
const BR_BLACK = 8;
const BR_RED = 9;
// parseLog turns spaces into non-breaking spaces.
const NBSP = String.fromCharCode(0xa0);

test("reset without a code clears the style", () => {
  const cells = parseLog("\u001b[1;31mfailed\u001b[m done");

  expect(cells[0]).toMatchObject({ content: "failed", fg: BR_RED });
  expect(cells[1]).toMatchObject({ content: `${NBSP}done`, fg: WHITE });
});

test("applies every code of a combined sequence", () => {
  const cells = parseLog(
    "\u001b[33mwarn\u001b[0;32mok\u001b[0m\u001b[1;4;31mfail\u001b[0m"
  );

  expect(cells[0]).toMatchObject({ content: "warn", fg: YELLOW });
  expect(cells[1]).toMatchObject({
    content: "ok",
    fg: GREEN,
    bold: false,
    underline: false,
  });
  expect(cells[2]).toMatchObject({
    content: "fail",
    fg: RED,
    bold: true,
    underline: true,
  });
});

test("parse bright colors set with 90-97 and 100-107", () => {
  const cells = parseLog("\u001b[90mdebug\u001b[0m\u001b[101malert\u001b[0m");

  expect(cells[0]).toMatchObject({ content: "debug", fg: BR_BLACK });
  expect(cells[1]).toMatchObject({ content: "alert", bg: BR_RED });
});

test("skips the arguments of 256 and true color codes", () => {
  const cells = parseLog(
    "\u001b[1m\u001b[38;2;0;0;0mbold\u001b[38;5;0mstill bold"
  );

  expect(cells[0]).toMatchObject({ content: "bold", bold: true });
  expect(cells[1]).toMatchObject({ content: `still${NBSP}bold`, bold: true });
});

test("escape sequences that don't set a style leave it unchanged", () => {
  const cells = parseLog(
    "\u001b[31mdownloading\u001b[0K\u001b[1Gdone\u001b[2K"
  );

  expect(cells).toHaveLength(2);
  expect(cells[0]).toMatchObject({ content: "downloading", fg: RED });
  expect(cells[1]).toMatchObject({ content: "done", fg: RED, bold: false });
});

test("skips the arguments of underline color codes", () => {
  const cells = parseLog("\u001b[1;4;31mfail\u001b[58;5;0mstill");

  expect(cells[1]).toMatchObject({
    content: "still",
    fg: RED,
    bold: true,
    underline: true,
  });
});

test("reversed swaps the text and background colors", () => {
  const cells = parseLog(
    "\u001b[91;7mdebug\u001b[0m\u001b[101;7malert\u001b[0m"
  );

  expect(cells[0]).toMatchObject({ content: "debug", fg: BLACK, bg: BR_RED });
  expect(cells[1]).toMatchObject({ content: "alert", fg: BR_RED, bg: WHITE });
});

test("an empty code at the start of a sequence is a reset", () => {
  const cells = parseLog("\u001b[1;44mbold\u001b[;32mok");

  expect(cells[1]).toMatchObject({
    content: "ok",
    fg: GREEN,
    bg: BLACK,
    bold: false,
  });
  expect(parseLog("\u001b[1;44mbold\u001b[;32mok")).toEqual(
    parseLog("\u001b[1;44mbold\u001b[0;32mok")
  );
});
