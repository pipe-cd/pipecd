import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { StatusCode } from "grpc-web";
import {
  ApplicationActiveStatus,
  ApplicationKind,
} from "pipe/pkg/app/web/model/common_pb";
import { InsightApplicationCountLabelKey } from "pipe/pkg/app/web/model/insight_pb";
import { getApplicationCount } from "~/api/insight";
import { APPLICATION_ACTIVE_STATUS_NAME } from "~/constants/application-active-status";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";

const MODULE_NAME = "applicationCounts";

export const INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT: Record<
  InsightApplicationCountLabelKey,
  string
> = {
  [InsightApplicationCountLabelKey.KIND]: "KIND",
  [InsightApplicationCountLabelKey.ACTIVE_STATUS]: "ACTIVE_STATUS",
};

interface ApplicationCounts {
  updatedAt: number;
  counts: Record<string, Record<string, number>>;
}

const createInitialCount = (): Record<string, number> => ({
  [APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.ENABLED]]: 0,
  [APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.DISABLED]]: 0,
});

const createInitialCounts = (): Record<string, Record<string, number>> => ({
  [APPLICATION_KIND_TEXT[ApplicationKind.KUBERNETES]]: createInitialCount(),
  [APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM]]: createInitialCount(),
  [APPLICATION_KIND_TEXT[ApplicationKind.LAMBDA]]: createInitialCount(),
  [APPLICATION_KIND_TEXT[ApplicationKind.CLOUDRUN]]: createInitialCount(),
  [APPLICATION_KIND_TEXT[ApplicationKind.ECS]]: createInitialCount(),
});

const initialState: ApplicationCounts = {
  updatedAt: 0,
  counts: createInitialCounts(),
};

export const fetchApplicationCount = createAsyncThunk(
  `${MODULE_NAME}/fetch`,
  async (): Promise<ApplicationCounts> => {
    const res = await getApplicationCount().catch((e: { code: number }) => {
      // NOT_FOUND is the initial normal state, so it is excluded from error handling.
      if (e.code !== StatusCode.NOT_FOUND) {
        throw e;
      }
    });

    if (!res) {
      // Handling of errors with code NOT_FOUND.
      return initialState;
    }

    const counts: Record<
      string,
      Record<string, number>
    > = createInitialCounts();

    res.countsList.forEach((count) => {
      const [, kindName] =
        count.labelsMap.find(
          (val) =>
            val[0] ===
            INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
              InsightApplicationCountLabelKey.KIND
            ]
        ) || [];
      const [, activeStatusName] =
        count.labelsMap.find(
          (val) =>
            val[0] ===
            INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
              InsightApplicationCountLabelKey.ACTIVE_STATUS
            ]
        ) || [];
      if (!kindName || !activeStatusName) {
        return;
      }

      counts[kindName][activeStatusName] = count.count;
    });

    return { updatedAt: res.updatedAt, counts };
  }
);

export const applicationCountsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchApplicationCount.fulfilled, (_, action) => {
      return action.payload;
    });
  },
});

export {
  InsightApplicationCount,
  InsightApplicationCountLabelKey,
} from "pipe/pkg/app/web/model/insight_pb";
