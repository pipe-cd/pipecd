import React, { FC, memo } from "react";
import { makeStyles } from "@material-ui/core";
import { useSelector } from "react-redux";
import { AppState } from "../../modules";
import { ApplicationCount } from "../application-count";
import {
  ApplicationActiveStatus,
  ApplicationKind,
} from "../../modules/applications";
import { APPLICATION_ACTIVE_STATUS_NAME } from "../../constants/application-active-status";
import { APPLICATION_KIND_BY_NAME } from "../../constants/application-kind";

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
    const counts = useSelector<
      AppState,
      Record<string, Record<string, number>>
    >((state) => state.applicationCounts.counts);

    return (
      <div className={classes.root}>
        {Object.keys(counts).map((kindName) => {
          if (
            counts[kindName][
              APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.ENABLED]
            ] === 0 &&
            counts[kindName][
              APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.DISABLED]
            ] === 0
          ) {
            return null;
          }

          return (
            <ApplicationCount
              key={kindName}
              kind={kindName}
              totalCount={
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
