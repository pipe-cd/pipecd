import dayjs, { Dayjs } from "dayjs";
import faker from "faker";

faker.seed(1);

export function randomUUID(): string {
  return faker.random.uuid();
}

export function randomKeyHash(): string {
  return faker.random.alphaNumeric(128);
}

export function randomWords(num: number): string {
  return faker.lorem.words(num);
}

export function randomText(times: number): string {
  return faker.lorem.text(times);
}

export function createRandTime(): Dayjs {
  return dayjs().subtract(faker.random.number({ min: 1, max: 10 }), "minute");
}

export const randomNumber = faker.random.number;

function subtractRandTimeFrom(t: Dayjs): Dayjs {
  return t.subtract(faker.random.number({ min: 5, max: 30 }), "minute");
}

/**
 * Generate times that are ascending ordered.
 */
export function createRandTimes(count: number): Dayjs[] {
  const times: Dayjs[] = [createRandTime()];

  if (count === 1) {
    return times;
  }

  for (let i = 1; i < count; i++) {
    times.push(subtractRandTimeFrom(times[i - 1]));
  }

  return times.reverse();
}
