export const sortFunc = (
  a: string,
  b: string,
  direction: "ASC" | "DESC" = "ASC"
): number => {
  if (direction === "ASC") return a > b ? 1 : -1;
  return a > b ? -1 : 1;
};
