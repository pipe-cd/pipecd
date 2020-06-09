import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";

export const setupDayjs = () => {
  dayjs.extend(relativeTime);
};
