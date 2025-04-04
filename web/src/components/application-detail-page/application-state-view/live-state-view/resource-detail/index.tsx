import { IconButton, makeStyles, Paper, Typography } from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import { FC } from "react";
import { findMetadataByKey } from "~/utils/find-metadata-by-key";
import { ResourceState } from "~~/model/application_live_state_pb";

const DETAIL_WIDTH = 400;

const useStyles = makeStyles((theme) => ({
  root: {
    width: DETAIL_WIDTH,
    padding: "16px 24px",
    height: "100%",
    overflow: "auto",
    position: "relative",
  },
  closeButton: {
    position: "absolute",
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
  name: {
    paddingRight: theme.spacing(4),
    wordBreak: "break-all",
    paddingBottom: theme.spacing(2),
  },
  section: {
    paddingTop: theme.spacing(1),
    display: "flex",
    alignItems: "center",
  },
  sectionTitle: {
    color: theme.palette.text.secondary,
    minWidth: 120,
  },
  sectionBody: {
    flex: 1,
    wordBreak: "break-all",
  },
  multilineSection: {
    paddingTop: theme.spacing(1),
  },
}));

export interface ResourceDetailProps {
  resource: ResourceState.AsObject;
  onClose: () => void;
}

export const ResourceDetail: FC<ResourceDetailProps> = ({
  resource,
  onClose,
}) => {
  const classes = useStyles();
  return (
    <Paper className={classes.root} square>
      <IconButton className={classes.closeButton} onClick={onClose}>
        <CloseIcon />
      </IconButton>
      <Typography variant="h6" className={classes.name}>
        {resource.name}
      </Typography>

      <div className={classes.section}>
        <Typography variant="subtitle1" className={classes.sectionTitle}>
          Kind
        </Typography>
        <Typography variant="body1" className={classes.sectionBody}>
          {findMetadataByKey(resource.resourceMetadataMap, "Kind")}
        </Typography>
      </div>

      <div className={classes.section}>
        <Typography variant="subtitle1" className={classes.sectionTitle}>
          Namespace
        </Typography>
        <Typography variant="body1" className={classes.sectionBody}>
          {findMetadataByKey(resource.resourceMetadataMap, "Namespace")}
        </Typography>
      </div>

      <div className={classes.section}>
        <Typography variant="subtitle1" className={classes.sectionTitle}>
          Api Version
        </Typography>
        <Typography variant="body1" className={classes.sectionBody}>
          {findMetadataByKey(resource.resourceMetadataMap, "API Version")}
        </Typography>
      </div>

      <div className={classes.multilineSection}>
        <Typography variant="subtitle1" className={classes.sectionTitle}>
          Health Description
        </Typography>
        <Typography variant="body1">
          {resource.healthDescription || "Empty"}
        </Typography>
      </div>
    </Paper>
  );
};
