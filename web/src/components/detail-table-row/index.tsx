import { FC } from "react";
import * as React from "react";
import { Typography } from "@mui/material";

import makeStyles from "@mui/styles/makeStyles";

const useStyles = makeStyles((theme) => ({
  head: {
    textAlign: "left",
    whiteSpace: "nowrap",
    display: "block",
  },
  value: {
    marginLeft: theme.spacing(1),
  },
}));

export interface DetailTableRowProps {
  label: string;
  value: React.ReactChild;
}

export const DetailTableRow: FC<DetailTableRowProps> = ({ label, value }) => {
  const classes = useStyles();
  return (
    <tr>
      <th className={classes.head}>
        <Typography variant="subtitle2">{label}</Typography>
      </th>
      <td>
        <Typography variant="body2" className={classes.value}>
          {value}
        </Typography>
      </td>
    </tr>
  );
};
