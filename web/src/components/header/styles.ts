import { Link, LinkProps, MenuItemProps } from "@mui/material";
import { styled } from "@mui/material/styles";
import { LinkProps as RouterLinkProps } from "react-router-dom";
import { OpenInNew } from "@mui/icons-material";

export type StyledLinkProps = LinkProps &
  Partial<RouterLinkProps> & { active?: boolean };

export const StyledLink = styled(Link, {
  shouldForwardProp: (prop) => prop !== "active",
})<StyledLinkProps>(({ active, theme }) => ({
  display: "inline-flex",
  height: "100%",
  alignItems: "center",
  fontSize: theme.typography.fontSize,
  marginRight: theme.spacing(2),
  "&:hover": {
    color: theme.palette.grey[300],
  },
  textDecoration: "none",
  borderBottom: "2px solid",
  borderColor: active ? theme.palette.background.paper : "transparent",
}));

export const LogoImage = styled("img")({
  height: 56,
});

export const IconOpenNewTab = styled(OpenInNew)({
  fontSize: "0.95rem",
  marginLeft: "5px",
  marginBottom: "-2px",
  color: "rgba(0, 0, 0, 0.5)",
});

export type StyledMenuItemProps = MenuItemProps &
  Partial<RouterLinkProps> & { active?: boolean };
