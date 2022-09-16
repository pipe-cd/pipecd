import {
  Paper,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import EditIcon from "@material-ui/icons/Edit";
import { FC, memo } from "react";
import { RBAC_DESCRIPTION } from "~/constants/text";
import { useAppSelector } from "~/hooks/redux";
import { useProjectSettingStyles } from "~/styles/project-setting";
import {
  RBACPolicy,
  RBAC_RESOURCE_TYPE_TEXT,
  RBAC_ACTION_TYPE_TEXT,
} from "~/modules/project";

const SECTION_TITLE = "Role-Based Access Control";
const RESOURCE_ACTION_SEPARATOR = ";";
const RESOURCES_KEY = "resources";
const ACTIONS_KEY = "actions";

const useStyles = makeStyles(() => ({
  selectTableCell: {
    width: 120,
  },
}));

export const RBACForm: FC = memo(function RBACForm() {
  const classes = useStyles();
  const rbacRoles = useAppSelector((state) => state.project.rbacRoles);
  const projectSettingClasses = useProjectSettingStyles();
  const formalizePoliciesFromPoliciesList = ({
    policiesList,
  }: {
    policiesList: Array<RBACPolicy>;
  }): Array<string> => {
    const policies: Array<string> = [];
    policiesList.map((policy) => {
      const resources: Array<string> = [];
      policy.resourcesList.map((resource) => {
        resources.push(RBAC_RESOURCE_TYPE_TEXT[resource.type]);
      });

      const actions: Array<string> = [];
      policy.actionsList.map((action) => {
        actions.push(RBAC_ACTION_TYPE_TEXT[action]);
      });

      const resource = RESOURCES_KEY + "=" + resources.join(",");
      const action = ACTIONS_KEY + "=" + actions.join(",");
      policies.push(resource + RESOURCE_ACTION_SEPARATOR + action);
    });
    return policies;
  };

  return (
    <>
      <div className={projectSettingClasses.title}>
        <Typography
          variant="h5"
          className={projectSettingClasses.titleWithIcon}
        >
          {SECTION_TITLE}
        </Typography>
      </div>

      <Typography
        variant="body1"
        color="textSecondary"
        className={projectSettingClasses.description}
      >
        {RBAC_DESCRIPTION}
      </Typography>

      <TableContainer component={Paper} square>
        <Table size="small" stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell align="left" />
              <TableCell>Role</TableCell>
              <TableCell>Policies</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rbacRoles.map((role) => (
              <TableRow key={role.name}>
                <TableCell align="left" className={classes.selectTableCell}>
                  <IconButton aria-label="edit" disabled={role.isBuiltin}>
                    <EditIcon />
                  </IconButton>
                </TableCell>
                <TableCell>{role.name}</TableCell>
                <TableCell>
                  {formalizePoliciesFromPoliciesList({
                    policiesList: role.policiesList,
                  }).map((policy, i) => (
                    <p key={i}>{policy}</p>
                  ))}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
});
