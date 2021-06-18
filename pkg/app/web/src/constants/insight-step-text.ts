import { InsightStep } from "~/modules/insight";

export const INSIGHT_STEP_TEXT: Record<InsightStep, string> = {
  [InsightStep.DAILY]: "DAILY",
  [InsightStep.MONTHLY]: "MONTHLY",
  [InsightStep.WEEKLY]: "WEEKLY",
  [InsightStep.YEARLY]: "YEARLY",
};
