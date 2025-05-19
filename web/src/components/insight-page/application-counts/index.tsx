import { FC, memo } from "react";
import { APPLICATION_ACTIVE_STATUS_NAME } from "~/constants/application-active-status";
import { APPLICATION_KIND_BY_NAME } from "~/constants/application-kind";
import { useAppSelector } from "~/hooks/redux";
import {
  ApplicationActiveStatus,
  ApplicationKind,
} from "~/modules/applications";
import { ApplicationCount } from "./application-count";
import { Box } from "@mui/material";

interface ApplicationCountsProps {
  onClick: (kind: ApplicationKind) => void;
}

export const ApplicationCounts: FC<ApplicationCountsProps> = memo(
  function ApplicationCounts({ onClick }) {
    const counts = useAppSelector<Record<string, Record<string, number>>>(
      (state) => state.applicationCounts.counts
    );

    return (
      <Box
        sx={{
          marginBottom: 2,
          display: "flex",
          justifyContent: "center",
        }}
      >
        {Object.keys(counts).map((kindName) => {
          return (
            <Box
              key={kindName}
              sx={{
                marginRight: 2,
                display: "flex",
                "&:last-child": {
                  marginRight: 0,
                },
              }}
            >
              <ApplicationCount
                kindName={kindName}
                enabledCount={
                  counts[kindName][
                    APPLICATION_ACTIVE_STATUS_NAME[
                      ApplicationActiveStatus.ENABLED
                    ]
                  ]
                }
                disabledCount={
                  counts[kindName][
                    APPLICATION_ACTIVE_STATUS_NAME[
                      ApplicationActiveStatus.DISABLED
                    ]
                  ]
                }
                onClick={() => {
                  onClick(APPLICATION_KIND_BY_NAME[kindName]);
                }}
              />
            </Box>
          );
        })}
      </Box>
    );
  }
);
