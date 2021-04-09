import * as React from "react";
import { Log } from "./";
import { LogSeverity } from "../../modules/stage-logs";

export default {
  title: "DEPLOYMENT/Log",
  component: Log,
};

export const overview: React.FC = () => (
  <Log
    logs={[
      "Hello, World",
      "Hello, World",
      "Hello, World",
      "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!" +
        "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!" +
        "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!" +
        "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!" +
        "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!",
      "Hello, World",
    ].map((v, i) => ({
      log: v,
      index: i,
      severity: LogSeverity.INFO,
      createdAt: 0,
    }))}
    loading={false}
  />
);

export const severity: React.FC = () => (
  <Log
    logs={["Hello, World", "Hello, World", "Hello, World"].map((v, i) => ({
      log: v,
      index: i,
      severity: i,
      createdAt: 0,
    }))}
    loading={false}
  />
);

export const loading: React.FC = () => (
  <Log
    logs={["Hello, World", "Hello, World", "Hello, World", "Hello, World"].map(
      (v, i) => ({
        log: v,
        index: i,
        severity: LogSeverity.INFO,
        createdAt: 0,
      })
    )}
    loading
  />
);

export const wordWrap: React.FC = () => (
  <Log
    logs={[
      "Hello, World",
      "Hello, World",
      "Hello, World",
      "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!" +
        "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!" +
        "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!" +
        "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!" +
        "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!",
      "Hello, World",
    ].map((v, i) => ({
      log: v,
      index: i,
      severity: LogSeverity.INFO,
      createdAt: 0,
    }))}
    loading={false}
  />
);

export const ansiCodes: React.FC = () => (
  <Log
    logs={[
      "\u001b[30m A \u001b[31m B \u001b[32m C \u001b[33m D \u001b[0m",
      "\u001b[34m E \u001b[35m F \u001b[36m G \u001b[37m H \u001b[0m",
      "\u001b[40m A \u001b[41m B \u001b[42m C \u001b[43m D \u001b[0m",
      "\u001b[44m A \u001b[45m B \u001b[46m C \u001b[47m D \u001b[0m",
      "\u001b[1m BOLD \u001b[0m\u001b[4m Underline \u001b[0m\u001b[7m Reversed \u001b[0m",
      "\u001b[1m\u001b[4m\u001b[7m BOLD Underline Reversed \u001b[0m",
      "\u001b[1m\u001b[31m Red Bold \u001b[0m",
      "\u001b[4m\u001b[44m Blue Background Underline \u001b[0m",
    ].map((v, i) => ({
      log: v,
      index: i,
      severity: LogSeverity.INFO,
      createdAt: 0,
    }))}
    loading={false}
  />
);

export const allColors: React.FC = () => (
  <Log
    logs={[
      "                 40m     41m     42m     43m     44m     45m     46m     47m",
      "     m \u001b[m  gYw   \u001b[m\u001b[40m  gYw  \u001b[0m \u001b[m\u001b[41m  gYw  \u001b[0m \u001b[m\u001b[42m  gYw  \u001b[0m \u001b[m\u001b[43m  gYw  \u001b[0m \u001b[m\u001b[44m  gYw  \u001b[0m \u001b[m\u001b[45m  gYw  \u001b[0m \u001b[m\u001b[46m  gYw  \u001b[0m \u001b[m\u001b[47m  gYw  \u001b[0m",
      "    1m \u001b[1m  gYw   \u001b[1m\u001b[40m  gYw  \u001b[0m \u001b[1m\u001b[41m  gYw  \u001b[0m \u001b[1m\u001b[42m  gYw  \u001b[0m \u001b[1m\u001b[43m  gYw  \u001b[0m \u001b[1m\u001b[44m  gYw  \u001b[0m \u001b[1m\u001b[45m  gYw  \u001b[0m \u001b[1m\u001b[46m  gYw  \u001b[0m \u001b[1m\u001b[47m  gYw  \u001b[0m",
      "   30m \u001b[30m  gYw   \u001b[30m\u001b[40m  gYw  \u001b[0m \u001b[30m\u001b[41m  gYw  \u001b[0m \u001b[30m\u001b[42m  gYw  \u001b[0m \u001b[30m\u001b[43m  gYw  \u001b[0m \u001b[30m\u001b[44m  gYw  \u001b[0m \u001b[30m\u001b[45m  gYw  \u001b[0m \u001b[30m\u001b[46m  gYw  \u001b[0m \u001b[30m\u001b[47m  gYw  \u001b[0m",
      " 1;30m \u001b[1;30m  gYw   \u001b[1;30m\u001b[40m  gYw  \u001b[0m \u001b[1;30m\u001b[41m  gYw  \u001b[0m \u001b[1;30m\u001b[42m  gYw  \u001b[0m \u001b[1;30m\u001b[43m  gYw  \u001b[0m \u001b[1;30m\u001b[44m  gYw  \u001b[0m \u001b[1;30m\u001b[45m  gYw  \u001b[0m \u001b[1;30m\u001b[46m  gYw  \u001b[0m \u001b[1;30m\u001b[47m  gYw  \u001b[0m",
      "   31m \u001b[31m  gYw   \u001b[31m\u001b[40m  gYw  \u001b[0m \u001b[31m\u001b[41m  gYw  \u001b[0m \u001b[31m\u001b[42m  gYw  \u001b[0m \u001b[31m\u001b[43m  gYw  \u001b[0m \u001b[31m\u001b[44m  gYw  \u001b[0m \u001b[31m\u001b[45m  gYw  \u001b[0m \u001b[31m\u001b[46m  gYw  \u001b[0m \u001b[31m\u001b[47m  gYw  \u001b[0m",
      " 1;31m \u001b[1;31m  gYw   \u001b[1;31m\u001b[40m  gYw  \u001b[0m \u001b[1;31m\u001b[41m  gYw  \u001b[0m \u001b[1;31m\u001b[42m  gYw  \u001b[0m \u001b[1;31m\u001b[43m  gYw  \u001b[0m \u001b[1;31m\u001b[44m  gYw  \u001b[0m \u001b[1;31m\u001b[45m  gYw  \u001b[0m \u001b[1;31m\u001b[46m  gYw  \u001b[0m \u001b[1;31m\u001b[47m  gYw  \u001b[0m",
      "   32m \u001b[32m  gYw   \u001b[32m\u001b[40m  gYw  \u001b[0m \u001b[32m\u001b[41m  gYw  \u001b[0m \u001b[32m\u001b[42m  gYw  \u001b[0m \u001b[32m\u001b[43m  gYw  \u001b[0m \u001b[32m\u001b[44m  gYw  \u001b[0m \u001b[32m\u001b[45m  gYw  \u001b[0m \u001b[32m\u001b[46m  gYw  \u001b[0m \u001b[32m\u001b[47m  gYw  \u001b[0m",
      " 1;32m \u001b[1;32m  gYw   \u001b[1;32m\u001b[40m  gYw  \u001b[0m \u001b[1;32m\u001b[41m  gYw  \u001b[0m \u001b[1;32m\u001b[42m  gYw  \u001b[0m \u001b[1;32m\u001b[43m  gYw  \u001b[0m \u001b[1;32m\u001b[44m  gYw  \u001b[0m \u001b[1;32m\u001b[45m  gYw  \u001b[0m \u001b[1;32m\u001b[46m  gYw  \u001b[0m \u001b[1;32m\u001b[47m  gYw  \u001b[0m",
      "   33m \u001b[33m  gYw   \u001b[33m\u001b[40m  gYw  \u001b[0m \u001b[33m\u001b[41m  gYw  \u001b[0m \u001b[33m\u001b[42m  gYw  \u001b[0m \u001b[33m\u001b[43m  gYw  \u001b[0m \u001b[33m\u001b[44m  gYw  \u001b[0m \u001b[33m\u001b[45m  gYw  \u001b[0m \u001b[33m\u001b[46m  gYw  \u001b[0m \u001b[33m\u001b[47m  gYw  \u001b[0m",
      " 1;33m \u001b[1;33m  gYw   \u001b[1;33m\u001b[40m  gYw  \u001b[0m \u001b[1;33m\u001b[41m  gYw  \u001b[0m \u001b[1;33m\u001b[42m  gYw  \u001b[0m \u001b[1;33m\u001b[43m  gYw  \u001b[0m \u001b[1;33m\u001b[44m  gYw  \u001b[0m \u001b[1;33m\u001b[45m  gYw  \u001b[0m \u001b[1;33m\u001b[46m  gYw  \u001b[0m \u001b[1;33m\u001b[47m  gYw  \u001b[0m",
      "   34m \u001b[34m  gYw   \u001b[34m\u001b[40m  gYw  \u001b[0m \u001b[34m\u001b[41m  gYw  \u001b[0m \u001b[34m\u001b[42m  gYw  \u001b[0m \u001b[34m\u001b[43m  gYw  \u001b[0m \u001b[34m\u001b[44m  gYw  \u001b[0m \u001b[34m\u001b[45m  gYw  \u001b[0m \u001b[34m\u001b[46m  gYw  \u001b[0m \u001b[34m\u001b[47m  gYw  \u001b[0m",
      " 1;34m \u001b[1;34m  gYw   \u001b[1;34m\u001b[40m  gYw  \u001b[0m \u001b[1;34m\u001b[41m  gYw  \u001b[0m \u001b[1;34m\u001b[42m  gYw  \u001b[0m \u001b[1;34m\u001b[43m  gYw  \u001b[0m \u001b[1;34m\u001b[44m  gYw  \u001b[0m \u001b[1;34m\u001b[45m  gYw  \u001b[0m \u001b[1;34m\u001b[46m  gYw  \u001b[0m \u001b[1;34m\u001b[47m  gYw  \u001b[0m",
      "   35m \u001b[35m  gYw   \u001b[35m\u001b[40m  gYw  \u001b[0m \u001b[35m\u001b[41m  gYw  \u001b[0m \u001b[35m\u001b[42m  gYw  \u001b[0m \u001b[35m\u001b[43m  gYw  \u001b[0m \u001b[35m\u001b[44m  gYw  \u001b[0m \u001b[35m\u001b[45m  gYw  \u001b[0m \u001b[35m\u001b[46m  gYw  \u001b[0m \u001b[35m\u001b[47m  gYw  \u001b[0m",
      " 1;35m \u001b[1;35m  gYw   \u001b[1;35m\u001b[40m  gYw  \u001b[0m \u001b[1;35m\u001b[41m  gYw  \u001b[0m \u001b[1;35m\u001b[42m  gYw  \u001b[0m \u001b[1;35m\u001b[43m  gYw  \u001b[0m \u001b[1;35m\u001b[44m  gYw  \u001b[0m \u001b[1;35m\u001b[45m  gYw  \u001b[0m \u001b[1;35m\u001b[46m  gYw  \u001b[0m \u001b[1;35m\u001b[47m  gYw  \u001b[0m",
      "   36m \u001b[36m  gYw   \u001b[36m\u001b[40m  gYw  \u001b[0m \u001b[36m\u001b[41m  gYw  \u001b[0m \u001b[36m\u001b[42m  gYw  \u001b[0m \u001b[36m\u001b[43m  gYw  \u001b[0m \u001b[36m\u001b[44m  gYw  \u001b[0m \u001b[36m\u001b[45m  gYw  \u001b[0m \u001b[36m\u001b[46m  gYw  \u001b[0m \u001b[36m\u001b[47m  gYw  \u001b[0m",
      " 1;36m \u001b[1;36m  gYw   \u001b[1;36m\u001b[40m  gYw  \u001b[0m \u001b[1;36m\u001b[41m  gYw  \u001b[0m \u001b[1;36m\u001b[42m  gYw  \u001b[0m \u001b[1;36m\u001b[43m  gYw  \u001b[0m \u001b[1;36m\u001b[44m  gYw  \u001b[0m \u001b[1;36m\u001b[45m  gYw  \u001b[0m \u001b[1;36m\u001b[46m  gYw  \u001b[0m \u001b[1;36m\u001b[47m  gYw  \u001b[0m",
      "   37m \u001b[37m  gYw   \u001b[37m\u001b[40m  gYw  \u001b[0m \u001b[37m\u001b[41m  gYw  \u001b[0m \u001b[37m\u001b[42m  gYw  \u001b[0m \u001b[37m\u001b[43m  gYw  \u001b[0m \u001b[37m\u001b[44m  gYw  \u001b[0m \u001b[37m\u001b[45m  gYw  \u001b[0m \u001b[37m\u001b[46m  gYw  \u001b[0m \u001b[37m\u001b[47m  gYw  \u001b[0m",
      " 1;37m \u001b[1;37m  gYw   \u001b[1;37m\u001b[40m  gYw  \u001b[0m \u001b[1;37m\u001b[41m  gYw  \u001b[0m \u001b[1;37m\u001b[42m  gYw  \u001b[0m \u001b[1;37m\u001b[43m  gYw  \u001b[0m \u001b[1;37m\u001b[44m  gYw  \u001b[0m \u001b[1;37m\u001b[45m  gYw  \u001b[0m \u001b[1;37m\u001b[46m  gYw  \u001b[0m \u001b[1;37m\u001b[47m  gYw  \u001b[0m",
    ].map((v, i) => ({
      log: v,
      index: i,
      severity: LogSeverity.INFO,
      createdAt: 0,
    }))}
    loading={false}
  />
);
