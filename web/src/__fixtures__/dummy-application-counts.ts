import * as jspb from "google-protobuf";
import { APPLICATION_ACTIVE_STATUS_NAME } from "~/constants/application-active-status";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import {
  InsightApplicationCount,
  InsightApplicationCountLabelKey,
} from "~~/model/insight_pb";
import { ApplicationActiveStatus, ApplicationKind } from "~/types/applications";
import { INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT } from "~/queries/application-counts/use-get-application-counts";

export const dummyApplicationCounts: InsightApplicationCount.AsObject[] = [
  {
    count: 123,
    labelsMap: [
      [
        INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
          InsightApplicationCountLabelKey.KIND
        ],
        APPLICATION_KIND_TEXT[ApplicationKind.KUBERNETES],
      ],
      [
        INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
          InsightApplicationCountLabelKey.ACTIVE_STATUS
        ],
        APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.ENABLED],
      ],
    ],
  },
  {
    count: 8,
    labelsMap: [
      [
        INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
          InsightApplicationCountLabelKey.KIND
        ],
        APPLICATION_KIND_TEXT[ApplicationKind.KUBERNETES],
      ],
      [
        INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
          InsightApplicationCountLabelKey.ACTIVE_STATUS
        ],
        APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.DISABLED],
      ],
    ],
  },
  {
    count: 75,
    labelsMap: [
      [
        INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
          InsightApplicationCountLabelKey.KIND
        ],
        APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM],
      ],
      [
        INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
          InsightApplicationCountLabelKey.ACTIVE_STATUS
        ],
        APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.ENABLED],
      ],
    ],
  },
  {
    count: 2,
    labelsMap: [
      [
        INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
          InsightApplicationCountLabelKey.KIND
        ],
        APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM],
      ],
      [
        INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
          InsightApplicationCountLabelKey.ACTIVE_STATUS
        ],
        APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.DISABLED],
      ],
    ],
  },
];

export function createInsightApplicationCountFromObject(
  o: InsightApplicationCount.AsObject
): InsightApplicationCount {
  const count = new InsightApplicationCount();
  count.setCount(o.count);
  const map: jspb.Map<string, string> = count.getLabelsMap();
  o.labelsMap.map((m) => {
    map.set(m[0], m[1]);
  });
  return count;
}
