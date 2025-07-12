import useAuth from "./use-auth";

const useProjectName = (): string => {
  const { me } = useAuth();
  if (me && me?.isLogin) {
    return me?.projectId;
  }

  return "";
};

export default useProjectName;
