import { Box, IconButton, makeStyles, TextField } from "@material-ui/core";
import { FC, memo, useState } from "react";
import EditIcon from "@material-ui/icons/Edit";
import SaveIcon from "@material-ui/icons/Save";
import ReactMarkdown from "react-markdown";
import { useFormik } from "formik";

interface ApplicationDescriptionProps {
  description: string;
  onUpdate: (description: string) => void;
}

const useStyles = makeStyles(() => ({
  markdown: { flex: 1 },
  textField: {
    flex: 1,
  },
  form: {
    display: "flex",
  },
}));

export const ApplicationDescription: FC<ApplicationDescriptionProps> = memo(
  function ApplicationDescription({ description, onUpdate }) {
    const [editing, setEditing] = useState(false);
    const classes = useStyles();
    const formik = useFormik({
      initialValues: { description },
      async onSubmit({ description }) {
        await onUpdate(description);
        setEditing(false);
      },
    });

    if (editing) {
      return (
        <form onSubmit={formik.handleSubmit} className={classes.form}>
          <TextField
            id="description"
            name="description"
            label="Description"
            multiline
            variant="outlined"
            value={formik.values.description}
            onChange={formik.handleChange}
            placeholder="# Input description by Markdown"
            className={classes.textField}
            disabled={formik.isSubmitting}
          />
          <div>
            <IconButton
              aria-label="Save description"
              type="submit"
              disabled={formik.isSubmitting}
            >
              <SaveIcon />
            </IconButton>
          </div>
        </form>
      );
    }

    return (
      <Box borderLeft="2px solid" borderColor="divider" pl={2} display="flex">
        <ReactMarkdown linkTarget="_blank" className={classes.markdown}>
          {description || "No description."}
        </ReactMarkdown>
        <div>
          <IconButton
            aria-label="Edit description"
            onClick={() => setEditing(!editing)}
          >
            <EditIcon />
          </IconButton>
        </div>
      </Box>
    );
  }
);
