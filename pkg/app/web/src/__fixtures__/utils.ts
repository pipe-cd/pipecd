import faker from "faker";
import dayjs, { Dayjs } from "dayjs";

faker.seed(1);

export function createdRandTime(): Dayjs {
  return dayjs().subtract(faker.random.number({ min: 1, max: 10 }), "minute");
}

export function subtractRandTimeFrom(t: Dayjs): Dayjs {
  return t.subtract(faker.random.number({ min: 5, max: 30 }), "minute");
}
