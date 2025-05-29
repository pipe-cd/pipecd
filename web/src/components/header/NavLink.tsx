import React, { FC } from "react";
import { StyledLink, StyledLinkProps } from "./styles";
import { Link, useLocation } from "react-router-dom";

type Props = StyledLinkProps & {
  href: string;
  children: React.ReactNode;
};

const NavLink: FC<Props> = ({ href, children, ...props }) => {
  const location = useLocation();

  return (
    <StyledLink
      component={Link}
      to={href}
      active={location.pathname === href}
      {...props}
      color={"inherit"}
    >
      {children}
    </StyledLink>
  );
};

export default NavLink;
