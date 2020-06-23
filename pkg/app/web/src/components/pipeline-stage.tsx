import React, { FC, memo } from "react";
import { makeStyles, Paper, Typography, Box } from "@material-ui/core";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";
import { StageStatusIcon } from "./stage-status-icon";

const useStyles = makeStyles((theme) => ({
  container: (props: { active: boolean }) => ({
    display: "inline-flex",
    cursor: "pointer",
    "&:hover": {
      backgroundColor: theme.palette.action.hover,
    },
    backgroundColor: props.active
      ? // NOTE: 12%
        theme.palette.primary.main + "1e"
      : undefined,
  }),
  name: {
    marginLeft: theme.spacing(1),
  },
}));

interface Props {
  id: string;
  name: string;
  status: StageStatus;
  active: boolean;
  onClick: (stageId: string, stageName: string) => void;
}

export const PipelineStage: FC<Props> = memo(function PipelineStage({
  id,
  name,
  status,
  onClick,
  active,
}) {
  const classes = useStyles({ active });

  function handleOnClick(): void {
    onClick(id, name);
  }

  return (
    <Paper square className={classes.container} onClick={handleOnClick}>
      <Box alignItems="center" display="flex" justifyContent="center" p={2}>
        <StageStatusIcon status={status} />
        <Typography variant="subtitle2" className={classes.name}>
          <Box fontFamily="Roboto Mono">{name}</Box>
        </Typography>
      </Box>
    </Paper>
  );
});
