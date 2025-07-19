import { FC, PropsWithChildren, useEffect, useState } from "react";
import { AuthContext, AuthContextType } from "./auth-context";
import { GetMeResponse } from "~~/api_client/service_pb";
import { useGetMe } from "~/queries/me/use-get-me";

export const AuthProvider: FC<PropsWithChildren<unknown>> = ({ children }) => {
  const [me, setMe] = useState<
    (GetMeResponse.AsObject & { isLogin: boolean }) | null
  >(null);

  const { data, isInitialLoading } = useGetMe({
    retry: false,
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
