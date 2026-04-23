import { MemoryRouter, MemoryRouterProps } from "react-router-dom";
import React from "react";

type Props = MemoryRouterProps;

const MemoryRouterTest = (props: Props): React.ReactElement => {
  return (
    <MemoryRouter {...props}>
      {props.children}
    </MemoryRouter>
  );
};

export default MemoryRouterTest;
