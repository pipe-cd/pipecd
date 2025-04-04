import dayjs from "dayjs";
import { PIPED_VERSION } from "~/modules/applications";
import { Application } from "~~/model/application_pb";

export const sortFunc = (
  a: string | number,
  b: string | number,
  direction: "ASC" | "DESC" = "ASC"
): number => {
  if (direction === "ASC") return a > b ? 1 : -1;
  return a > b ? -1 : 1;
};

export const sortDateFunc = (
  a: string | number,
  b: string | number,
  direction: "ASC" | "DESC" = "ASC"
): number => {
  const dateA = dayjs(a).valueOf();
  const dateB = dayjs(b).valueOf();
  return sortFunc(dateA, dateB, direction);
};

export const getPercentage = (
  num: number,
  total: number,
  precision?: number
): number => {
  const percentage = (num / (total || 1)) * 100;
  return precision ? Number(percentage.toFixed(precision)) : percentage;
};

export const getPipedAppVersion = (
  app?: Application.AsObject
): PIPED_VERSION => {
  if (!app) return PIPED_VERSION.V0;
  if (!app.platformProvider) return PIPED_VERSION.V1;
  if (app?.deployTargetsByPluginMap?.length > 0) return PIPED_VERSION.V1;
  return PIPED_VERSION.V0;
};

export const checkPipedAppVersion = (
  app?: Application.AsObject
): Record<PIPED_VERSION, boolean> => {
  const pipedVersion = getPipedAppVersion(app);
  return {
    [PIPED_VERSION.V0]: pipedVersion === PIPED_VERSION.V0,
    [PIPED_VERSION.V1]: pipedVersion === PIPED_VERSION.V1,
  };
};
