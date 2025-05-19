import { styled } from "@mui/material";
import { green, red, yellow } from "@mui/material/colors";

export const LineWrap = styled("div")(({ theme }) => ({
  minHeight: `${theme.typography.body2.lineHeight}em`,
  wordBreak: "break-all",
  whiteSpace: "pre-wrap",
}));

const Line = styled("span")(({ theme }) => ({
  fontFamily: theme.typography.fontFamilyMono,
  fontSize: theme.typography.body2.fontSize,
}));

export const AddedLine = styled(Line)(() => ({
  color: green[800],
  backgroundColor: green[50],
}));

export const DeletedLine = styled(Line)(() => ({
  color: red[800],
  backgroundColor: red[50],
}));

export const ChangedLine = styled(Line)(() => ({
  color: yellow[900],
  backgroundColor: yellow[50],
}));
