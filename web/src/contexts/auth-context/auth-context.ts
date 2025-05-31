import React from "react";
import { GetMeResponse } from "~~/api_client/service_pb";

export type AuthContextType = {
  me: (GetMeResponse.AsObject & { isLogin: boolean }) | null;
};

export const AuthContext = React.createContext<AuthContextType>({
  me: null,
});
