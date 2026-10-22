// eslint-disable-next-line no-control-regex
const ANSI_COLOR_REGEX = /[\u001B\u009B][[()#?]*(\d{0,4}(?:;\d{0,4})*)?[\dA-PRZcf-nqry=><]/g;
const NON_BREAK_SPACE = "\u00a0";

export interface Cell {
  content: string;
  fg: number;
  bg: number;
  bold: boolean;
  underline: boolean;
}

enum ANSI_CODE {
  RESET = "0",
  BOLD = "1",
  UNDERLINE = "4",
  REVERSED = "7",
}

enum TERM_COLOR {
  BLACK,
  RED,
  GREEN,
  YELLOW,
  BLUE,
  MAGENTA,
  CYAN,
  WHITE,
  BR_BLACK,
  BR_RED,
  BR_GREEN,
  BR_YELLOW,
  BR_BLUE,
  BR_MAGENTA,
  BR_CYAN,
  BR_WHITE,
}

const TERM_COLOR_FG: Record<string, TERM_COLOR> = {
  30: TERM_COLOR.BLACK,
  31: TERM_COLOR.RED,
  32: TERM_COLOR.GREEN,
  33: TERM_COLOR.YELLOW,
  34: TERM_COLOR.BLUE,
  35: TERM_COLOR.MAGENTA,
  36: TERM_COLOR.CYAN,
  37: TERM_COLOR.WHITE,
  39: TERM_COLOR.WHITE,
  // bright
  "1;30": TERM_COLOR.BR_BLACK,
  "1;31": TERM_COLOR.BR_RED,
  "1;32": TERM_COLOR.BR_GREEN,
  "1;33": TERM_COLOR.BR_YELLOW,
  "1;34": TERM_COLOR.BR_BLUE,
  "1;35": TERM_COLOR.BR_MAGENTA,
  "1;36": TERM_COLOR.BR_CYAN,
  "1;37": TERM_COLOR.BR_WHITE,
  90: TERM_COLOR.BR_BLACK,
  91: TERM_COLOR.BR_RED,
  92: TERM_COLOR.BR_GREEN,
  93: TERM_COLOR.BR_YELLOW,
  94: TERM_COLOR.BR_BLUE,
  95: TERM_COLOR.BR_MAGENTA,
  96: TERM_COLOR.BR_CYAN,
  97: TERM_COLOR.BR_WHITE,
};

const TERM_COLOR_BG: Record<string, TERM_COLOR> = {
  40: TERM_COLOR.BLACK,
  41: TERM_COLOR.RED,
  42: TERM_COLOR.GREEN,
  43: TERM_COLOR.YELLOW,
  44: TERM_COLOR.BLUE,
  45: TERM_COLOR.MAGENTA,
  46: TERM_COLOR.CYAN,
  47: TERM_COLOR.WHITE,
  49: TERM_COLOR.BLACK,
  // bright
  "1;40": TERM_COLOR.BR_BLACK,
  "1;41": TERM_COLOR.BR_RED,
  "1;42": TERM_COLOR.BR_GREEN,
  "1;43": TERM_COLOR.BR_YELLOW,
  "1;44": TERM_COLOR.BR_BLUE,
  "1;45": TERM_COLOR.BR_MAGENTA,
  "1;46": TERM_COLOR.BR_CYAN,
  "1;47": TERM_COLOR.BR_WHITE,
  100: TERM_COLOR.BR_BLACK,
  101: TERM_COLOR.BR_RED,
  102: TERM_COLOR.BR_GREEN,
  103: TERM_COLOR.BR_YELLOW,
  104: TERM_COLOR.BR_BLUE,
  105: TERM_COLOR.BR_MAGENTA,
  106: TERM_COLOR.BR_CYAN,
  107: TERM_COLOR.BR_WHITE,
};

// Codes that select a 256-color or a true color for the text, the background or
// the underline take extra arguments ("5;n" or "2;r;g;b"). Those colors aren't
// supported, but the arguments must be skipped so that a 0 in them isn't taken
// as a reset.
const EXTENDED_COLOR_CODES = ["38", "48", "58"];
const EXTENDED_COLOR_ARGS: Record<string, number> = { 5: 2, 2: 4 };

export function parseLog(logStr: string): Cell[] {
  let match = null;
  const cells: Cell[] = [];
  let prev = 0;
  let fg = TERM_COLOR.WHITE;
  let bg = TERM_COLOR.BLACK;
  let bold = false;
  let underline = false;

  const pushToCells = (content: string): void => {
    if (content.length > 0) {
      cells.push({
        content: content.replace(/\n/g, "\\n").replace(/\s/g, NON_BREAK_SPACE),
        fg,
        bg,
        bold,
        underline,
      });
    }
  };

  const applyCode = (code: string): void => {
    switch (code) {
      case "":
      case ANSI_CODE.RESET:
        fg = TERM_COLOR.WHITE;
        bg = TERM_COLOR.BLACK;
        bold = false;
        underline = false;
        break;
      case ANSI_CODE.BOLD:
        bold = true;
        break;
      case ANSI_CODE.UNDERLINE:
        underline = true;
        break;
      case ANSI_CODE.REVERSED:
        [fg, bg] = [bg, fg];
        break;
    }

    if (typeof TERM_COLOR_FG[code] !== "undefined") {
      fg = TERM_COLOR_FG[code];
    }

    if (typeof TERM_COLOR_BG[code] !== "undefined") {
      bg = TERM_COLOR_BG[code];
    }
  };

  while ((match = ANSI_COLOR_REGEX.exec(logStr))) {
    if (match.index > 0) {
      const content = logStr.slice(prev, match.index);
      pushToCells(content);
    }
    prev = match.index + match[0].length;

    // Only sequences ending with "m" set the style. Others, like the ones that
    // move the cursor or erase the line, are dropped without changing it.
    if (!match[0].endsWith("m")) {
      continue;
    }

    // Set code for next cell process. An empty code is a reset, so "ESC[m" is
    // the same as "ESC[0m" and "ESC[;32m" the same as "ESC[0;32m".
    const code = match[1] ?? ANSI_CODE.RESET;
    const isKnownCode =
      typeof TERM_COLOR_FG[code] !== "undefined" ||
      typeof TERM_COLOR_BG[code] !== "undefined";
    if (isKnownCode || !code.includes(";")) {
      applyCode(code);
      continue;
    }

    // A sequence such as "ESC[0;32m" or "ESC[1;4;31m" sets several codes at once.
    const codes = code.split(";");
    for (let i = 0; i < codes.length; i++) {
      if (EXTENDED_COLOR_CODES.includes(codes[i])) {
        i += EXTENDED_COLOR_ARGS[codes[i + 1]] ?? 0;
        continue;
      }
      applyCode(codes[i]);
    }
  }

  const lastContent = logStr.slice(prev);
  pushToCells(lastContent);

  return cells;
}
