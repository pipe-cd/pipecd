import dayjs from "dayjs";

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
