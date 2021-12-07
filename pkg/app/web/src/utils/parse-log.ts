// eslint-disable-next-line no-control-regex
const ANSI_COLOR_REGEX = /[\u001B\u009B][[()#;?]*(\d{1,4}(?:;\d{0,4})*)?[\dA-PRZcf-nqry=><]/g;
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
};

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

  while ((match = ANSI_COLOR_REGEX.exec(logStr))) {
    if (match.index > 0) {
      const content = logStr.slice(prev, match.index);
      pushToCells(content);
    }
    prev = match.index + match[0].length;

    // Set code for next cell process.
    const code = match[1];
    switch (code) {
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
        fg = 7 - fg;
        bg = 7 - bg;
        break;
    }

    if (typeof TERM_COLOR_FG[code] !== "undefined") {
      fg = TERM_COLOR_FG[code];
    }

    if (typeof TERM_COLOR_BG[code] !== "undefined") {
      bg = TERM_COLOR_BG[code];
    }
  }

  const lastContent = logStr.slice(prev);
  pushToCells(lastContent);

  return cells;
}
