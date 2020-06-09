import { addDecorator } from "@storybook/react";
import { ThemeDecorator } from "./ThemeDecorator";
import { setupDayjs } from "../src/utils/setup-dayjs";

setupDayjs();

addDecorator(ThemeDecorator);
