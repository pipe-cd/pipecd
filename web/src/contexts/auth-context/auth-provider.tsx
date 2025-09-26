import { FC, PropsWithChildren, useEffect, useState } from "react";
import { AuthContext, AuthContextType } from "./auth-context";
import { GetMeResponse } from "~~/api_client/service_pb";
import { useGetMe } from "~/queries/me/use-get-me";
import { useLocation } from "react-router-dom";
import {
  LOGIN_ENDPOINT,
  LOGOUT_ENDPOINT,
  PAGE_PATH_LOGIN,
  STATIC_LOGIN_ENDPOINT,
} from "~/constants/path";

const PATH_PUBLIC = [
  PAGE_PATH_LOGIN,
  STATIC_LOGIN_ENDPOINT,
  LOGIN_ENDPOINT,
  LOGOUT_ENDPOINT,
];

export const AuthProvider: FC<PropsWithChildren<unknown>> = ({ children }) => {
  const path = useLocation();
  const [me, setMe] = useState<
    (GetMeResponse.AsObject & { isLogin: boolean }) | null
  >(null);

  const { data, isInitialLoading } = useGetMe({
    retry: false,
    enabled: !PATH_PUBLIC.includes(path.pathname),
  });

  useEffect(() => {
    if (data) {
      setMe({ ...data, isLogin: true });
    } else if (!data && isInitialLoading === false) {
      setMe({ isLogin: false } as AuthContextType["me"]);
    }
  }, [data, isInitialLoading]);

  return <AuthContext.Provider value={{ me }}>{children}</AuthContext.Provider>;
};
