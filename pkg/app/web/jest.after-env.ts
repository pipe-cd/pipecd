import { setAutoFreeze } from "immer";
import { setupDayjs } from "./src/utils/setup-dayjs";

setupDayjs();
setAutoFreeze(false);
