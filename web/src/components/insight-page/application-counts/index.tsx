import makeStyles from "@mui/styles/makeStyles";
import { FC, memo } from "react";
import { APPLICATION_ACTIVE_STATUS_NAME } from "~/constants/application-active-status";
import { APPLICATION_KIND_BY_NAME } from "~/constants/application-kind";
import { useAppSelector } from "~/hooks/redux";
import {
  ApplicationActiveStatus,
  ApplicationKind,
} from "~/modules/applications";
import { ApplicationCount } from "./application-count";

const useStyles = makeStyles((theme) => ({
  root: {
    marginBottom: theme.spacing(2),
    display: "flex",
    justifyContent: "center",
  },
  count: {
    marginRight: theme.spacing(2),
    display: "flex",
    "&:last-child": {
      marginRight: 0,
    },
  },
}));

interface ApplicationCountsProps {
  onClick: (kind: ApplicationKind) => void;
}

export const ApplicationCounts: FC<ApplicationCountsProps> = memo(
  function ApplicationCounts({ onClick }) {
    const classes = useStyles();
    const counts = useAppSelector<Record<string, Record<string, number>>>(
      (state) => state.applicationCounts.counts
    );

    return (
      <div className={classes.root}>
        {Object.keys(counts).map((kindName) => {
          return (
            <ApplicationCount
              key={kindName}
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
              className={classes.count}
            />
          );
        })}
      </div>
    );
  }
);
