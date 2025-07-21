import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import isBetween from "dayjs/plugin/isBetween";
import relativeTime from "dayjs/plugin/relativeTime";
import advancedFormat from "dayjs/plugin/advancedFormat";

export const setupDayjs = (): void => {
  dayjs.extend(relativeTime);
  dayjs.extend(isBetween);
  dayjs.extend(advancedFormat);
  dayjs.extend(utc);
};
