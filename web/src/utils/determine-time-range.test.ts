import dayjs from "dayjs";
import { InsightRange } from "~/queries/insight/insight.config";
import { InsightResolution } from "~~/model/insight_pb";
import { determineTimeRange } from "./determine-time-range";

beforeAll(() => {
  jest.useFakeTimers();
  jest.setSystemTime(new Date("2024-06-15T12:00:00.000Z"));
});

afterAll(() => {
  jest.useRealTimers();
});

describe("determineTimeRange", () => {
  describe("DAILY resolution", () => {
    it("determines time range for LAST_1_WEEK", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_1_WEEK,
        InsightResolution.DAILY
      );

      expect(dayjs(from).toISOString()).toBe("2024-06-09T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_1_MONTH", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_1_MONTH,
        InsightResolution.DAILY
      );

      expect(dayjs(from).toISOString()).toBe("2024-05-16T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_3_MONTHS", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_3_MONTHS,
        InsightResolution.DAILY
      );

      expect(dayjs(from).toISOString()).toBe("2024-03-16T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_6_MONTHS", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_6_MONTHS,
        InsightResolution.DAILY
      );

      expect(dayjs(from).toISOString()).toBe("2023-12-16T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_1_YEAR", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_1_YEAR,
        InsightResolution.DAILY
      );

      expect(dayjs(from).toISOString()).toBe("2023-06-16T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_2_YEARS", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_2_YEARS,
        InsightResolution.DAILY
      );

      expect(dayjs(from).toISOString()).toBe("2022-06-16T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });
  });

  describe("MONTHLY resolution", () => {
    it("determines time range for LAST_1_WEEK (maps to month-to-date)", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_1_WEEK,
        InsightResolution.MONTHLY
      );

      expect(dayjs(from).toISOString()).toBe("2024-06-01T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_1_MONTH (maps to month-to-date)", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_1_MONTH,
        InsightResolution.MONTHLY
      );

      expect(dayjs(from).toISOString()).toBe("2024-06-01T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_3_MONTHS", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_3_MONTHS,
        InsightResolution.MONTHLY
      );

      expect(dayjs(from).toISOString()).toBe("2024-04-01T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_6_MONTHS", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_6_MONTHS,
        InsightResolution.MONTHLY
      );

      expect(dayjs(from).toISOString()).toBe("2024-01-01T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_1_YEAR", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_1_YEAR,
        InsightResolution.MONTHLY
      );

      expect(dayjs(from).toISOString()).toBe("2023-07-01T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });

    it("determines time range for LAST_2_YEARS", () => {
      const [from, to] = determineTimeRange(
        InsightRange.LAST_2_YEARS,
        InsightResolution.MONTHLY
      );

      expect(dayjs(from).toISOString()).toBe("2022-07-01T00:00:00.000Z");
      expect(dayjs(to).toISOString()).toBe("2024-06-15T23:59:59.999Z");
    });
  });
});
