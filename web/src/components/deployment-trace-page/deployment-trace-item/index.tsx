import {
  Box,
  Collapse,
  IconButton,
  List,
  ListItem,
  makeStyles,
  Typography,
} from "@material-ui/core";
import dayjs from "dayjs";
import React, { FC, useEffect, useMemo, useState } from "react";
import { ListDeploymentTracesResponse } from "~~/api_client/service_pb";
import MoreHorizIcon from "@material-ui/icons/MoreHoriz";
import { ArrowDropDown } from "@material-ui/icons";
import DeploymentItem from "./deployment-item";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "~/constants/path";

const useStyles = makeStyles((theme) => ({
  btnActive: {
    backgroundColor: theme.palette.grey[300],
  },
  btnRotate: {
    transform: "rotate(180deg)",
  },
  list: {
    listStyle: "none",
    padding: theme.spacing(3),
    paddingTop: 0,
    margin: 0,
    flex: 1,
    overflowY: "scroll",
  },
  listItem: {
    borderColor: theme.palette.grey[300],
  },
  traceStickyTop: {
    position: "sticky",
    top: 0,
    zIndex: 50,
    backgroundColor: theme.palette.background.paper,
    paddingBottom: theme.spacing(1),
    borderBottom: `1px solid ${theme.palette.grey[300]}`,
  },
}));

type Props = {
  trace: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject["trace"];
  deploymentList: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject["deploymentsList"];
};

const DeploymentTraceItem: FC<Props> = ({ trace, deploymentList }) => {
  const classes = useStyles();
  const [visibleMessage, setVisibleMessage] = useState(false);
  const [visibleDeployments, setVisibleDeployments] = useState(false);

  useEffect(() => {
    if (visibleDeployments) {
      setVisibleMessage(false);
    }
  }, [visibleDeployments]);

  const onViewCommitMessage = (
    e: React.MouseEvent<HTMLButtonElement>
  ): void => {
    e.stopPropagation();
    setVisibleMessage(!visibleMessage);
  };

  const timeStampCommit = useMemo(() => {
    if (!trace?.commitTimestamp) return "-";
    const diff = dayjs().diff(trace.commitTimestamp, "month");
    const date = dayjs(trace.commitTimestamp);
    const isCurrentYear = dayjs().isSame(date, "year");

    if (!isCurrentYear) {
      return date.format("MMM D, YYYY");
    }
    if (diff > 1) {
      return date.format("MMM D");
    }

    return date.fromNow();
  }, [trace?.commitTimestamp]);

  return (
    <Box flex={1} pl={2} py={2} width={"100%"}>
      <Box
        display="flex"
        flexDirection="row"
        justifyContent="space-between"
        alignItems="center"
        className={visibleDeployments ? classes.traceStickyTop : ""}
      >
        <Box>
          <Box display="flex" gridColumnGap={10} alignItems={"start"}>
            <div>
              <Typography variant="h6">{trace?.title}</Typography>
              <RouterLink to={trace?.commitUrl || "#"} target="_blank">
                <Typography variant="body2" color="textSecondary">
                  {trace?.commitHash}
                </Typography>
              </RouterLink>
            </div>

            <Box display={visibleDeployments ? "none" : "flex"}>
              <IconButton
                size="small"
                className={visibleMessage ? classes.btnActive : ""}
                onClick={onViewCommitMessage}
              >
                <MoreHorizIcon />
              </IconButton>
            </Box>
          </Box>
          <Box
            sx={{ display: !visibleMessage ? "none" : "block" }}
            borderLeft={0.5}
            borderColor={"grey.300"}
            pl={1}
            py={1}
            my={1}
            ml={1}
          >
            <Typography variant="body2" color="textSecondary">
              {trace?.commitMessage}
            </Typography>
          </Box>

          <Box display={"flex"} gridColumnGap={3}>
            <Typography variant="body2" color="textSecondary">
              {trace?.author} authored
            </Typography>
            <Typography
              variant="body2"
              color="textSecondary"
              title={dayjs(trace?.commitTimestamp).format("MMM D, YYYY h:mm A")}
            >
              {timeStampCommit}
            </Typography>
          </Box>
        </Box>
        <IconButton
          aria-label="expand"
          className={visibleDeployments ? classes.btnRotate : ""}
          onClick={() => setVisibleDeployments(!visibleDeployments)}
        >
          <ArrowDropDown />
        </IconButton>
      </Box>

      <Collapse in={visibleDeployments} unmountOnExit key={trace?.id}>
        <List className={classes.listItem}>
          {deploymentList.map((deployment) => (
            <ListItem
              key={deployment?.id}
              button
              dense
              divider
              component={RouterLink}
              to={`${PAGE_PATH_DEPLOYMENTS}/${deployment.id}`}
            >
              <DeploymentItem key={deployment.id} deployment={deployment} />
            </ListItem>
          ))}
        </List>
      </Collapse>
    </Box>
  );
};

export default DeploymentTraceItem;
