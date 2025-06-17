import { useContext } from "react";
import { AuthContext, AuthContextType } from "./auth-context";

const useAuth = (): AuthContextType => useContext(AuthContext);

export default useAuth;
