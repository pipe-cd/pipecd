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
