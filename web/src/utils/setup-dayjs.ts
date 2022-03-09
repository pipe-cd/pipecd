import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import isBetween from "dayjs/plugin/isBetween";
import advancedFormat from "dayjs/plugin/advancedFormat";

export const setupDayjs = (): void => {
  dayjs.extend(relativeTime);
  dayjs.extend(isBetween);
  dayjs.extend(advancedFormat);
};
