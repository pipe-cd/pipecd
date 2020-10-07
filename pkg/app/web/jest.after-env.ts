import { setAutoFreeze } from "immer";
import { setupDayjs } from "./src/utils/setup-dayjs";
import "@testing-library/jest-dom";

setupDayjs();
setAutoFreeze(false);
