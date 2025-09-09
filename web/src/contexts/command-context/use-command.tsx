import { useContext } from "react";
import { CommandContext, CommandContextType } from "./command-context";

export const useCommand = (): CommandContextType => {
  return useContext(CommandContext);
};

export default useCommand;
