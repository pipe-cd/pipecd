import React, { memo, FC, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  fetchApplications,
  selectAll,
  Application,
} from "../../modules/applications";
import { AppState } from "../../modules";
import { Link } from "@material-ui/core";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "../../constants";

export const ApplicationIndexPage: FC = memo(() => {
  const dispatch = useDispatch();
  const applications = useSelector<AppState, Application[]>((state) =>
    selectAll(state.applications)
  );

  useEffect(() => {
    dispatch(fetchApplications());
  }, []);

  return (
    <div>
      <ul>
        {applications.map((application) => (
          <li>
            <Link
              component={RouterLink}
              to={`${PAGE_PATH_APPLICATIONS}/${application.id}`}
            >
              {application.name}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
});
